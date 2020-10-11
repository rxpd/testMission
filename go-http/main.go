package main

import (
	"avito/router"
	"avito/utils"
	"fmt"
	"net/http"
)

func Run(address string, port string) {
	fullAddress := address + ":" + port
	fmt.Printf("\nRun server on %s \n", fullAddress)
	err := http.ListenAndServe(fullAddress, router.GetRoutes())
	utils.LogFatalIf(err)
}

func main() {
	Run("127.0.0.1", "8080")
}

