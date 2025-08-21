package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// prepare a collections snapshot for import (ImportCollectionsByMarshaledJSON)
		coll := []map[string]any{
			{
				"name": "containers",
				"type": "base",
				"fields": []map[string]any{
					{"name": "container_id", "type": "text", "required": true, "unique": true, "options": map[string]any{}},
					{"name": "name", "type": "text", "required": true, "unique": true, "options": map[string]any{}},
					{"name": "image", "type": "text", "required": false, "unique": false, "options": map[string]any{}},
					{"name": "created_at", "type": "text", "required": false, "unique": false, "options": map[string]any{}},
					{"name": "environment", "type": "json", "required": false, "unique": false, "options": map[string]any{}},
					{"name": "ports", "type": "json", "required": false, "unique": false, "options": map[string]any{}},
					{"name": "volumes", "type": "json", "required": false, "unique": false, "options": map[string]any{}},
					{"name": "networks", "type": "json", "required": false, "unique": false, "options": map[string]any{}},
					{"name": "status", "type": "text", "required": false, "unique": false, "options": map[string]any{}},
					{"name": "type", "type": "text", "required": false, "unique": false, "options": map[string]any{}},
				},
			},
		}

		b, err := json.Marshal(coll)
		if err != nil {
			return err
		}

		return app.ImportCollectionsByMarshaledJSON(b, false)
	}, func(app core.App) error {
		// down: remove the collection if exists
		c, err := app.FindCollectionByNameOrId("containers")
		if err != nil {
			return err
		}
		return app.Delete(c)
	}, "1687800000_create_containers_collection.go")
}
