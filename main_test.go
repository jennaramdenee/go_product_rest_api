package main_test

import (
  "os"
  "testing"
  "log"
  "net/http"
  "."
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
  id SERIAL,
  name TEXT NOT NULL,
  price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
  CONSTRAINT products_pkey PRIMARY KEY (id)
)`

var a main.App

// Ensures that database is correctly set up and cleared before running tests
func TestMain(m *testing.M) {
  main.SetEnvironmentVariables()

  a = main.App{}
  a.Initialize(
    os.Getenv("TEST_DB_USERNAME"),
    os.Getenv("TEST_DB_PASSWORD"),
    os.Getenv("TEST_DB_NAME"))

  ensureTableExists()

  code := m.Run()

  clearTable()

  os.Exit(code)
}

func ensureTableExists() {
  if _, err := a.DB.Exec(tableCreationQuery); err != nil {
    log.Fatal(err)
  }
}

func clearTable() {
  // Remember that a has a DB property for the database, as per struct
  a.DB.Exec("DELETE FROM products")
  a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
  clearTable()

  req, _ := http.NewRequest("GET", "/products", nil)

}
