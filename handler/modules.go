package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

/* Gets the base module_part and all of its children by id.
Can only query base components. Nodes that have null as children
represent leaf nodes */
func GetModuleById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := pgxpool.Connect(ctx, db.ConnectionString)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	// Load components that have given name from database
	componentIdString := chi.URLParam(r, "componentId")
	componentId, err := strconv.Atoi(componentIdString)
	if err != nil {
		http.Error(w, "Could not parse id value!", http.StatusBadRequest)
		return
	}

	rows, err := conn.Query(ctx, `with recursive module_part_recursive as (
			select id, name, value_type, parent_id from module_part
			where id=$1 and parent_id is null
			union select c.id, c.name, c.value_type, c.parent_id from module_part c inner join module_part_recursive cp on cp.id=c.parent_id 
		) select * from module_part_recursive;`, componentId)

	if err != nil {
		http.Error(w, "Could not fetch data from database!", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type modulePart struct {
		Id        int
		Name      string
		ValueType int
		ParentId  int
	}
	moduleParts := make([]modulePart, 0, 10)

	var id int
	var name string
	var valueType int
	var parentId sql.NullInt32
	for rows.Next() {
		err = rows.Scan(&id, &name, &valueType, &parentId)
		if err != nil {
			http.Error(w, "Error while fetching from database!"+err.Error(), http.StatusInternalServerError)
			return
		}
		pid := 0
		if parentId.Valid {
			pid = int(parentId.Int32)
		}
		moduleParts = append(moduleParts, modulePart{Id: id, Name: name, ValueType: valueType, ParentId: pid})
	}
	componentCount := rows.CommandTag().RowsAffected()
	if componentCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	type componentResponse struct {
		Id        int                  `json:"id"`
		Name      string               `json:"name"`
		ValueType int                  `json:"value_type"`
		Children  []*componentResponse `json:"children"`
	}

	/* Create object tree from module part list. Components are always ordered from root to child*/
	root := componentResponse{
		Id:        moduleParts[0].Id,
		Name:      moduleParts[0].Name,
		ValueType: moduleParts[0].ValueType,
		Children:  make([]*componentResponse, 0),
	}
	nodeMap := make(map[int]*componentResponse, componentCount)
	nodeMap[root.Id] = &root
	for _, c := range moduleParts[1:] {
		if val, ok := nodeMap[c.ParentId]; ok {
			cr := componentResponse{
				Id:        c.Id,
				Name:      c.Name,
				ValueType: c.ValueType,
				Children:  nil,
			}
			val.Children = append(val.Children, &cr)
			nodeMap[c.Id] = &cr
		} else {
			http.Error(w, "Could not create object tree from components! Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&root)
}

func ListComponents(w http.ResponseWriter, r *http.Request) {

}

// Returns all the root components.
// Root components are module_part that have no parent. Used to group
// other components and editor can only filter using root components
func GetRootComponents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := pgxpool.Connect(ctx, db.ConnectionString)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	rows, err := conn.Query(ctx, `select id, name from "module_part" where parent_id is null`)
	if err != nil {
		http.Error(w, "Could not fetch data from database!", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type modulePart struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	components := make([]modulePart, 0, 100)

	var id int
	var name string
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			http.Error(w, "Error while fetching from database! "+err.Error(), http.StatusInternalServerError)
			return
		}
		components = append(components, modulePart{Id: id, Name: name})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&components)
}

func DeleteComponent(w http.ResponseWriter, r *http.Request) {
}

func UpdateComponent(w http.ResponseWriter, r *http.Request) {

}

type node struct {
	id         int
	value_type int
	parent_id  sql.NullInt32
	key        string
	value      json.RawMessage
}

// Creates a new module part structure
func CreateComponent(w http.ResponseWriter, r *http.Request) {
	type createComponentRequest struct {
		Name      string
		Structure json.RawMessage
	}

	type partInfo struct {
		ValueType int
		Child     json.RawMessage
	}

	// Parse request body
	var req createComponentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Could not parse request body!", http.StatusBadRequest)
		return
	}

	// Start database connection.
	ctx := r.Context()
	conn, err := db.GetConnectionPool(ctx)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	var exists bool
	err = conn.QueryRow(ctx, "select exists(select 1 from module_part where name=$1)", req.Name).Scan(&exists)
	if exists || err != nil {
		http.Error(w, "Given module part already exists!", http.StatusBadRequest)
		return
	}

	// Start transcation. If all components can not be added successfully, don't
	// insert anything
	tx, err := conn.Begin(ctx)
	if err != nil {
		http.Error(w, "Could not start transaction!", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	// Create processing channel
	// TODO(selim): Is 500 correct for channel size?
	c := make(chan node, 500)
	c <- node{
		id:         0,
		value_type: 0,
		parent_id:  sql.NullInt32{},
		key:        req.Name,
		value:      req.Structure,
	}
	processedObjects := 0

	// Process nodes in json object tree
	for n := range c {
		// Insert current node to database and store its id
		id, err := insertComponent(conn, ctx, n)
		if err != nil {
			http.Error(w, "Could not insert module part with key: "+n.key+" Error:"+err.Error(), http.StatusBadRequest)
			return
		}
		processedObjects++

		// Check if current value is a json object
		m := make(map[string]json.RawMessage, 20)
		err = json.Unmarshal(n.value, &m)
		if err != nil {
			// If there is no more elements to process, close the channel
			if len(c) == 0 {
				close(c)
			}

			// If there is an error during json parsing for root node, exit early with error
			if processedObjects == 1 {
				http.Error(w, "Root object doesn't have a valid json", http.StatusBadRequest)
				return
			}
			continue
		}

		// Add all values to process channel
		for k := range m {
			var pi partInfo
			err = json.Unmarshal(m[k], &pi)
			if err != nil {
				http.Error(w, "Invalid module part structure!", http.StatusBadRequest)
				return
			}
			c <- node{
				parent_id:  id,
				key:        k,
				value_type: pi.ValueType,
				value:      pi.Child,
			}
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Commit failed with message: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully created module part!")
}

func insertComponent(conn *pgxpool.Pool, ctx context.Context, node node) (sql.NullInt32, error) {
	var id sql.NullInt32
	err := conn.QueryRow(ctx,
		`insert into "module_part" 
		 (name, value_type, parent_id) 
		 values($1, $2, $3) returning id`,
		node.key, node.value_type, node.parent_id).Scan(&id)
	return id, err
}
