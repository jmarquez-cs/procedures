package sqlprocessor

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/afero"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SslMode  string
	Fs       afero.Fs       // Add this line
	DB       *sql.DB        // Add this line
}

func ProcessSQLFile(filename string, config Config) error {
	if !isValidFile(filename) {
		return fmt.Errorf("error: only .sql files are accepted")
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DbName, config.SslMode)
		db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping the database: %v", err)
	}

	sqlContent, err := afero.ReadFile(config.Fs, filename)
	if err != nil {
		return fmt.Errorf("failed to read the .sql file: %v", err)
	}

	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id:   filepath.Base(filename),
				Up:   []string{string(sqlContent)},
				Down: []string{},
			},
		},
	}

	n, err := migrate.Exec(config.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	if n > 0 {
		fmt.Printf("Applied %d migrations from file: %s\n", n, filename)
	} else {
		fmt.Printf("No migrations to apply from file: %s\n", filename)
	}

	return nil
}

func isValidFile(filename string) bool {
	extension := filepath.Ext(filename)
	return strings.ToLower(extension) == ".sql"
}

func ProcessSQLDirectory(directory string, config Config) error {
	files, err := afero.ReadDir(config.Fs, directory)
	if err != nil {
			return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
			if file.IsDir() {
					continue
			}

			filename := filepath.Join(directory, file.Name())
			if !isValidFile(filename) {
					continue
			}

			if err := ProcessSQLFile(filename, config); err != nil {
					return fmt.Errorf("failed to process file '%s': %v", filename, err)
			}
	}

	return nil
}