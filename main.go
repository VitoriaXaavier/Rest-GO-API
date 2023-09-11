package main

import (
	"log"
	"os"
)

func main() {
	args := Args{
		conn: "postgres://postgres:@localhost:5432/postgres?sslmode=disable",
		port: ":8000", 
	}

	if conn := os.Getenv("DB_CONN"); conn != "" {
		args.conn = conn
	}
	if port := os.Getenv("PORT"); port != "" {
		args.port = ":" + port
	}

	if err := Run(args); err != nil {
		log.Panicln(err)
	}
}