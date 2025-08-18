package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/pocketbase/pocketbase"
)

var collections = []string{
	"containers",
}

func InitializeCollections(pb *pocketbase.PocketBase) error {
	if pb == nil || pb.App == nil {
		return errors.New("pocketbase not initialized")
	}

	for _, collName := range collections {
		if err := EnsureContainersCollection(pb, collName); err != nil {
			log.Printf("Error ensuring collection %s: %v", collName, err)
			return err
		}
	}

	return nil
}

func EnsureContainersCollection(pb *pocketbase.PocketBase, collectionName string) error {
	if _, err := pb.App.FindCollectionByNameOrId(collectionName); err == nil {
		return nil
	} else if errors.Is(err, sql.ErrNoRows) {
		log.Printf("%s collection not found; will defer creation to migrations", collectionName)
		return nil
	} else {
		return err
	}
}
