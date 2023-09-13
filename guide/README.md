##### GO + HTMX

En el tutorial de hoy vamos a estar haciendo un simple CRUD con Go, htmx y tailwind css,
es nesesario que tengan instalado Go y Docker en su sistema

-   Lo primero que debemos hacer es tener una terminal a nuestra disposicion y vamos a crear un
    nuevo directorio y vamos a meternos dentro con los comandos:

```bash
mkdir ~/go-htmx-crud
cd ~/go-htmx-crud
```

-   Lo siguiente seria iniciar un nuevo modulo con go con el comando <strong>go mod init <url unica>
    </strong>, en mi caso en la url unica voy a poner github.com/agustfricke/go-htmx-crud
    en tu caso seria <strong>github.com/tu-username/go-htmx-crud</strong>

```bash
go mod init github.com/agustfricke/go-htmx-crud
```

-   Perfecto, una vez que tengamos el nuevo modulo de go, podemos comenzar a instalar las dependencias
    que para este proyecto vamos a estar utilizando GORM, el driver de Postgres y dotenv
    GORM es el ORM y dotenv es para manejar las variables de entorono

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get github.com/joho/godotenv
```

-   Ahora podemos crear una nueva base de datos con Postgres y Docker con el comando:

```bash
sudo docker run --name some-postgres -e POSTGRES_USER=agust -e POSTGRES_PASSWORD=agust -e POSTGRES_DB=agust -p 5432:5432 -d postgres
```

-   Creemos un nuevo arvhivo llamado .env, donde vamos a guardar las credenciales de la base de datos

#### ~/go-crud-htmx/.env

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=agust
DB_PASSWORD=agust
DB_NAME=agust
```

-   Exelente, ahora creemos una nueva carpeta para obtener las credenciales de la base de datos

#### ~/go-crud-htmx/config/config.go

```go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
```

-   Creemos los modelo Task que va a estar en la base de datos:

#### ~/go-crud-htmx/models/task.go

```go
package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name    string `json:"name"`
}
```

-   Ahora creemos el siguiente archivo

#### ~/go-crud-htmx/database/database.go

```go
package database

import "gorm.io/gorm"

var DB *gorm.DB
```

-   Creemos el siguiente archivo para conectarnos a la base de datos

#### ~/go-crud-htmx/database/connect.go

```go
package database

import (
	"fmt"
	"strconv"

	"github.com/agustfricke/go-htmx-crud/config"
	"github.com/agustfricke/go-htmx-crud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() {

	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		panic("failed to parse database port")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"),
	config.Config("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.AutoMigrate(&models.Task{})
	fmt.Println("Database Migrated")
}
```

-   Perfecto, ahora debemos crear un nuevo servidor con go

#### ~/go-crud-htmx/main.go

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/agustfricke/go-htmx-crud/database"
)


func main() {

    database.ConnectDB()

    fs := http.FileServer(http.Dir("public"))
    http.Handle("/public/", http.StripPrefix("/public/", fs))

	fmt.Println("Runnning in port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
```

-   Ahora podemos correr nuevo servidor de go para ver que todo funciona correcto

#### ~/go-crud-htmx/

```bash
   ./~/go-crud-htmx/main.go
```

-   Ahora vamos a crear el archvio que va a contener la logica para realizar las operaciones
    en la base de datos, y de retornar el html correspondiente

#### ~/go-crud-htmx/handlers/task.go

```go
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

	    tmpl := template.Must(template.ParseFiles("templates/item.html"))
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

```

-   Perfecto, ahora podemos crear las rutas

#### ~/go-crud-htmx/main.go

```go
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

    // agregar las rutas
	http.HandleFunc("/add/", handlers.CreateTask)
	http.HandleFunc("/delete/", handlers.DeleteTask)
	http.HandleFunc("/edit/form/", handlers.FormEditTask)
	http.HandleFunc("/put", handlers.EditTask)
    http.HandleFunc("/", handlers.GetTasks)
    // fin

	fmt.Println("Runnning in port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
```

-   Perfecto, ahora nos quedaria crear los archivos html, traer htmx y tailwind css

```bash
   mkdir ~/go-crud-htmx/templates
   mkdir ~/go-crud-htmx/public
   touch  ~/go-crud-htmx/templates/index.html
   touch  ~/go-crud-htmx/templates/item.html
   touch  ~/go-crud-htmx/templates/edit.html
```

-   dentro del archivo index.html:

#### ~/go-crud-htmx/templates/index.html

```html
<!doctype html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <script src="public/htmx.js"></script>
        <link rel="stylesheet" href="public/output.css" />

        <title>HTMX & Go</title>
    </head>

    <style>
        body {
            background-color: black;
        }

        .spinner {
            display: none;
        }
        .htmx-request .spinner {
            display: inline;
        }
        .htmx-request.spinner {
            display: inline;
        }
    </style>

    <body class="container mx-auto px-[300px]">
        <form
            id="task-form"
            hx-post="/add/"
            hx-target="#task-list"
            hx-swap="beforeend"
            hx-indicator="#spinner"
            class="mt-11"
        >
            <div class="flex justify-between gap-2">
                <input
                    required
                    type="text"
                    name="name"
                    id="name"
                    class="rounded-lg block text-slate-200 w-full p-2.5 bg-gray-700"
                    placeholder="Name"
                />
                <button
                    class="rounded-lg bg-gray-700 flex justify-between hover:bg-gray-900 py-4 px-8 text-sm capitalize text-white shadow shadow-black/60"
                >
                    <span> Create </span>
                    <div role="status" id="spinner" class="spinner text-white">
                        <svg
                            aria-hidden="true"
                            class="w-5 h-5 ml-2 text-black animate-spin fill-slate-200"
                            viewBox="0 0 100 101"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                                fill="currentColor"
                            />
                            <path
                                d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                                fill="currentFill"
                            />
                        </svg>
                    </div>
                </button>
            </div>
        </form>

        <div class="row mt-4 g-4">
            <div class="">
                <ul class="" id="task-list">
                    {{ range . }}
                    <li class="text-slate-200" id="item-{{ .ID }}">
                        <div class="flex justify-center">
                            <div
                                class="w-[300px] mb-2 border border-gray-200 rounded-lg shadow bg-gray-800 border-gray-700"
                            >
                                <div class="flex flex-col items-center py-2">
                                    <span
                                        class="font-poppis text-gray-500 dark:text-gray-400"
                                    >
                                        {{ .Name }} - {{ .ID }}
                                    </span>
                                    <div class="flex jusitfy-between mt-2">
                                        <button
                                            hx-get="/edit/form?name={{ .Name }}&ID={{ .ID }}"
                                            hx-target="#item-{{ .ID }}"
                                            class="mr-2 inline-flex items-center px-4 py-2 text-sm font-medium text-center text-gray-900 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-gray-200 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-700 dark:focus:ring-gray-700"
                                        >
                                            Edit
                                        </button>

                                        <div class="flex justify-between">
                                            <button
                                                hx-delete="/delete/{{ .ID }}"
                                                hx-swap="delete"
                                                hx-target="#item-{{ .ID }}"
                                                hx-indicator="#spinner-delete-{{ .ID }}"
                                                class="inline-flex items-center px-4 py-2 text-sm font-medium text-center text-gray-900 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-gray-200 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-700 dark:focus:ring-gray-700"
                                            >
                                                Delete

                                                <div
                                                    role="status"
                                                    id="spinner-delete-{{ .ID }}"
                                                    class="spinner"
                                                >
                                                    <svg
                                                        aria-hidden="true"
                                                        class="w-5 h-5 ml-2 text-black animate-spin fill-slate-200"
                                                        viewBox="0 0 100 101"
                                                        fill="none"
                                                        xmlns="http://www.w3.org/2000/svg"
                                                    >
                                                        <path
                                                            d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                                                            fill="currentColor"
                                                        />
                                                        <path
                                                            d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                                                            fill="currentFill"
                                                        />
                                                    </svg>
                                                </div>
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </li>
                    {{ end }}
                </ul>
            </div>
        </div>
    </body>
</html>
```

-   dentro del archivo index.html:

#### ~/go-crud-htmx/templates/item.html

```html
<li class="text-slate-200" id="item-{{ .Task.ID }}">
    <div class="flex justify-center">
        <div
            class="w-[300px] mb-2 border border-gray-200 rounded-lg shadow bg-gray-800 border-gray-700"
        >
            <div class="flex flex-col items-center py-2">
                <span class="font-poppis text-gray-500 dark:text-gray-400">
                    {{ .Task.Name }} - {{ .Task.ID }}
                </span>
                <div class="flex jusitfy-between mt-2">
                    <button
                        hx-get="/edit/form?name={{ .Task.Name }}&ID={{ .Task.ID }}"
                        hx-target="#item-{{ .Task.ID }}"
                        class="mr-2 inline-flex items-center px-4 py-2 text-sm font-medium text-center text-gray-900 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-gray-200 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-700 dark:focus:ring-gray-700"
                    >
                        Edit
                    </button>

                    <div class="flex justify-between">
                        <button
                            hx-delete="/delete/{{ .Task.ID }}"
                            hx-swap="delete"
                            hx-target="#item-{{ .Task.ID }}"
                            hx-indicator="#spinner-delete-{{ .Task.ID }}"
                            class="inline-flex items-center px-4 py-2 text-sm font-medium text-center text-gray-900 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-gray-200 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-700 dark:focus:ring-gray-700"
                        >
                            Delete

                            <div
                                role="status"
                                id="spinner-delete-{{ .Task.ID }}"
                                class="spinner"
                            >
                                <svg
                                    aria-hidden="true"
                                    class="w-5 h-5 ml-2 text-black animate-spin fill-slate-200"
                                    viewBox="0 0 100 101"
                                    fill="none"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path
                                        d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                                        fill="currentColor"
                                    />
                                    <path
                                        d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                                        fill="currentFill"
                                    />
                                </svg>
                            </div>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</li>
```

-   Dentro del archivo edit.html:

#### ~/go-crud-htmx/templates/edit.html

```html
<form hx-put="/put" hx-indicator="#spinner-{{ .ID }}">
    <ul class="" id="task-list">
        <li class="text-slate-200">
            <div class="flex justify-center">
                <div
                    class="w-[300px] mb-2 border border-gray-200 rounded-lg shadow bg-gray-800 border-gray-700"
                >
                    <div class="flex flex-col items-center py-2">
                        <span
                            class="flex m-1 justify-between font-poppis text-gray-500 dark:text-gray-400"
                        >
                            <input
                                required
                                value="{{ .Name }}"
                                type="text"
                                name="name"
                                class="rounded-lg block text-slate-200 w-full p-2.5 bg-gray-700"
                                placeholder="Name"
                            />
                            <input
                                type="text"
                                hidden
                                name="ID"
                                value="{{ .ID }}"
                            />
                            <button
                                type="submit"
                                class="ml-2 inline-flex items-center px-4 py-2 text-sm font-medium text-center text-gray-900 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-gray-200 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-700 dark:focus:ring-gray-700"
                            >
                                Save

                                <div
                                    role="status"
                                    class="spinner ml-2 text-white"
                                    id="spinner-{{ .ID }}"
                                >
                                    <svg
                                        aria-hidden="true"
                                        class="w-5 h-5 ml-2 text-black animate-spin fill-slate-200"
                                        viewBox="0 0 100 101"
                                        fill="none"
                                        xmlns="http://www.w3.org/2000/svg"
                                    >
                                        <path
                                            d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                                            fill="currentColor"
                                        />
                                        <path
                                            d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                                            fill="currentFill"
                                        />
                                    </svg>
                                </div>
                            </button>
                        </span>
                    </div>
                </div>
            </div>
        </li>
    </ul>
</form>
```

Ahora podemos descargar htmx en nuestro proyecto y compliar el codigo de tailwind css

```bash
    cd ~/go-crud-htmx/public
    wget  https://unpkg.com/htmx.org@1.9.5/dist/htmx.min.js
    cd ..
    npm install -D tailwindcss
    npx tailwindcss init
```

-   Dentro del tailwind.config.js pon lo siguiente:

#### ~/go-crud-htmx/tailwind.config.js

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./templates/*.html"],
    theme: {
        extend: {},
    },
    plugins: [],
};
```

-   Crea el archivo input.css dentro de public

#### ~/go-crud-htmx/public/input.css

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

-   Ahora por ultimo podemos correr el comando para hacer un build de tailwind css:

```css
npx tailwindcss -i ./public/input.css -o ./public/output.css
```

-   Perfecto ahora podemos correr nuestro servidor de go y a probar la app!

```css
go run main.go
```
