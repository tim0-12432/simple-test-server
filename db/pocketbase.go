package db

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/tim0-12432/simple-test-server/config"
)

var DB *pocketbase.PocketBase

func InitializeDatabase() {
	db := pocketbase.New()

	go func() {
		if err := apis.Serve(db, apis.ServeConfig{
			HttpAddr:        config.EnvConfig.PbHost + ":" + config.EnvConfig.PbPort,
			ShowStartBanner: false,
		}); err != nil {
			log.Fatalf("pocketbase: %v", err)
		}
	}()

	DB = db

	log.Println("PocketBase started successfully")
}

func CloseDatabase() error {
	// if DB != nil {
	// 	if err := DB.Stop(); err != nil {
	// 		log.Printf("Error stopping PocketBase: %v", err)
	// 	}
	// }
	log.Println("PocketBase stopped successfully")
	return nil
}
