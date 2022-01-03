package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/a98c14/hyperion/db"
	"github.com/a98c14/hyperion/db/query"
	"github.com/go-chi/chi/v5"
)

func GetTextures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := db.GetConnectionPool(ctx)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	rows, err := conn.Query(ctx, "select id, unity_name from texture")
	if err != nil {
		http.Error(w, "Could not fetch textures,"+err.Error(), http.StatusInternalServerError)
		return
	}

	type textureResponse struct {
		Id   int
		Name string
	}

	var id int
	var name string
	result := make([]textureResponse, 0)
	for rows.Next() {
		rows.Scan(&id, &name)
		result = append(result, textureResponse{
			Id:   id,
			Name: name,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&result)
}

func GetTextureFile(w http.ResponseWriter, r *http.Request) {
	type textureRequest struct {
		Id int
	}

	ctx := r.Context()

	textureIdString := chi.URLParam(r, "textureId")
	textureId, err := strconv.Atoi(textureIdString)
	if err != nil {
		http.Error(w, "Could not parse id value!", http.StatusBadRequest)
		return
	}
	conn, err := db.GetConnectionPool(ctx)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	var imagePath string
	err = conn.QueryRow(ctx, "select image_path from texture where id=$1", textureId).Scan(&imagePath)
	if err != nil {
		http.Error(w, "Could not fetch texture,"+err.Error(), http.StatusBadRequest)
		return
	}

	file, err := os.ReadFile(imagePath)
	if err != nil {
		http.Error(w, "Error while reading the file,"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

func CreateTexture(w http.ResponseWriter, r *http.Request) {
	type textureResponse struct {
		Id int
	}

	r.ParseMultipartForm(100 << 20)
	unityName := r.FormValue("name")
	unityGuid := r.FormValue("guid")

	if unityName == "" {
		http.Error(w, "Texture name should not be empty!", http.StatusBadRequest)
		return
	}

	if unityGuid == "" {
		http.Error(w, "Texture guid should not be empty!", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("texture")
	if err != nil {
		http.Error(w, "Could not read form file named texture, Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Connect to database
	ctx := r.Context()
	conn, err := db.GetConnectionPool(ctx)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	// Check if texture with given name already existing
	var existing int
	conn.QueryRow(ctx, "select id from texture where unity_name=$1 limit 1", unityName).Scan(&existing)
	if existing != 0 {
		w.WriteHeader(http.StatusOK)
		resp := textureResponse{
			Id: existing,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(".\\textures", unityName+"-*.png")
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer tempFile.Close()

	// Read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	path := tempFile.Name()

	id, err := query.InsertTexture(ctx, conn, path, unityGuid, unityName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not insert texture to database, Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := textureResponse{
		Id: int(id.Int32),
	}
	json.NewEncoder(w).Encode(&resp)
}
