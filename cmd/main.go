package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ozontz/app/graph"
	"ozontz/app/storage"
	"syscall"
	"time"

	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
)

func main() {
	storageType := flag.String(
		"storage",
		"inmemory",
		"Select storage type: 'inmemory' or 'postgres'. 'inmemory' by default",
	)
	flag.Parse()

	var store storage.Storage
	switch *storageType {
	case "inmemory":
		log.Println("Initializing in-memory store...")
		store = storage.NewStorageInMemory()
	case "postgres":
		log.Println("Initializing postgres store...")
		db, err := storage.InitPostgresDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		store = storage.NewStoragePostgres(db)
		log.Println("Connectet to db")
	default:
		log.Fatalf("Invalid storage type: %s", *storageType)
	}

	graph.SetStore(store)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graph.QueryType,
		Mutation: graph.MutationType,
	})
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	http.Handle("/query", graph.GraphQLHandler(&schema))

	log.Println("Initializing server...")

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	log.Println("Server stopped")
}
