gen_docs:
	swag init -g ./cmd/main.go --dir ./ --exclude ./ent
swaggo_fmt:
	 go run github.com/swaggo/swag/cmd/swag@v1.8.1 fmt  -g ./cmd/main.go --dir ./ --exclude ./ent

build:
	go build -o ./bin/cmd ./cmd

run_app:
	./bin/cmd

build_image:
	docker build -t toufiq-austcse/go-gin-boilerplate:latest .

docker_up:
	 docker compose up --build -d

docker_down:
	 docker compose down