package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"fixtures_from_db/config"
	"fixtures_from_db/database"
	"fixtures_from_db/pkg/graph"
)

var dbConnectionString string

var IntrospectCmd = &cobra.Command{
	Use: "introspect",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		db, err := database.NewDB(dbConnectionString)
		if err != nil {
			log.Fatalf("new db failed: %v", err)
		}
		defer func(db *database.DB) {
			err := db.Close()
			if err != nil {
				log.Printf("db close failed: %v", err)
			}
		}(db)

		allTables, err := db.GetAllTableNames(ctx)
		if err != nil {
			log.Fatalf("get all table names failed: %v", err)
		}

		g := graph.NewGraph[string]()
		columnInfoByTableName := make(map[string][]database.ColumnInfo)
		for i := range allTables {
			info, err := db.GetTableColumnsInfo(ctx, allTables[i])
			if err != nil {
				log.Fatalf("get table columns info failed: %v", err)
			}

			columnInfoByTableName[allTables[i]] = info
			for j := range info {
				if info[j].IsForeignKey && info[j].ForeignTable != nil {
					g.AddEdge(allTables[i], *info[j].ForeignTable)
				} else {
					g.AddNode(allTables[i])
				}
			}
		}

		tablesOrder, err := g.TopSort()
		if err != nil {
			log.Fatalf("top sort failed: %v", err)
		}

		var cfg config.Config
		cfg.Tables = make(config.TablesDescription)
		for _, tableName := range tablesOrder {
			cfg.TablesOrder = append(cfg.TablesOrder, config.TableName(tableName))
			tableDescription := make(config.TableDescription)
			for _, columnInfo := range columnInfoByTableName[tableName] {
				columnName := config.ColumnName(columnInfo.ColumnName)
				columnDescription := config.ColumnDescription{
					Type: columnInfo.DataType,
					RuleSet: config.RuleSet{
						Nullable: !columnInfo.IsPrimaryKey || !columnInfo.IsNotNull,
						Unique:   columnInfo.IsPrimaryKey,
						Values:   columnInfo.Values,
					},
				}
				if columnInfo.IsForeignKey && columnInfo.ForeignTable != nil && columnInfo.ForeignColumnName != nil {
					columnDescription.Reference = &config.ColumnDescriptionReference{
						Table:  config.TableName(*columnInfo.ForeignTable),
						Column: config.ColumnName(*columnInfo.ForeignColumnName),
					}
				}
				tableDescription[columnName] = columnDescription
			}
			cfg.Tables[config.TableName(tableName)] = tableDescription
		}

		configContent, err := config.SaveToBytes(&cfg)
		if err != nil {
			log.Fatalf("save to bytes failed: %v", err)
		}

		fmt.Println(string(configContent))

		return nil
	},
}

func init() {
	IntrospectCmd.Flags().StringVarP(&dbConnectionString, "db-connection-string", "c", "", "Connection string for database")
}
