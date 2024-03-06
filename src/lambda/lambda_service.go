package lambda

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/mongodb"
	"strings"
)

const collection string = "lambda"

func validateInsertLambda(lambda *Lambda) error {
	if strings.TrimSpace(lambda.Title) == "" ||
		strings.TrimSpace(lambda.Runtime) == "" ||
		strings.TrimSpace(lambda.Version) == "" ||
		strings.TrimSpace(lambda.Code) == "" {
		return errors.New("lambda is not a valid format")
	}
	if strings.TrimSpace(lambda.Id) != "" ||
		strings.TrimSpace(lambda.Endpoint) != "" {
		return errors.New("lambda id, endpoint param should not be set")
	}
	return nil
}

func insertLambda(lambda *Lambda) (primitive.ObjectID, error) {
	// query
	cursor, err := mongodb.Database.
		Collection(collection).
		InsertOne(context.Background(), *lambda)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	// return
	return cursor.InsertedID.(primitive.ObjectID), nil
}

func getLambda(objectID primitive.ObjectID) (Lambda, error) {
	var result Lambda
	var err error
	// query
	cursor := mongodb.Database.
		Collection(collection).
		FindOne(context.Background(), bson.M{"_id": objectID})
	// decode
	if err = cursor.Decode(&result); err != nil {
		return Lambda{}, err
	}
	// return
	return result, nil
}
