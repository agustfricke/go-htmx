package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/agustfricke/go-htmx-crud/database"
	"github.com/agustfricke/go-htmx-crud/models"
)


func GetTasks(w http.ResponseWriter, r *http.Request) {
        db := database.DB 
	    var tasks []models.Task
	    db.Find(&tasks)

	    tmpl := template.Must(template.ParseFiles("templates/index.html"))

	    err := tmpl.Execute(w, tasks)
	    if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	    }
	}

func CreateTask(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		name := r.PostFormValue("name")
        
        var task models.Task 
        if name != "" {
            db := database.DB
            task = models.Task{Name: name} 
            db.Create(&task)
        }
	    data := struct {Task models.Task}{Task: task,}

	    tmpl := template.Must(template.ParseFiles("templates/index.html"))
	    err :=  tmpl.Execute(w, data)
	    if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	    }

	}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
    time.Sleep(1 * time.Second)
    
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) != 3 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    ID := parts[2]
    
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
	    err :=  tmpl.Execute(w, data)
	    if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		    return
	    }
    }
    
func EditTask(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)

        name := r.PostFormValue("name")
        ID := r.PostFormValue("ID") 
        db := database.DB

        var task models.Task
        if err := db.First(&task, ID).Error; err != nil {
            fmt.Printf("NOp")
        }

        task.Name = name 

        if err := db.Save(&task).Error; err != nil {
            fmt.Printf("Error al guardar la tarea: %s\n", err)
        }

	    data := struct {Task models.Task}{Task: task}

	    tmpl := template.Must(template.ParseFiles("templates/item.html"))
	    err :=  tmpl.Execute(w, data)
	    if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		    return
	    }

	}

