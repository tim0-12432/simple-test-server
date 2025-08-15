package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/tim0-12432/simple-test-server/config"
)

var DB *pocketbase.PocketBase

func InitializeDatabase() {
	db := pocketbase.New()
	DB = db

	if err := db.Bootstrap(); err != nil {
		log.Fatalf("Pocketbase bootstrap: %v", err)
	}

	if err := createAdminFromConfig(db); err != nil {
		log.Printf("Failed to create admin user: %v", err)
	}

	go func() {
		if err := apis.Serve(db, apis.ServeConfig{
			HttpAddr:        config.EnvConfig.PbHost + ":" + config.EnvConfig.PbPort,
			ShowStartBanner: false,
		}); err != nil {
			log.Fatalf("pocketbase: %v", err)
		}
	}()

	log.Println("PocketBase started successfully")
}

func CloseDatabase() error {
	if DB == nil {
		return nil
	}

	event := new(core.TerminateEvent)
	event.App = DB

	err := DB.OnTerminate().Trigger(event, func(e *core.TerminateEvent) error {
		return e.App.ResetBootstrapState()
	})
	if err != nil {
		log.Printf("Error during PocketBase shutdown: %v", err)
		return err
	}

	time.Sleep(250 * time.Millisecond)
	log.Println("PocketBase stopped successfully")
	return nil
}

func createAdminFromConfig(pb *pocketbase.PocketBase) error {
	if config.EnvConfig == nil {
		return nil
	}

	adminUser := config.EnvConfig.AdminUser
	adminPass := config.EnvConfig.AdminPass

	if adminUser == "" || adminPass == "" {
		return nil
	}

	var tryNames = []string{"_superusers", "_users", "users"}

	for _, name := range tryNames {
		coll, err := pb.App.FindCollectionByNameOrId(name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return err
		}
		if coll == nil {
			continue
		}

		err = pb.App.RunInTransaction(func(txApp core.App) error {
			var existing *core.Record

			if coll.IsAuth() {
				rec, err := txApp.FindAuthRecordByEmail(coll, adminUser)
				if err == nil {
					existing = rec
				} else if !errors.Is(err, sql.ErrNoRows) {
					return err
				}
			} else {
				rec, err := txApp.FindFirstRecordByData(coll, "email", adminUser)
				if err == nil {
					existing = rec
				} else if !errors.Is(err, sql.ErrNoRows) {
					return err
				}
			}

			if existing != nil {
				existing.Set("email", adminUser)
				existing.Set("username", adminUser)
				existing.Set("verified", true)
				existing.Set("system", true)
				existing.Set(core.FieldNamePassword, adminPass)

				if err := txApp.SaveWithContext(context.Background(), existing); err != nil {
					return err
				}

				log.Println("Marked existing user as system superuser")
				return nil
			}

			newRec := core.NewRecord(coll)
			newRec.Set("email", adminUser)
			newRec.Set("username", adminUser)
			newRec.Set("verified", true)
			newRec.Set("system", true)
			newRec.Set(core.FieldNamePassword, adminPass)

			if err := txApp.SaveWithContext(context.Background(), newRec); err != nil {
				return err
			}

			log.Println("Created admin superuser")
			return nil
		})

		return err
	}

	log.Println("no suitable collection found for admin creation; skipping")
	return nil
}
