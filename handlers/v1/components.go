package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/db"
	"github.com/a98c14/hyperion/utils"
	"github.com/jackc/pgx/v4"
)

type ComponentsHandler struct{}

func (h *ComponentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)

	if head == "" {
		switch r.Method {
		case utils.GET:
			handleGetComponents(w)
		case utils.POST:
			handleCreateComponent(w, r)
		}
		return
	}

	// Get by Id
	id, err := strconv.Atoi(head)
	if err != nil {
		http.Error(w, "Invalid component id: "+head, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case utils.GET:
		handleGetComponent(id, w)
	}

}

func handleGetComponent(id int, w http.ResponseWriter) {

}

func handleGetComponents(w http.ResponseWriter) {

}

type createComponentRequest struct {
	Name      string
	Structure json.RawMessage
}

type node struct {
	id        int
	parent_id sql.NullInt32
	key       string
	val       json.RawMessage
}

func handleCreateComponent(w http.ResponseWriter, r *http.Request) {
	// Start database connection.
	// TODO(selim): Use connection pool
	ctx := r.Context()
	conn, err := pgx.Connect(ctx, db.ConnectionString)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
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

	// Parse request body
	var req createComponentRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Could not parse request body!", http.StatusBadRequest)
		return
	}

	// Create processing channel
	// TODO(selim): Is 500 correct for channel size?
	c := make(chan node, 500)
	c <- node{
		id:        0,
		parent_id: sql.NullInt32{},
		key:       req.Name,
		val:       req.Structure,
	}
	processedObjects := 0

	// Process nodes in json object tree
	for n := range c {
		// Insert current node to database and store its id
		id, err := insertComponent(conn, ctx, n)
		if err != nil {
			http.Error(w, "Could not insert component with key: "+n.key, http.StatusBadRequest)
			return
		}
		processedObjects++

		// Check if current value is a json object
		m := make(map[string]json.RawMessage, 20)
		err = json.Unmarshal(n.val, &m)
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
			c <- node{
				parent_id: id,
				key:       k,
				val:       m[k],
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Commit failed with message: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully created component!")
}

func insertComponent(conn *pgx.Conn, ctx context.Context, node node) (sql.NullInt32, error) {
	var id sql.NullInt32
	err := conn.QueryRow(ctx,
		`insert into "component" 
		 (name, type, parent_id, is_hidden) 
		 values($1, $2, $3, $4) returning id`,
		node.key, 0, node.parent_id, false).Scan(&id)
	return id, err
}
