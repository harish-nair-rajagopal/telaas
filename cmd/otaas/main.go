package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	var (
		listen = flag.String("listen", "0.0.0.0:8080", "address of to listen on")
	)
	flag.Parse()
	ctx := context.Background()

	// create new router
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	// GET pipeline

	// GET pipelines

	// POST managed pipeline (create/ update)

	// POST unmanaged pipeline (add)

	// DELETE pipeline

	router.Run(*listen)
}
