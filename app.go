// Holds our application
package main

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strconv"

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
  a.initializeRoutes()
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
func (a *App) Run(addr string) {
  log.Fatal(http.ListenAndServe(":8000", a.Router))
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
    return
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

  p := product{ID:id}
  if err := p.deleteProduct(a.DB); err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  respondWithJSON(w, http.StatusOK, map[string]string{ "result": "success"})
}

// Handler method for many products
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
  // get count and start parameters from querystring
  // r.FormValue returns the first value for the named component of the query
  count, _ := strconv.Atoi(r.FormValue("count"))
  start, _ := strconv.Atoi(r.FormValue("start"))

  if count > 10 || count < 1 {
    count = 10
  }

  if start < 0 {
    start = 0
  }

  products, err := getProducts(a.DB, start, count)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  respondWithJSON(w, http.StatusOK, products)
}
