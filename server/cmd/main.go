package main

import (
	"fmt"
	"runtime/debug"

	"github.com/toufiq-austcse/deployit/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Println("error in running application ", err)
		debug.PrintStack()
		return
	}
}
