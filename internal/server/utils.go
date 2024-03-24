package server

import (
	"context"

	"github.com/anantadwi13/gorong2/pkg/utils"
)

type key string

const (
	keyConnID key = "conn_id"
)

func getConnID(ctx context.Context) string {
	val, ok := ctx.Value(keyConnID).(string)
	if !ok || val == "" {
		return ""
	}
	return val
}

func initConnID(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyConnID, utils.GenerateUniqueID())
}
