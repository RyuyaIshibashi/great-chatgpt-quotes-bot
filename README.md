## start development env
```bash
docker compose build
docker compose run bot bash
```

## run go main
```bash
# create go.mod file
go mod init great-chatgpt-quotes-bot

# go.sum file
go mod tidy

go run main.go
```