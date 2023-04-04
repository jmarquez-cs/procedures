package sqlprocessor

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"

	"github.com/doug-martin/goqu/v9"
	"github.com/rocketlaunchr/dataframe-go/pandas"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SslMode  string
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

	sqlContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read the .sql file: %v", err)
	}

	queries := strings.Split(string(sqlContent), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		// Convert MySQL dialect to PostgreSQL dialect using pandas
		convertedQuery, err := pandas.ConvertDialect(query, pandas.MySQL, pandas.PostgreSQL)
		if err != nil {
			return fmt.Errorf("failed to convert MySQL query to PostgreSQL query: %v", err)
		}

		// Use goqu to execute the converted query
		dialect := goqu.Dialect("postgres")
		exec := dialect.Exec(db)

		_, err = exec.Exec(convertedQuery)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return nil
}

func isValidFile(filename string) bool {
	extension := filepath.Ext(filename)
	return strings.ToLower(extension) == ".sql"
}
