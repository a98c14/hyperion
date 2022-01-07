package data

import (
	"database/sql"

	"github.com/a98c14/hyperion/common"
	"github.com/jackc/pgx/v4"
)

type Animation struct {
	Id             int    `json:"id"`
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
	err := state.Conn.QueryRow(state.Context,
		`insert into "animation"
		(name, priority, transition_type)
		values($1, $2, $3) 
		on conflict do nothing 
		returning id`,
		animation.Name, animation.Priority, animation.TransitionType).Scan(&id)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
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
