package main

import (
	"os"
)

func main() {
	port := os.Getenv("CLOUDOGS_HTTP_PORT")
	if port == "" {
		panic("CLOUDOGS_HTTP_PORT environment variable is required")
	}
	storage := FileStorage{}
	err := storage.Init("dogs.json")
	if err != nil {
		panic(err)
	}
	httpServer := CreateServer(":"+port, &storage)
	err = httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

