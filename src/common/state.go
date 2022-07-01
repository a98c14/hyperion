package common

import (
	"context"
	"net/http"

	"github.com/a98c14/hyperion/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type State struct {
	Context context.Context
	Conn    *pgxpool.Pool
}

func InitState(r *http.Request) (State, error) {
	ctx := r.Context()
	pool, err := db.GetConnectionPool(ctx)
	if err != nil {
		return State{}, err
	}

	return State{
		Context: ctx,
		Conn:    pool,
	}, nil
}
