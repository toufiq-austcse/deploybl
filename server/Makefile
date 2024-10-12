gen-docs:
	swag init  -g ./cmd/main.go
swaggo-fmt:
	 go run github.com/swaggo/swag/cmd/swag@v1.8.1 fmt  -g ./cmd/main.go --dir ./ --exclude ./ent

build:
	go build -o ./bin/cmd ./cmd

run-app:
	./bin/cmd

tools-install:
	chmod +x ./scripts/install-tools.sh
	./scripts/install-tools.sh
format:
	chmod +x ./scripts/format.sh
	./scripts/format.sh