package main

import (
	"log/slog"
	"os"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/sqlite"
	"github.com/go-jet/jet/v2/generator/template"
	sqlite2 "github.com/go-jet/jet/v2/sqlite"
	"github.com/m4rc3l05/media-follower/internal/common"
	"github.com/m4rc3l05/media-follower/internal/store"
	_ "modernc.org/sqlite"
)

func run() (statusCode int) {
	log := common.NewLogger("gen-db-types")
	db, err := store.New()

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing database", slog.Any("err", err))
			statusCode = 1

			return
		}
	}()

	if err != nil {
		log.Error("Error creating database", slog.Any("err", err))

		return 1
	}

	template := template.Default(sqlite2.Dialect).
		UseSchema(func(schemaMetaData metadata.Schema) template.Schema {
			return template.DefaultSchema(schemaMetaData).
				UseSQLBuilder(template.DefaultSQLBuilder().
					UseTable(func(table metadata.Table) template.TableSQLBuilder {
						if table.Name == "schema_migrations" {
							return template.TableSQLBuilder{Skip: true}
						}
						return template.DefaultTableSQLBuilder(table)
					})).
				UseModel(template.DefaultModel().
					UseTable(func(table metadata.Table) template.TableModel {
						if table.Name == "schema_migrations" {
							return template.TableModel{Skip: true}
						}
						return template.DefaultTableModel(table)
					}))
		})

	if err := sqlite.GenerateDB(db.DB, "./.gen/jetdb", template); err != nil {
		log.Error("Error generating database types", slog.Any("err", err))

		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
