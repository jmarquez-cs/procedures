package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ScooterHelmet/procedures/pkg/sqlprocessor"
)

func printHelp() {
	helpText := `
Usage: program [options] <filename.sql>

Options:
  --help                Display this help message.

Environment Variables:
  PG_HOST               Set the PostgreSQL database connection. Options: host.docker.internal (macOs & Windows). Default: localhost.
  SSL_MODE              Set the SSL mode for the database connection. Options: disable, allow, prefer, require, verify-ca, verify-full. Default: disable.
`
	fmt.Println(helpText)
}

func main() {
	helpFlag := flag.Bool("help", false, "Display help message")
	flag.Parse()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Missing filename.sql argument.")
		printHelp()
		os.Exit(1)
	}

	filename := args[0]

	pgHost := os.Getenv("PG_HOST")
	if pgHost == "" {
		pgHost = "localhost"
		fmt.Println("Warning: PG_HOST not set. Using default: localhost")
	}

	sslMode := os.Getenv("SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
		fmt.Println("Warning: SSL mode not set. Using default: disable")
	}

	config := sqlprocessor.Config{
		Host:     pgHost,
		Port:     5432,
		User:     "postgres",
		Password: "",
		DbName:   "postgres",
		SslMode:  sslMode,
	}

	err := sqlprocessor.ProcessSQLFile(filename, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Processed file: %s\n", filename)
}