package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"fixtures_from_db/database"
	"fixtures_from_db/fixture"
	"fixtures_from_db/pkg/graph"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage: go run main.go <connection string>")
		return
	}
	fmt.Println(os.Args[1])
	ctx := context.Background()

	db, err := database.NewDB(os.Args[1])
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

	order, err := g.TopSort()

	for i := range order {
		tableColumnsInfo := columnInfoByTableName[order[i]]
		fmt.Println(fixture.ToSQL(order[i], tableColumnsInfo))
	}
}
