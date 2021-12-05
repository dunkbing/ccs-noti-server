package main

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
)

func getDeviceTokens(db *firestore.Client, ctx context.Context, coll, doc string) ([]string, error) {
	var tokens []string
	snapshot, err := db.Collection(coll).Doc(doc).Get(ctx)
	if err != nil {
		return nil, err
	}
	if snapshot.Exists() {
		data := snapshot.Data()
		if data["tokens"] != nil {
			for _, token := range data["tokens"].([]interface{}) {
				tokens = append(tokens, token.(string))
			}
			return tokens, nil
		} else {
			return nil, errors.New("no tokens found")
		}
	} else {
		return nil, errors.New("no data found")
	}
}
