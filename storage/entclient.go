package storage

import (
	"context"
	"log"
	"somapay-backend/ent"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	client *ent.Client
	once   sync.Once
)

func GetClient() *ent.Client {
	once.Do(func() {
		var err error
		client, err = ent.Open("mysql", "root:1234@tcp(localhost:3306)/somapay?parseTime=True")
		if err != nil {
			log.Fatalf("failed opening connection to mysql: %v", err)
		}

		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	})

	return client
}
