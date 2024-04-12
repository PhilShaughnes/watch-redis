# https://cheatography.com/linux-china/cheat-sheets/justfile/

set dotenv-load := true

default:
	@just --list

# run it!
r ORG:
	go run ./main.go -host=$REDIS_HOST -port=$REDIS_PORT -key=seed:{{ ORG }}

# watch and run a go file
watch PATH:
	ls {{PATH}}/* | entr -c go run {{PATH}}/*.go

# watch and run a go file
wtest PATH:
	ls {{PATH}}/* | entr -c go test {{PATH}}/*.go
