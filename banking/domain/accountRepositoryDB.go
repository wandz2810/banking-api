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

type AccountRepositoryDb struct {
	client *mongo.Client
	ctx    context.Context
}

func generateUniqueAccountId() string {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

func (d AccountRepositoryDb) SaveAccount(a Account) (*Account, *errs.AppError) {
	customerCollection := d.client.Database("asa").Collection("customer")
	accountCollection := d.client.Database("asa").Collection("account")
	filter := bson.D{{"customer_id", a.CustomerId}}
	var customer Customer
	err := customerCollection.FindOne(d.ctx, filter).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.NewNotFoundError("Customer not found to create account")
		} else {
			logger.Error("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}

	}
	a.AccountId = generateUniqueAccountId()
	result, err2 := accountCollection.InsertOne(d.ctx, a)
	if err2 != nil {
		logger.Error("Error while creating new account: " + err2.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	fmt.Println("Account created successfully with ID:", result.InsertedID)
	return &a, nil

}

func (d AccountRepositoryDb) SaveTransaction(t Transaction) (*Transaction, *errs.AppError) {
	transactionCollection := d.client.Database("asa").Collection("transaction")
	accountCollection := d.client.Database("asa").Collection("account")

	t.TransactionId = generateUniqueAccountId()

	filter := bson.M{"account_id": t.AccountId}
	var update bson.M
	if t.IsWithdrawal() {
		update = bson.M{"$inc": bson.M{"amount": -t.Amount}}
	} else {
		update = bson.M{"$inc": bson.M{"amount": t.Amount}}
	}

	_, err := accountCollection.UpdateOne(d.ctx, filter, update)
	if err != nil {
		logger.Error("Error while saving transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	account, appErr := d.FindByAccountId(t.AccountId)
	if appErr != nil {
		return nil, appErr
	}
	t.Balance = account.Amount
	var result *mongo.InsertOneResult
	result, err = transactionCollection.InsertOne(d.ctx, t)
	if err != nil {
		logger.Error("Error while creating transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	fmt.Println("Transaction created successfully with ID:", result.InsertedID)
	return &t, nil
}

func (d AccountRepositoryDb) FindByCustomerIdAndAccountId(customerId string, accountId string) (*Account, *errs.AppError) {
	collection := d.client.Database("asa").Collection("account")
	filter := bson.D{{"customer_id", customerId}, {"account_id", accountId}}
	var account Account
	err := collection.FindOne(d.ctx, filter).Decode(&account)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.NewNotFoundError("Account not found")
		} else {
			logger.Error("Error while scanning account " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}

	}
	return &account, nil
}

func (d AccountRepositoryDb) FindByAccountId(accountId string) (*Account, *errs.AppError) {
	collection := d.client.Database("asa").Collection("account")
	filter := bson.D{{"account_id", accountId}}
	var account Account
	err := collection.FindOne(d.ctx, filter).Decode(&account)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.NewNotFoundError("Account not found")
		} else {
			logger.Error("Error while scanning account " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}

	}
	return &account, nil
}

func NewAccountRepositoryDb(client *mongo.Client, ctx context.Context) AccountRepositoryDb {

	return AccountRepositoryDb{client, ctx}
}
