package main

import (
	"fmt"
	"runtime/debug"

	"github.com/toufiq-austcse/deployit/internal/app"
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
