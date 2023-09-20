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

    if err := db.Find(&tasks).Error; err != nil {
        http.Error(w, "Error getting tasks from database", http.StatusInternalServerError)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    if err := tmpl.Execute(w, tasks); err != nil {
        http.Error(w, "Render error", http.StatusInternalServerError)
        return
    }
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
    time.Sleep(2 * time.Second)

    name := r.PostFormValue("name")

    if name == "" {
        http.Error(w, "Can't create task without a name", http.StatusBadRequest)
        return
    }

    db := database.DB
    task := models.Task{Name: name}

    if err := db.Create(&task).Error; err != nil {
        http.Error(w, "Error creating task in database", http.StatusInternalServerError)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/item.html"))
    if err := tmpl.Execute(w, task); err != nil {
        http.Error(w, "Render error", http.StatusInternalServerError)
        return
    }
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
    time.Sleep(2 * time.Second)

    ID := r.URL.Query().Get("ID")

    if ID == "" {
        http.Error(w, "ID not found", http.StatusBadRequest)
        return
    }

    db := database.DB
    var task models.Task

    if err := db.First(&task, ID).Error; err != nil {
            http.Error(w, "Task not found", http.StatusNotFound)
            return
    }

    if err := db.Delete(&task).Error; err != nil {
        http.Error(w, "Error deleting task from database", http.StatusInternalServerError)
        return
    }

}

func FormEditTask(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    ID := r.URL.Query().Get("ID")

    if ID == "" || name == "" {
        http.Error(w, "ID or Name not found", http.StatusBadRequest)
        return
    }

    data := struct{ ID string; Name string }{ID: ID, Name: name}

    tmpl := template.Must(template.ParseFiles("templates/edit.html"))
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Render error", http.StatusInternalServerError)
        return
    }
}
    
func EditTask(w http.ResponseWriter, r *http.Request) {
    time.Sleep(2 * time.Second)

    name := r.PostFormValue("name")
    ID := r.URL.Query().Get("ID")

    if ID == "" || name == "" {
        http.Error(w, "ID or Name not found", http.StatusBadRequest)
        return
    }

    db := database.DB

    var task models.Task
    if err := db.First(&task, ID).Error; err != nil {
            http.Error(w, "Task not found", http.StatusNotFound)
            return
    }

    task.Name = name
    if err := db.Save(&task).Error; err != nil {
        http.Error(w, "Error saving task in database", http.StatusInternalServerError)
        return
    }

    tmpl := template.Must(template.ParseFiles("templates/item.html"))
    if err := tmpl.Execute(w, task); err != nil {
        http.Error(w, "Render error", http.StatusInternalServerError)
        return
    }
}
