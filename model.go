// Represent the product
package main

import (
  "database/sql"
  "strconv"
)

type product struct {
  ID    int     `json:"id"`
  Name  string  `json:"name"`
  Price float64 `json:"price"`
}

// Functions that deal with a single product
func (p *product) getProduct(db *sql.DB) error {
  // Scan converts columns read from the database into Go types, and saves to p.Name, p.Price
  return db.QueryRow("SELECT name, price FROM products WHERE id=$1", p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(db *sql.DB) error {
  _, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)
  return err
}

func (p *product) deleteProduct(db *sql.DB) error {
  _, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)
  return err
}

func (p *product) createProduct(db *sql.DB) error {
  err := db.QueryRow("INSERT INTO products(name, price) VALUES($1, $2) RETURNING id", p.Name, p.Price).Scan(&p.ID)
  if err != nil {
    return err
  }
  return nil
}

// Functions that deal with many products
func (p *product) getProducts(db *sql.DB, start, count int) ([]product, error) {
  // Fetch records from the products table, limiting number of records by count, and skipping the number of records defined by start
  rows, err := db.Query("SELECT id, name, price FROM products LIMIT $1 OFFSET $2", count, start)
  if err != nil {
    return nil, err
  }

  defer rows.Close()

  products := []product{}

  // Next prepares the next result row for reading with the Scan method
  for rows.Next() {
    var p product
    if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
      return nil, err
    }
    // Create slice of product objects
    products = append(products, p)
  }
  return products, nil
}

// Handler method for a single product
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
  // Get id of the product from requested URL
  vars := mux.Vars(r)
  // Atoi is string to int
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid product ID")
    return
  }

  p := product{ ID: id }
  // Uses the getProduct method defined above
  if err := p.getProduct(a.DB); err != nil {
    switch err {
    case sql.ErrNoRows:
      respondWithError(w, http.StatusNotFound, "Product not found")
    default:
      respondWithError(w, http.StatusInternalServerError, err.Error())
    }
    return
  }
  respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request){
  var p product
  // Decoding a stream of data - assumes that the input is in JSON format
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&p); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }
  // Remember to close stream
  defer r.Body.Close()

  if err := p.createProduct(a.DB); err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  }

  respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid product ID")
    return
  }

  var p product
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&p); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }
  defer r.Body.Close()
  p.ID = id

  if err := p.updateProduct(a.DB); err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "ID not valid")
    return
  }

  p := product{ID: id}
  err := p.deleteProduct(a.DB); err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  respondWithJSON(w, http.StatusOK, map[string]string{ "result": "success"})
}
// Handler method for many products
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
  // get count and start parameters from querystring
  // r.FormValue returns the first value for the named component of the query
  count = strconv.Atoi(r.FormValue("count"))
  start = strconv.Atoi(r.FormValue("start"))

  if count > 10 || count < 1 {
    count = 10
  }

  if start < 0 {
    start = 0
  }

  products, err := p.getProduct(a.DB, start, count)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  respondWithJSON(w, http.StatusOK, products)
}
