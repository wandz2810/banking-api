package app

import (
	"banking/domain"
	"banking/service"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wandz2810/banking-lib/logger"
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
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient, ctx)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient, ctx)
	//wiring
	ch := CustomerHandler{service.NewCustomerService(customerRepositoryDb)}
	ah := AccountHandler{service.NewAccountService(accountRepositoryDb)}
	// define routes
	router.HandleFunc("/customer", ch.getAllCustomer).Methods(http.MethodGet).Name("GetAllCustomers")
	router.HandleFunc("/customer/{customer_id:[0-9]+}", ch.getCustomerById).Methods(http.MethodGet).Name("GetCustomer")
	router.HandleFunc("/customer/{customer_id:[0-9]+}/update", ch.UpdateCustomerById).Methods(http.MethodPost).Name("UpdateCustomer")
	router.HandleFunc("/customer/{customer_id:[0-9]+}/delete", ch.DeleteAccountById).Methods(http.MethodGet).Name("DeleteCustomer")
	router.HandleFunc("/customer/{customer_id:[0-9]+}/account", ah.NewAccount).Methods(http.MethodPost).Name("NewAccount")
	router.HandleFunc("/customer/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).Methods(http.MethodPost).Name("NewTransaction")
	am := AuthMiddleware{domain.NewAuthRepository()}
	router.Use(am.authorizationHandler())
	// starting server
	// $env:SERVER_ADDRESS="localhost"
	// $env:SERVER_PORT="8080"
	address := os.Getenv("SER VER_ADDRESS")
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
		logger.Error("Connect DB fail: %v" + err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Error("Connect DB fail: %v" + err.Error())
	}
	fmt.Println("Connect DB successfully")
	return client, ctx
}
