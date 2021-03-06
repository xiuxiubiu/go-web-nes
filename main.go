package main

import (
	"go-web-nes/router"
	"log"
	"net/http"
)

func main() {

	r := router.InitRouter()

	srv := &http.Server{
		Addr:    ":8181",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen: %\n", err)
	}
}
