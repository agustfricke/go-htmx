package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/agustfricke/go-htmx-crud/database"
	"github.com/agustfricke/go-htmx-crud/models"
)

type Film struct {
	Title    string
	Director string
}

func main() {

    database.ConnectDB()

	fmt.Println("Go app...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
        db := database.DB 
	    var tasks []models.Task
	    db.Find(&tasks)

	    tmpl := template.Must(template.ParseFiles("index.html"))

	    err := tmpl.Execute(w, tasks)
	    if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	    }
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		name := r.PostFormValue("name")

        if name != "" {
            db := database.DB
            task := models.Task{Name: name} 
            db.Create(&task)
        }
		htmlStr := fmt.Sprintf("<li>%s</li>", name)
		tmpl, _ := template.New("t").Parse(htmlStr)
	    tmpl.Execute(w, nil)

	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
        ID := r.URL.Query().Get("ID")
	    db := database.DB

	    var task models.Task
	    db.First(&task, ID)
	    db.Delete(&task)
	}

    h4 := func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Query().Get("name")
        ID := r.URL.Query().Get("ID")
        htmlContent := `
        <form hx-put="/put">
            <input type="text" name="name" value="` + name + `">
            <input type="text" hidden name="ID" value="` + ID + `">
            <button type="submit">Update</button>
        </form>`

        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(htmlContent))
    }

	h5 := func(w http.ResponseWriter, r *http.Request) {

name := r.PostFormValue("name")
ID := r.PostFormValue("ID") 
db := database.DB

var task models.Task
if err := db.First(&task, ID).Error; err != nil {
    fmt.Printf("NOp")
}

task.Name = name // Asignar el valor del nombre directamente

if err := db.Save(&task).Error; err != nil {
    fmt.Printf("Error al guardar la tarea: %s\n", err)
}

htmlStr := fmt.Sprintf("<li>%s</li>", task.Name) // Usar task.Name en lugar de name
tmpl, _ := template.New("t").Parse(htmlStr)
tmpl.Execute(w, nil)
	}


	// define handlers
	http.HandleFunc("/", h1)
	http.HandleFunc("/add/", h2)
	http.HandleFunc("/delete/", h3)
	http.HandleFunc("/edit/form/", h4)
	http.HandleFunc("/put/", h5)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
