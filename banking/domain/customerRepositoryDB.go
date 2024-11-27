package domain

import (
	"context"
	"github.com/wandz2810/banking-lib/errs"
	"github.com/wandz2810/banking-lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepositoryDb struct {
	client *mongo.Client
	ctx    context.Context
}

func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	//connect db

	collection := d.client.Database("asa").Collection("customer")
	//query db

	var customers []Customer
	var filter bson.D
	if status == "" {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{Key: "status", Value: status}}
	}
	cur, err := collection.Find(d.ctx, filter)
	if err != nil {
		logger.Error("Error while finding customers " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	for cur.Next(d.ctx) {
		var c Customer
		if err := cur.Decode(&c); err != nil {
			logger.Error("Error while decoding customer " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
		customers = append(customers, c)
	}
	//if err := cur.Err(); err != nil {
	//	log.Println("Cursor error " + err.Error())
	//	return nil, errs.NewUnexpectedError("Unexpected database error")
	//}
	//
	//if len(customers) == 0 {
	//	return customers, errs.NewNotFoundError("No customers found")
	//}

	return customers, nil
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	collection := d.client.Database("asa").Collection("customer")
	filter := bson.D{{"customer_id", id}}
	var customer Customer
	err := collection.FindOne(d.ctx, filter).Decode(&customer)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}

	}
	return &customer, nil
}

func (d CustomerRepositoryDb) UpdateById(c Customer) (*Customer, *errs.AppError) {
	collection := d.client.Database("asa").Collection("customer")
	filter := bson.D{{"customer_id", c.Id}}
	update := bson.D{
		{"$set", bson.D{
			{"name", c.Name},
			{"city", c.City},
			{"zip_code", c.Zipcode},
			{"date_of_birth", c.DateofBirth},
			{"status", c.Status},
		}},
	}
	_, err := collection.UpdateOne(d.ctx, filter, update)
	if err != nil {
		logger.Error("Error while updating customer " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &c, nil
}

func (d CustomerRepositoryDb) DeleteById(customerId string) *errs.AppError {
	// Xóa khách hàng theo customerId
	userCollection := d.client.Database("asa").Collection("user")
	customerCollection := d.client.Database("asa").Collection("customer")
	accountCollection := d.client.Database("asa").Collection("account")
	filter := bson.D{{"customer_id", customerId}}

	// Thực hiện xóa
	_, err := userCollection.DeleteOne(d.ctx, filter)
	if err != nil {
		logger.Error("Error while deleting customer in user collection " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	_, err2 := customerCollection.DeleteOne(d.ctx, filter)
	if err2 != nil {
		logger.Error("Error while deleting customer in customer collection " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	_, err3 := accountCollection.DeleteMany(d.ctx, filter)
	if err3 != nil {
		logger.Error("Error while deleting customer in account collection " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	return nil
}

func NewCustomerRepositoryDb(client *mongo.Client, ctx context.Context) CustomerRepositoryDb {

	return CustomerRepositoryDb{client, ctx}
}
