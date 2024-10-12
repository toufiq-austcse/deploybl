# https://github.com/segmentio/golines
# # will do `gofmt` internally
golines -m 120 -w --ignore-generated .

go run github.com/swaggo/swag/cmd/swag@v1.8.1 fmt  -g ./cmd/main.go --dir ./
# https://github.com/mvdan/gofumpt
# will do `gofmt` internally
gofumpt -l -w .