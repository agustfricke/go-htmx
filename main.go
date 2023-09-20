package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/agustfricke/go-htmx-crud/database"
	"github.com/agustfricke/go-htmx-crud/handlers"
)


func main() {

    database.ConnectDB()

    fs := http.FileServer(http.Dir("public"))
    http.Handle("/public/", http.StripPrefix("/public/", fs))

	  http.HandleFunc("/add/", handlers.CreateTask)
	  http.HandleFunc("/delete/", handlers.DeleteTask)
	  http.HandleFunc("/edit/form/", handlers.FormEditTask)
	  http.HandleFunc("/put", handlers.EditTask)
    http.HandleFunc("/", handlers.GetTasks)

	  fmt.Println("Runnning in port 8000")
	  log.Fatal(http.ListenAndServe(":8000", nil))
}
