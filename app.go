// Holds our application
package main

import (
  "fmt"
  "log"
  "database/sql"
  "github.com/gorilla/mux"
  _ "github.com/lib/pq"
)

// References to the router and database that the application uses
type App struct {
  Router  *mux.Router
  DB      *sql.DB
}

// Takes details required to run the database
// Establish connection with database and initialize router
func (a *App) Initialize(user, password, dbname string) {
  connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

  var err error
  // Open with driver name and data source strings
  a.DB, err = sql.Open("postgres", connectionString)
  if err != nil {
    log.Fatal(err)
  }

  a.Router = mux.NewRouter()
}

// Define routes that will use handlers in model.go
func (a *App) initializeRoutes(){
  a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
  a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
  a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
  a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
  a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

// Starts the application
func (a *App) Run(addr string) { }
