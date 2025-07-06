package migrations

import (
	"database/sql"
	"embed"
	"log"

	"github.com/cmrd-a/gophermart/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

func Migrate() {
	db, err := sql.Open("pgx", config.Config.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "sql"); err != nil {
		log.Fatal(err)
	}
}
