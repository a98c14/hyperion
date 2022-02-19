package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/a98c14/hyperion/api/render/data"
	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/response"
	"github.com/go-chi/chi/v5"
)

func GetTextures(state common.State, w http.ResponseWriter, r *http.Request) error {
	rows, err := state.Conn.Query(state.Context, "select id, unity_name from texture")
	if err != nil {
		return err
	}

	type textureResponse struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
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

	response.Json(w, &result)

	return nil
}

func GetTextureFile(state common.State, w http.ResponseWriter, r *http.Request) error {
	textureIdString := chi.URLParam(r, "textureId")
	textureId, err := strconv.Atoi(textureIdString)
	if err != nil {
		return err
	}

	var imagePath string
	err = state.Conn.QueryRow(state.Context, "select image_path from texture where id=$1", textureId).Scan(&imagePath)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(file)

	return nil
}

func CreateTexture(state common.State, w http.ResponseWriter, r *http.Request) error {
	type textureResponse struct {
		Id int `json:"id"`
	}

	r.ParseMultipartForm(100 << 20)
	unityName := r.FormValue("name")
	unityGuid := r.FormValue("guid")

	if unityName == "" {
		return errors.New("texture name should not be empty")
	}

	if unityGuid == "" {
		return errors.New("texture guid should not be empty")
	}

	file, handler, err := r.FormFile("texture")
	if err != nil {
		return errors.New("could not read form file named texture")
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Check if texture with given name already existing
	var existing int
	state.Conn.QueryRow(state.Context, "select id from texture where unity_name=$1 limit 1", unityName).Scan(&existing)
	if existing != 0 {
		resp := textureResponse{
			Id: existing,
		}
		response.Json(w, &resp)
		return nil
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(".\\textures", unityName+"-*.png")
	if err != nil {
		return err
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

	id, err := data.InsertTexture(state, path, unityGuid, 0, unityName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp := textureResponse{
		Id: int(id.Int32),
	}

	response.Json(w, &resp)
	return nil
}
