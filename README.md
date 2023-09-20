##### GO + HTMX

-   Simple CRUD de go y htmx usando una base de datos postgres con docker.

### Instalar y usar

```bash
git clone https://github.com/agustfricke/go-htmx-crud.git
cd go-htmx-crud
mv .env.example .env
docker run --name postgres_db -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -e POSTGRES_DB=super_db -p 5432:5432 -d postgres
go run main.go
```

## Dale una estrella ⭐
