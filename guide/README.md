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

-   Ahora creemos otra carpeta para conectarnos a la base de datos

#### ~/go-crud-htmx/database/connect.go

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
