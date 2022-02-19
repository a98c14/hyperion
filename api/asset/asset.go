package asset

import (
	"database/sql"

	"github.com/a98c14/hyperion/common"
	e "github.com/a98c14/hyperion/common/errors"
	"github.com/jackc/pgx/v4"
)

type AssetDb struct {
	Id              int32     `json:"id"`
	Guid            string    `json:"guid"`
	UnityGuid       string    `json:"unityGuid"`
	UnityInternalId int64     `json:"unityInternalId"`
	Name            string    `json:"name"`
	Type            AssetType `json:"type"`
}

type Asset struct {
	Id   int `json:"id"`
	Name int `json:"name"`
}

type AssetUnity struct {
	Name         string `json:"name"`
	Guid         string `json:"guid"`
	InternalGuid string `json:"internalGuid"`
	InternalId   int64  `json:"internalId"`
}

type AssetType int32

const (
	Material          AssetType = 0
	MaterialAnimation AssetType = 1
	ParticleSystem    AssetType = 2
	TrailSystem       AssetType = 3
	ItemPool          AssetType = 4
	Sprite            AssetType = 5
	Texture           AssetType = 6
	Animation         AssetType = 7
	Prefab            AssetType = 8
)

func DbGetAssetName(state common.State, id int32) (string, error) {
	var name string
	err := state.Conn.QueryRow(state.Context, `select name from asset where id=$1`, id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

// TODO(selim): Add paging
func DbGetAssets(state common.State, assetType AssetType) ([]AssetDb, error) {
	rows, err := state.Conn.Query(state.Context, `select id, name, guid, unity_guid, unity_internal_id from asset where type=$1 and deleted_date is null`, assetType)
	if err != nil {
		return nil, e.Wrap("DbGetAssets", err)
	}
	defer rows.Close()

	assets := make([]AssetDb, 0)
	for rows.Next() {
		a := AssetDb{}
		err = rows.Scan(&a.Id, &a.Name, &a.Guid, &a.UnityGuid, &a.UnityInternalId)
		if err != nil {
			return nil, e.Wrap("DbGetAssets", err)
		}
		assets = append(assets, a)
	}

	return assets, nil
}

// Inserts the asset if it doesn't exist. Otherwise updates the name or guid whichever don't match
// TODO(selim): Use batching
func DbSyncAsset(state common.State, assetType AssetType, asset *AssetUnity) (int32, error) {
	var id sql.NullInt32
	var name string
	var guid string
	var internalId int64
	// CHECK(selim): Is checking the name correct here? It is not guaranteed to be unique.
	err := state.Conn.QueryRow(state.Context, `select id, name, unity_guid, unity_internal_id from asset where (name=$1 or unity_internal_id=$2) and type=$3`, asset.Name, asset.InternalId, assetType).Scan(&id, &name, &guid, &internalId)
	if err != nil || !id.Valid {
		err = state.Conn.QueryRow(state.Context, `insert into asset (name, unity_guid, unity_internal_id, type) values ($1, $2, $3, $4) returning id`, asset.Name, asset.InternalGuid, asset.InternalId, assetType).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id.Int32, nil
	}
	if name != asset.Name {
		state.Conn.Exec(state.Context, `update asset (name) values ($2) where id=$1`, id, asset.Name)
	}

	if guid != asset.InternalGuid {
		state.Conn.Exec(state.Context, `update asset (unity_guid) values ($2) where id=$1`, id, asset.InternalGuid)
	}

	if internalId != asset.InternalId {
		state.Conn.Exec(state.Context, `update asset (internalId) values ($2) where id=$1`, id, asset.InternalId)
	}

	return id.Int32, nil
}

// Creates given asset in database
func DbCreateAsset(state common.State, assetType AssetType, asset *AssetUnity) (int32, error) {
	var id int32
	err := state.Conn.QueryRow(state.Context, `insert into "asset" (name, type, unity_guid) 
		values($1, $2, $3) returning id`, asset.Name, assetType, asset.Name).Scan(&id)
	if err != nil {
		return 0, e.Wrap("DbCreateAsset", err)
	}

	return id, nil
}

// Queues the insertion of asset to given batch
func DbCreateAssetBatched(state common.State, batch *pgx.Batch, asset *AssetDb) {
	batch.Queue(`insert into "asset" (name, type, unity_guid) values($1, $2, $3) returning id`,
		asset.Name, asset.Type, asset.UnityGuid)
}
