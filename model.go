// Represent the product
package main

import (
  "database/sql"
  "errors"
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
  _, err := db.QueryRow("INSERT INTO products(name, price) VALUES($1, $2) RETURNING id", p.Name, p.Price).Scan(&p.ID))
  if err != nil {
    return err
  }
  return nil
}

// Functions that deal with many products
func (p *product) getProducts(db *sql.DB, start, count int) ([]product, error) {
  return nil, errors.New("Not implemented")
}
