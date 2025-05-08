package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/maxischmaxi/ljtime-api/customer/v1/customerv1connect"
	"github.com/maxischmaxi/ljtime-api/project/v1/projectv1connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	customerServer := &CustomerServer{
		MongoClient: client,
		DBName: "lj-time",
	}

	projectServer := &ProjectServer{
		MongoClient: client,
		DBName: "lj-time",
	}

	mux := http.NewServeMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},	
		AllowCredentials: true,
	}).Handler(h2c.NewHandler(mux, &http2.Server{}))

	customerPath, customerHandler := customerv1connect.NewCustomerServiceHandler(customerServer)
	projectPath, projectHandler := projectv1connect.NewProjectServiceHandler(projectServer)

	mux.Handle(customerPath, customerHandler)
	mux.Handle(projectPath, projectHandler)

	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
