package main

import (
	"fmt"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/app"
	"runtime/debug"
)

const configPath = ".env"

func main() {
	err := app.Run(configPath)
	if err != nil {
		fmt.Println("error in running application ", err)
		debug.PrintStack()
		return
	}

}
