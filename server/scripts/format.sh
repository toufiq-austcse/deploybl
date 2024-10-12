# https://github.com/segmentio/golines
# # will do `gofmt` internally
golines -m 120 -w --ignore-generated .

# https://github.com/mvdan/gofumpt
# will do `gofmt` internally
gofumpt -l -w .