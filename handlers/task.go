package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/agustfricke/go-htmx-crud/database"
	"github.com/agustfricke/go-htmx-crud/models"
)


func GetTasks(w http.ResponseWriter, r *http.Request) {
  db := database.DB 
	var tasks []models.Task
	db.Find(&tasks) 

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
  time.Sleep(2 * time.Second)
	name := r.PostFormValue("name")
  db := database.DB
  task := models.Task{Name: name} 
  db.Create(&task)

	tmpl := template.Must(template.ParseFiles("templates/item.html"))
	tmpl.Execute(w, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
    time.Sleep(2 * time.Second)
    ID := r.URL.Query().Get("ID")
    db := database.DB
    var task models.Task
    db.First(&task, ID)
    db.Delete(&task)
}

func FormEditTask(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")
  ID := r.URL.Query().Get("ID")
	data := struct {ID string; Name string}{ID: ID, Name: name}

	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
  tmpl.Execute(w, data)
}
    
func EditTask(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)

  name := r.PostFormValue("name")
  ID := r.URL.Query().Get("ID")
  db := database.DB
  var task models.Task
  db.First(&task, ID)

  task.Name = name 
  db.Save(&task)

	tmpl := template.Must(template.ParseFiles("templates/item.html"))
	tmpl.Execute(w, task)
}

