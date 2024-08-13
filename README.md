## GO + HTMX
Simple CRUD with Go and Htmx

### Requirements
- Go
- Docker

### How to use
The commands are for Mac and Linux.

- Clone the repository
```bash
git clone https://github.com/agustfricke/go-htmx-crud.git ~/
cd ~/go-htmx-crud
```

- Setup the environment, run the database and the application
```bash
mv .example.env .env
docker run --name postgres_db -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -e POSTGRES_DB=super_db -p 5432:5432 -d postgres
go run main.go
```

## Give it a ⭐
