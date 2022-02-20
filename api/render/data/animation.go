package data

import (
	"database/sql"
	"fmt"

	"github.com/a98c14/hyperion/api/asset"
	"github.com/a98c14/hyperion/common"
	"github.com/jackc/pgx/v4"
)

type Animation struct {
	Id             int    `json:"id"`
	AssetId        int    `json:"assetId"`
	Name           string `json:"name"`
	Priority       int    `json:"priority"`
	TransitionType int    `json:"transitionType"`
	Sprites        []int  `json:"sprites"`
}

type ByName []*Animation

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func CreateAnimation(state common.State, animation *Animation) error {
	var id sql.NullInt32

	// Check if animation exists
	err := state.Conn.QueryRow(state.Context, `select animation.id from animation inner join asset on asset.id=animation.asset_id where asset.name=$1`, animation.Name).Scan(&id)
	if err != nil || !id.Valid {
		fmt.Println("Creating animation: ", animation.Name)
		// If animation doesn't exist insert it
		err := state.Conn.QueryRow(state.Context,
			`
			with ins as (
				insert into asset (name, unity_guid, unity_internal_id, type)
				values ($1, '0', -4, $4)
				returning id
			)
			insert into "animation"
			(asset_id, priority, transition_type)
			values((select id from ins), $2, $3) 
			on conflict do nothing 
			returning id`,
			animation.Name, animation.Priority, animation.TransitionType, asset.Animation).Scan(&id)
		if err != nil {
			return err
		}
	}

	batch := &pgx.Batch{}
	batch.Queue(`delete from animation_sprite where animation_id=$1`, id)
	for _, sprite := range animation.Sprites {
		batch.Queue(`insert into "animation_sprite"
			(animation_id, sprite_id) 
			values($1, $2)`, id, sprite)
	}

	br := state.Conn.SendBatch(state.Context, batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	br.Close()
	return nil
}
