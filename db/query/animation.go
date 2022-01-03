package query

import (
	"context"
	"database/sql"

	"github.com/a98c14/hyperion/model/render"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateAnimation(ctx context.Context, pool *pgxpool.Pool, animation *render.Animation) error {
	var id sql.NullInt32
	err := pool.QueryRow(ctx,
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

	br := pool.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	br.Close()
	return nil
}
