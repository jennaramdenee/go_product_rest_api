// Main entry point for the application
package main

import (
  "os"
)

func main() {
  SetEnvironmentVariables()

  a := App{}
  a.Initialize(
    os.Getenv("APP_DB_USERNAME"),
    os.Getenv("APP_DB_PASSWORD"),
    os.Getenv("APP_DB_NAME"))
}
