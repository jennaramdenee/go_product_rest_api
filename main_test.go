package main_test

import (
  "os"
  "testing"
  "log"
  "net/http"
  "net/http/httptest"
  "encoding/json"
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
  // Remember that 'a' has a DB property for the database, as per struct
  a.DB.Exec("DELETE FROM products")
  a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
  clearTable()

  req, _ := http.NewRequest("GET", "/products", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  if body := response.Body.String(); body != "[]" {
    t.Errorf("Expected an empty array. Got %s", body)
  }
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
  rr := httptest.NewRecorder()
  a.Router.ServeHTTP(rr, req)
  return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
  if actual != expected {
    t.Errorf("Expected response code %d. Got %d\n", expected, actual)
  }
}

func TestGetNonExistentProduct(t *testing.T) {
  clearTable()

  req, _ := http.NewRequest("GET", "/product/11", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusNotFound, response.Code)

  var m map[string]string
  // Parse JSON data into format of m; stores key value pairs into the map
  json.Unmarshal(response.Body.Bytes(), &m)
  if m["error"] != "Product not found" {
    t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got %s.", m["error"])
  }
}

// func TestCreateProduct(t *testing.T) {
//   clearTable()
//
//   payload := []byte(`{ "name": "test product", "price": 11.22 }`)
//
//
//
// }
