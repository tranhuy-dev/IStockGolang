package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/tranhuy-dev/IStockGolang/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	dbName = "IStock"
)

var Client *mongo.Client

func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	Client = client
}

func InsertCustomer(req models.CustomerReq) interface{} {
	newCustomer := models.Customer{
		FirstName: req.FirstName,
		LastName: req.LastName,
		Phone: req.Phone,
		Address: req.Address,
		Age: req.Age,
		Status:1}
	customerCollection := Client.Database("IStock").Collection("customer")
	insertQuery, errorQueryInsert := customerCollection.InsertOne(context.TODO(), newCustomer)
	if errorQueryInsert != nil {
		log.Fatal(errorQueryInsert)
	}
	return insertQuery.InsertedID
}

func RetrieveAllCustomer() interface{} {
	var customer []*models.Customer
	customerCollection := Client.Database("IStock").Collection("customer")
	findOptions := options.Find()
	findOptions.SetLimit(2)
	cur, err := customerCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Customer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		customer = append(customer, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return customer
}

func UpdateCustomer(req models.CustomerReq , id string) interface{}{
	customerCollection := Client.Database("IStock").Collection("customer")
	objID , _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id",objID}}

	updateBody := bson.D{
		{"$set" , bson.D{
			{"first_name" , req.FirstName},
			{"last_name",req.LastName},
			{"age",req.Age},
			{"phone",req.Phone},
			{"address",req.Address},
		}},
	}
	updateResult, err := customerCollection.UpdateOne(context.TODO() , filter , updateBody)
	if err != nil {
		log.Fatal(err)
	}

	return updateResult
}

func DeleteCustomer(id string) interface{} {
	customerCollection := Client.Database("IStock").Collection("customer")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id" , objID}}
	updateBody := bson.D{
		{"$set", bson.D{
			{"status",0},
		}},
	}
	deleteResult, err := customerCollection.UpdateOne(context.TODO() , filter , updateBody)
	if err != nil {
		log.Fatal(err)
	}
	return deleteResult
}