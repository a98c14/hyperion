package asset

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/response"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func GetAssets(state common.State, w http.ResponseWriter, r *http.Request) error {
	assetTypeString := r.URL.Query().Get("assetType")
	if assetTypeString == "" {
		return errors.ErrBadRequest
	}

	assetType, err := strconv.Atoi(assetTypeString)
	if err != nil {
		return errors.ErrNotFound
	}

	assets, err := DbGetAssets(state, AssetType(assetType))
	if err != nil {
		return err
	}

	fmt.Println(assetTypeString)
	response.Json(w, &assets)
	return nil
}

func GetAssetName(state common.State, w http.ResponseWriter, r *http.Request) error {
	assetIdString := chi.URLParam(r, "assetId")
	assetId, err := strconv.Atoi(assetIdString)
	if err != nil {
		return errors.ErrNotFound
	}

	resp := struct {
		Name string `json:"name"`
	}{}
	resp.Name, err = DbGetAssetName(state, int32(assetId))
	if err != nil {
		return err
	}

	response.Json(w, &resp)
	return nil
}

func SyncAssets(state common.State, w http.ResponseWriter, r *http.Request) error {
	req := struct {
		Type   AssetType
		Assets []AssetUnity
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	batch := &pgx.Batch{}
	for _, asset := range req.Assets {
		DbSyncAsset(state, batch, req.Type, &asset)
	}
	br := state.Conn.SendBatch(state.Context, batch)
	ct, err := br.Exec()
	if err != nil {
		return err
	}
	fmt.Printf("Inserted rows: %d", ct.RowsAffected())

	assets, err := DbGetAssets(state, AssetType(req.Type))
	if err != nil {
		return err
	}

	response.Json(w, assets)
	return nil
}
