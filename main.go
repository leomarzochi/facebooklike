package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leomarzochi/facebooklike/cmd/config"
	"github.com/leomarzochi/facebooklike/cmd/router"
)

func main() {
	config.Load()

	r := router.Generate()

	fmt.Println("Server listening on port: " + config.WebPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.WebPort), r))
}
