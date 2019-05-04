# Cloudogs

### Cloudability dog 

Basic CRUD REST API Server to keep track of Cloudability office dogs

## Setup

[Install Go](https://golang.org/doc/install)

```bash
git clone https://github.com/jperelshteyn/cloudogs.git
cd cloudogs
go get github.com/gorilla/mux
go get github.com/google/uuid
```

## Run

```bash
set CLOUDOGS_HTTP_PORT=8080
go build
./cloudogs
```

## API

### Dog
```json
{
    "id": "7f4fb9fa-0784-4257-bb90-1f5952f6d7ba",
    "name": "Max",
    "owner": "John Smith",
    "notes": "Max is a good boy"
}
```

### Routes

GET /dogs - List all dogs

POST /dogs - Add a new dog**

GET /dogs/:id - Get details for one dog

PUT /dogs/:id - Update details for one dog**

DELETE /dogs/:id - Remove a dog

** must include Dog JSON in the body