package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a98c14/hyperion/db"
	"github.com/a98c14/hyperion/db/query"
	"github.com/a98c14/hyperion/model/render"
)

func GenerateAnimationsFromSprites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := db.GetConnectionPool(ctx)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	rows, err := conn.Query(ctx, `select id, unity_name from sprite`)
	if err != nil {
		http.Error(w, "Could not fetch sprites from database. "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	animationMap := make(map[string]*render.Animation)
	var id int
	var unityName string
	for rows.Next() {
		err := rows.Scan(&id, &unityName)
		if err != nil {
			http.Error(w, "Could not read sprite row."+err.Error(), http.StatusInternalServerError)
			return
		}

		unityName = strings.Replace(unityName, " ", "_", -1)
		unityName = strings.Replace(unityName, "-", "_", -1)
		parts := strings.Split(unityName, "_")
		var animationName string
		if len(parts) < 3 {
			animationName = strings.Join(parts, "_")
		} else {
			animationName = strings.Join(parts[:len(parts)-2], "_")
		}
		if val, ok := animationMap[animationName]; ok {
			val.Sprites = append(val.Sprites, id)
		} else {
			anim := render.Animation{
				Name:           animationName,
				Priority:       1,
				TransitionType: 1,
				Sprites:        make([]int, 0),
			}
			anim.Sprites = append(anim.Sprites, id)
			animationMap[animationName] = &anim
		}
	}

	for _, a := range animationMap {
		err = query.CreateAnimation(ctx, conn, a)
		if err != nil {
			fmt.Println("Error while creating animation" + err.Error())
		}
	}

}
