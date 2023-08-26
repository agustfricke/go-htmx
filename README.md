##### GO + HTMX
- Simple CRUD de go y htmx con tailwind css, usando una base de datos postres con docker


### Instalar y usar
```bash
git clone https://github.com/agustfricke/go-htmx-crud.git
cd go-htmx-crud
docker run --name some-postgres -e POSTGRES_USER=agust -e POSTGRES_PASSWORD=agust -e POSTGRES_DB=agust -p 5432:5432 -d postgres
go run main.go
```
