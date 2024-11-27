package app

import (
	"banking_auth/domain"
	"banking_auth/service"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

func sanityCheck() {
	if os.Getenv("SERVER_ADDRESS") == "" ||
		os.Getenv("SERVER_PORT") == "" ||
		os.Getenv("MONGODB_URI") == "" {
		log.Fatal("Environment variables not set")
	}

}

func Start() {

	sanityCheck()

	//create multiplexer
	//mux := http.NewServeMux()
	router := mux.NewRouter()

	dbClient, ctx := getDbClient()
	authRepositoryDb := domain.NewAuthRepositoryDB(dbClient, ctx)
	//wiring
	ah := AuthHandler{service.NewAuthService(authRepositoryDb, domain.GetRolePermissions())}
	// define routes
	router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", ah.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	// starting server
	// $env:SERVER_ADDRESS="localhost"
	// $env:SERVER_PORT="8181"
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func getDbClient() (*mongo.Client, context.Context) {
	var ctx = context.TODO()
	mongoURI := os.Getenv("MONGODB_URI")

	// $env:MONGODB_URI="mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Connect DB fail: %v" + err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Connect DB fail: %v" + err.Error())
	}
	fmt.Println("Connect DB successfully")
	return client, ctx
}
