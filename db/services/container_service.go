package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/tim0-12432/simple-test-server/db"
	"github.com/tim0-12432/simple-test-server/db/dtos"
)

const containersCollectionName = "containers"

func CreateContainer(c *dtos.Container) (string, error) {
	if db.DB == nil {
		return "", errors.New("pocketbase not initialized")
	}

	var newID string
	err := db.DB.App.RunInTransaction(func(txApp core.App) error {
		coll, err := txApp.FindCollectionByNameOrId(containersCollectionName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("collection %s not found", containersCollectionName)
			}
			return err
		}

		rec := core.NewRecord(coll)
		rec.Set("name", c.Name)
		rec.Set("image", c.Image)
		rec.Set("created_at", c.CreatedAt)
		rec.Set("environment", c.Environment)
		rec.Set("ports", c.Ports)
		rec.Set("volumes", c.Volumes)
		rec.Set("networks", c.Networks)
		rec.Set("status", c.Status)

		if err := txApp.SaveWithContext(context.Background(), rec); err != nil {
			return err
		}

		newID = rec.Id
		return nil
	})

	return newID, err
}

func ListContainers() ([]*dtos.Container, error) {
	if db.DB == nil {
		return nil, errors.New("pocketbase not initialized")
	}

	coll, err := db.DB.App.FindCollectionByNameOrId(containersCollectionName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("collection %s not found", containersCollectionName)
		}
		return nil, err
	}

	recs := make([]*core.Record, 0)
	q := db.DB.App.RecordQuery(coll)
	if err := q.All(&recs); err != nil {
		return nil, err
	}

	out := make([]*dtos.Container, 0, len(recs))
	for _, r := range recs {
		c := &dtos.Container{}
		c.ID = r.Id
		c.Name = db.ToString(r.Get("name"))
		c.Image = db.ToString(r.Get("image"))
		c.CreatedAt = db.ToString(r.Get("created_at"))
		c.Status = dtos.ToStatus(db.ToString(r.Get("status")))

		c.Environment = db.ToStringMap(r.Get("environment"))
		c.Volumes = db.ToStringMap(r.Get("volumes"))
		c.Networks = db.ToStringSlice(r.Get("networks"))
		c.Ports = db.ToIntMap(r.Get("ports"))

		out = append(out, c)
	}

	return out, nil
}

func GetContainer(id string) (*dtos.Container, error) {
	if db.DB == nil {
		return nil, errors.New("pocketbase not initialized")
	}

	coll, err := db.DB.App.FindCollectionByNameOrId(containersCollectionName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("collection %s not found", containersCollectionName)
		}
		return nil, err
	}

	rec, err := db.DB.App.FindRecordById(coll, id)
	if err != nil {
		return nil, err
	}

	c := &dtos.Container{}
	c.ID = rec.Id
	c.Name = db.ToString(rec.Get("name"))
	c.Image = db.ToString(rec.Get("image"))
	c.CreatedAt = db.ToString(rec.Get("created_at"))
	c.Status = dtos.ToStatus(db.ToString(rec.Get("status")))

	c.Environment = db.ToStringMap(rec.Get("environment"))
	c.Volumes = db.ToStringMap(rec.Get("volumes"))
	c.Networks = db.ToStringSlice(rec.Get("networks"))
	c.Ports = db.ToIntMap(rec.Get("ports"))

	return c, nil
}

func UpdateContainer(id string, updates *dtos.Container) error {
	if db.DB == nil {
		return errors.New("pocketbase not initialized")
	}

	return db.DB.App.RunInTransaction(func(txApp core.App) error {
		coll, err := txApp.FindCollectionByNameOrId(containersCollectionName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("collection %s not found", containersCollectionName)
			}
			return err
		}

		rec, err := txApp.FindRecordById(coll, id)
		if err != nil {
			return err
		}

		// set fields from updates (overwrite)
		rec.Set("name", updates.Name)
		rec.Set("image", updates.Image)
		rec.Set("created_at", updates.CreatedAt)
		rec.Set("environment", updates.Environment)
		rec.Set("ports", updates.Ports)
		rec.Set("volumes", updates.Volumes)
		rec.Set("networks", updates.Networks)
		rec.Set("status", updates.Status)

		if err := txApp.SaveWithContext(context.Background(), rec); err != nil {
			return err
		}

		return nil
	})
}

func DeleteContainer(id string) error {
	if db.DB == nil {
		return errors.New("pocketbase not initialized")
	}

	return db.DB.App.RunInTransaction(func(txApp core.App) error {
		coll, err := txApp.FindCollectionByNameOrId(containersCollectionName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("collection %s not found", containersCollectionName)
			}
			return err
		}

		deleteQuery := fmt.Sprintf("DELETE FROM _records WHERE id = '%s' AND collection_id = '%s'", db.EscapeSQL(id), db.EscapeSQL(coll.Id))
		if _, err := txApp.DB().NewQuery(deleteQuery).Execute(); err != nil {
			return err
		}

		return nil
	})
}
