package main

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"fixtures_from_db/cmd"
)

func main() {
	cmd.RootCmd.AddCommand(cmd.GenerateCmd)
	cmd.RootCmd.AddCommand(cmd.IntrospectCmd)

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("root cmd execute failed %v", err)
	}
}
