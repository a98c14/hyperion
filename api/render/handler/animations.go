package handler

import (
	"net/http"
	"sort"
	"strings"

	"github.com/a98c14/hyperion/api/render/data"
	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/response"
)

func GetAnimations(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}

	rows, err := state.Conn.Query(state.Context,
		`select a.id, asset.name, a.priority, a.transition_type from animation a
		 inner join asset on asset.id=a.asset_id`)

	if err != nil {
		response.InternalError(w, err)
		return
	}
	defer rows.Close()

	// Query animations
	animationMap := make(map[int]*data.Animation)
	var id int
	var name string
	var priority int
	var transitionType int
	for rows.Next() {
		err = rows.Scan(&id, &name, &priority, &transitionType)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		animationMap[id] = &data.Animation{
			Id:             id,
			Name:           name,
			Priority:       priority,
			TransitionType: transitionType,
			Sprites:        make([]int, 0),
		}
	}

	rows, err = state.Conn.Query(state.Context,
		`select animation_id, sprite_id from animation_sprite`)

	if err != nil {
		response.InternalError(w, err)
		return
	}
	defer rows.Close()

	var animationId int
	var spriteId int
	for rows.Next() {
		err = rows.Scan(&animationId, &spriteId)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		if val, ok := animationMap[animationId]; ok {
			val.Sprites = append(val.Sprites, spriteId)
		}
	}

	animations := make([]*data.Animation, 0, len(animationMap))
	for _, v := range animationMap {
		animations = append(animations, v)
	}

	sort.Sort(data.ByName(animations))
	response.Json(w, animations)
}

func GenerateAnimationsFromSprites(state common.State, w http.ResponseWriter, r *http.Request) error {
	conn := state.Conn
	ctx := state.Context

	rows, err := conn.Query(ctx, `select id, unity_name from sprite`)
	if err != nil {
		return err
	}
	defer rows.Close()

	animationMap := make(map[string]*data.Animation)
	var id int
	var unityName string
	for rows.Next() {
		err := rows.Scan(&id, &unityName)
		if err != nil {
			return err
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
			anim := data.Animation{
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
		err = data.CreateAnimation(state, a)
		if err != nil {
			return err
		}
	}

	return nil
}
