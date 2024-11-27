package domain

import (
	"context"
	"fmt"
	"github.com/wandz2810/banking-lib/errs"
	"github.com/wandz2810/banking-lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
)

type AuthRepositoty interface {
	FindBy(string, string) (*Login, *errs.AppError)
	CreateUser(register Register) (*Register, *errs.AppError)
}

type AuthRepositoryDB struct {
	client *mongo.Client
	ctx    context.Context
}

func (d AuthRepositoryDB) FindBy(username, password string) (*Login, *errs.AppError) {
	userCollection := d.client.Database(
		"asa").Collection("user")
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{
				"username": username, // Replace with actual username
				"password": password, // Replace with actual password
			}},
		},
		{
			{"$lookup", bson.M{
				"from":         "account",     // Join with the 'account' collection
				"localField":   "customer_id", // Field from the 'user' collection
				"foreignField": "customer_id", // Field from the 'account' collection
				"as":           "accounts",    // Name of the resulting array
			}},
		},
		{
			{"$unwind", bson.M{
				"path":                       "$accounts", // Deconstruct the accounts array
				"preserveNullAndEmptyArrays": true,        // Keep users with no accounts
			}},
		},
		{
			{"$group", bson.M{
				"_id":             "$customer_id", // Group by customer_id
				"username":        bson.M{"$first": "$username"},
				"role":            bson.M{"$first": "$role"},
				"account_numbers": bson.M{"$push": "$accounts.account_id"}, // Collect account IDs into an array
			}},
		},
		{
			{"$project", bson.M{
				"_id":         0, // Exclude _id field from the output
				"username":    1,
				"customer_id": "$_id", // Rename _id to customer_id
				"role":        1,
				"account_numbers": bson.M{ // Concatenate account numbers into a single string
					"$reduce": bson.M{
						"input":        "$account_numbers",
						"initialValue": "",
						"in": bson.M{
							"$cond": []interface{}{
								bson.M{"$eq": []interface{}{"$$value", ""}},
								"$$this",
								bson.M{"$concat": []interface{}{"$$value", ", ", "$$this"}},
							},
						},
					},
				},
			}},
		},
	}

	cursor, err := userCollection.Aggregate(d.ctx, pipeline)
	if err != nil {
		logger.Error("error aggregating users")
		return nil, errs.NewUnexpectedError("Error occurred while retrieving user data")
	}
	var result Login
	if cursor.Next(d.ctx) {
		if err := cursor.Decode(&result); err != nil {
			logger.Error("error decoding result")
			return nil, errs.NewUnexpectedError("Error decoding user data")
		}
	} else {
		logger.Error("no documents found for the provided credentials")
		return nil, errs.NewNotFoundError("User not found")
	}
	return &result, nil
}

func generateUniqueId() string {
	return fmt.Sprintf("%05d", rand.Intn(100000))
}

func (d AuthRepositoryDB) CreateUser(r Register) (*Register, *errs.AppError) {
	userCollection := d.client.Database("asa").Collection("user")
	customerCollection := d.client.Database("asa").Collection("customer")
	customerid := generateUniqueId()
	type User struct {
		Username   string `bson:"username"`
		Password   string `bson:"password"`
		Role       string `bson:"role"`
		CustomerId string `bson:"customer_id"`
	}
	user := User{
		Username:   r.Username,
		Password:   r.Password,
		Role:       r.Role,
		CustomerId: customerid,
	}
	result, err := userCollection.InsertOne(d.ctx, user)
	if err != nil {
		logger.Error("Error while creating new user: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	type Customer struct {
		CustomerId  string `bson:"customer_id"`
		Name        string `bson:"name"`
		City        string `bson:"city"`
		Zipcode     string `bson:"zip_code"`
		DateofBirth string `bson:"date_of_birth"`
		Status      string `bson:"status"`
	}
	customer := Customer{
		CustomerId:  customerid,
		Name:        r.Name,
		City:        r.City,
		Zipcode:     r.Zipcode,
		DateofBirth: r.DateofBirth,
		Status:      r.Status,
	}
	_, err1 := customerCollection.InsertOne(d.ctx, customer)
	if err1 != nil {
		logger.Error("Error while creating new customer: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	fmt.Println("Register successfully with ID:", result.InsertedID)
	r.CustomerId = customerid
	return &r, nil
}

func NewAuthRepositoryDB(client *mongo.Client, ctx context.Context) AuthRepositoryDB {
	return AuthRepositoryDB{client, ctx}
}
