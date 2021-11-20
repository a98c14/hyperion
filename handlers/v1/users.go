package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

type Person struct {
	Name string `db:"name"`
}

type UserHandler struct{}

type UserResponse struct {
	X int `json:"Value"`
}

type UserRequest struct {
	Id int
}

const CONNECTION_STRING string = "postgres://postgres:123@127.0.0.1:5432/hyperion"

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), CONNECTION_STRING)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var u UserRequest
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println("Invalid request!", r.Body)
		return
	}
	var name string
	err = conn.QueryRow(context.Background(), "select name from users where id=$1", u.Id).Scan(&name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name)
}
