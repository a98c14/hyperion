package asset

import (
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
func DbSyncAsset(state common.State, batch *pgx.Batch, assetType AssetType, asset *AssetUnity) {
	batch.Queue(`
	with new_assets (name, unity_guid, unity_internal_id, type) as 
			 (values ($1, $2, $3::bigint, $4::int)),
		 upsert as (
			 update asset a 
				 set name=na.name,
					 unity_guid=na.unity_guid,
					unity_internal_id=na.unity_internal_id
			 from new_assets na 
			 where (a.unity_internal_id=na.unity_internal_id or a.name=na.name) and a.type=na.type
			 returning a.name, a.unity_guid, a.unity_internal_id, a.type)
		insert into asset (name, unity_guid, unity_internal_id, type)
		select name, unity_guid, unity_internal_id, type from new_assets
		where not exists (select 1 from upsert where upsert.unity_internal_id=new_assets.unity_internal_id and upsert.type=new_assets.type)
`, asset.Name, asset.InternalGuid, asset.InternalId, assetType)
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
