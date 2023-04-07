package sqlprocessor

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	MemFS = afero.NewMemMapFs()
)

func setupMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock
}

func TestIsValidFile(t *testing.T) {
	validExtensions := []string{".sql", ".SQL", ".Sql"}
	invalidExtensions := []string{".txt", ".docx", ".csv"}

	for _, ext := range validExtensions {
		assert.True(t, isValidFile("file"+ext), fmt.Sprintf("Expected %s to be valid", ext))
	}

	for _, ext := range invalidExtensions {
		assert.False(t, isValidFile("file"+ext), fmt.Sprintf("Expected %s to be invalid", ext))
	}
}

func TestProcessSQLDirectory(t *testing.T) {
	MockDB, Mock := setupMockDB()
	config := Config{
		Fs: MemFS,
		DB: MockDB,
	}

	Mock.ExpectExec("create table if not exists").WillReturnResult(sqlmock.NewResult(0, 0))

	testPath := "test"
	_ = MemFS.Mkdir(testPath, 0755)
	// trunk-ignore(golangci-lint/errcheck)
	afero.WriteFile(MemFS, filepath.Join(testPath, "file1.sql"), []byte("CREATE TABLE test_table (id serial PRIMARY KEY, name VARCHAR(50) NOT NULL);"), 0644)

	err := ProcessSQLDirectory(testPath, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = Mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProcessSQLFile(t *testing.T) {
	MockDB, Mock := setupMockDB()
	config := Config{
		Fs: MemFS,
		DB: MockDB,
	}

	testPath := "test"
	// trunk-ignore(golangci-lint/errcheck)
	afero.WriteFile(MemFS, filepath.Join(testPath, "file2.sql"), []byte("CREATE TABLE another_table (id serial PRIMARY KEY, description TEXT NOT NULL);"), 0644)

	Mock.ExpectExec("create table if not exists").WillReturnResult(sqlmock.NewResult(0, 0))
	Mock.ExpectExec("CREATE TABLE another_table").WillReturnResult(sqlmock.NewResult(0, 0))

	err := ProcessSQLFile(filepath.Join(testPath, "file2.sql"), config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = Mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProcessSQLFile_InvalidFile(t *testing.T) {
	config := Config{}

	err := ProcessSQLFile("invalid.txt", config)
	assert.Equal(t, errors.New("error: only .sql files are accepted"), err)
}

func TestProcessSQLFile_ErrorReadingFile(t *testing.T) {
	filePath := "nonexistent.sql"
	fs := afero.NewMemMapFs()

	config := Config{
		Fs: fs,
	}

	err := ProcessSQLFile(filePath, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read the .sql file")
}

func TestProcessSQLFile_ErrorExecutingQuery(t *testing.T) {
	sqlContent := "INVALID SQL QUERY;"
	filePath := "file.sql"
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, filePath, []byte(sqlContent), 0644)
	require.NoError(t, err)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	mock.ExpectExec("INVALID SQL QUERY").WillReturnError(errors.New("syntax error"))

	config := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DbName:   "postgres",
		SslMode:  "disable",
		Fs:       fs,
		DB:       db,
	}

	err = ProcessSQLFile(filePath, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query")
}