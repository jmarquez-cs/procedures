package sqlprocessor

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidFile(t *testing.T) {
	// Test valid file extension
	valid := isValidFile("test.sql")
	assert.True(t, valid, "Expected .sql file to be valid")

	// Test invalid file extension
	invalid := isValidFile("test.txt")
	assert.False(t, invalid, "Expected non .sql file to be invalid")
}

func TestProcessSQLFile(t *testing.T) {
	// Create temporary test .sql file
	tmpFile, err := ioutil.TempFile("", "test-*.sql")
	assert.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tmpFile.Name())

	// Write sample SQL content to temporary file
	sqlContent := "CREATE TABLE test (id serial PRIMARY KEY, name VARCHAR(50));\nDROP TABLE test;"
	_, err = tmpFile.Write([]byte(sqlContent))
	assert.NoError(t, err, "Failed to write SQL content to temp file")

	// Test configuration
	config := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DbName:   "postgres",
		SslMode:  "disable",
	}

	// Test valid processing of the SQL file
	err = ProcessSQLFile(tmpFile.Name(), config)
	assert.NoError(t, err, "Expected processing of SQL file to succeed")

	// Test invalid processing of a non .sql file
	err = ProcessSQLFile("test.txt", config)
	assert.Error(t, err, "Expected processing of non .sql file to fail")
}
