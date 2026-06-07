package commands

import (
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/sqlite"
	"github.com/go-jet/jet/v2/generator/template"
	sqlite2 "github.com/go-jet/jet/v2/sqlite"
	"github.com/m4rc3l05/media-follower/internal/storage"
	_ "modernc.org/sqlite"
)

func GenDbTypes(db *storage.Db) error {
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

	if err := sqlite.GenerateDB(db.DB, ".gen/go-jet", template); err != nil {
		return err
	}

	return nil
}
