// Represent the product
package main

import (
  "database/sql"
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
func getProducts(db *sql.DB, start, count int) ([]product, error) {
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
