package lambdaClone

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"src/mongodb"
	"strconv"
	"strings"
	"time"
)

const lambdaCollection string = "lambda"
const lambdaRuntimeCollection string = "lambda.runtime"
const lambdaDefaultCollection string = "lambda.default"

func validateInsertLambda(ctx *fiber.Ctx) (*Lambda, error) {
	var err error
	var lambda Lambda
	if err = json.Unmarshal((*ctx).Body(), &lambda); err != nil {
		return nil, err
	}
	if strings.TrimSpace(lambda.Title) == "" ||
		strings.TrimSpace(lambda.Runtime) == "" ||
		strings.TrimSpace(lambda.Version) == "" {
		return nil, errors.New("lambda is not a valid format")
	}
	if strings.TrimSpace(lambda.Id) != "" {
		return nil, errors.New("lambda id param should not be set")
	}
	return &lambda, nil
}

func insertLambda(lambda *Lambda) (*primitive.ObjectID, error) {
	// query
	cursor, err := mongodb.Database.
		Collection(lambdaCollection).
		InsertOne(context.Background(), bson.M{
			"title":       lambda.Title,
			"runtime":     lambda.Runtime,
			"version":     lambda.Version,
			"reg_date":    time.Now(),
			"update_date": time.Now(),
			"disabled":    false,
		})
	if err != nil {
		return nil, err
	}
	// return
	primitiveId := cursor.InsertedID.(primitive.ObjectID)
	return &primitiveId, nil
}

func validateGetLambda(ctx *fiber.Ctx) (*primitive.ObjectID, error) {
	var err error
	var id string
	if id = (*ctx).Query("id", ""); id == "" {
		return nil, errors.New("invalid parameter id")
	}
	var primitiveId primitive.ObjectID
	if primitiveId, err = primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	}
	return &primitiveId, nil
}

func getLambda(objectID *primitive.ObjectID) (*Lambda, error) {
	var err error
	// query
	cursor := mongodb.Database.
		Collection(lambdaCollection).
		FindOne(context.Background(), bson.M{"_id": *objectID})
	// decode
	var result Lambda
	if err = cursor.Decode(&result); err != nil {
		return nil, err
	}
	// return
	return &result, nil
}

func getDefaultLambda() (*Lambda, error) {
	var err error
	// query
	cursor := mongodb.Database.
		Collection(lambdaDefaultCollection).
		FindOne(context.Background(), bson.M{"disabled": false})
	// decode
	var result Lambda
	if err = cursor.Decode(&result); err != nil {
		return nil, err
	}
	// return
	return &result, nil
}

func validateGetLambdas(ctx *fiber.Ctx) (*options.FindOptions, error) {
	var err error

	var limit int
	if (*ctx).Queries()["limit"] == "" {
		limit = 30
	} else if limit, err = strconv.Atoi((*ctx).Queries()["limit"]); err != nil {
		return nil, err
	}

	var page int
	if (*ctx).Queries()["page"] == "" {
		page = 1
	} else if page, err = strconv.Atoi((*ctx).Queries()["page"]); err != nil {
		return nil, err
	}
	if page <= 0 {
		return nil, errors.New("parameter page must be > 0")
	}

	var sort string
	if (*ctx).Queries()["sort"] == "" {
		sort = "reg_date"
	} else {
		sort = (*ctx).Queries()["sort"]
	}

	var asc int
	if (*ctx).Queries()["asc"] == "" || (*ctx).Queries()["asc"] == "false" {
		asc = -1
	} else if (*ctx).Queries()["asc"] == "true" {
		asc = 1
	} else {
		return nil, errors.New("parameter asc must be true or false or not set")
	}

	option := options.FindOptions{}
	option.SetLimit(int64(limit))
	option.SetSkip(int64(page-1) * int64(limit))
	option.SetSort(bson.M{sort: asc})
	return &option, nil
}

func getLambdas(option *options.FindOptions) (*[]Lambda, error) {
	var err error
	// query
	cursor, err := mongodb.Database.
		Collection(lambdaCollection).
		Find(context.Background(), bson.M{"disabled": false}, option)
	if err != nil {
		return nil, err
	}
	// decode
	var result []Lambda
	for cursor.Next(context.Background()) {
		var lambda Lambda
		if err = cursor.Decode(&lambda); err != nil {
			return nil, err
		}
		result = append(result, lambda)
	}
	// no results
	if cursor == nil || result == nil {
		result = []Lambda{}
	}
	// return
	return &result, nil
}

func getLambdasCnt(limit *int64) (int64, int64, error) {
	cnt, err := mongodb.Database.Collection(lambdaCollection).
		CountDocuments(context.Background(), bson.M{"disabled": false})
	if err != nil {
		return 0, 0, err
	}
	var pages = cnt / *limit
	if cnt%*limit != 0 {
		pages++
	}
	return cnt, pages, nil
}

func validateUpdateLambda(ctx *fiber.Ctx) (*Lambda, error) {
	var err error
	var lambda Lambda
	if err = json.Unmarshal((*ctx).Body(), &lambda); err != nil {
		return nil, err
	}
	if strings.TrimSpace(lambda.Title) == "" ||
		strings.TrimSpace(lambda.Runtime) == "" ||
		strings.TrimSpace(lambda.Version) == "" ||
		strings.TrimSpace(lambda.Code) == "" {
		return nil, errors.New("lambda is not a valid format")
	}
	if strings.TrimSpace(lambda.Id) == "" {
		return nil, errors.New("lambda id param should be set")
	}
	return &lambda, nil
}

func updateLambda(lambda *Lambda) error {
	id, err := primitive.ObjectIDFromHex(lambda.Id)
	if err != nil {
		return err
	}
	// query
	result, err := mongodb.Database.
		Collection(lambdaCollection).
		UpdateOne(context.Background(),
			bson.M{"_id": id},
			bson.M{"$set": bson.M{
				"title":       lambda.Title,
				"runtime":     lambda.Runtime,
				"version":     lambda.Version,
				"update_date": time.Now(),
			}},
		)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("document with id not found")
	}
	return nil
}

func validateDeleteLambda(ctx *fiber.Ctx) (*primitive.ObjectID, error) {
	var err error
	var id string
	if id = (*ctx).Query("id", ""); id == "" {
		return nil, errors.New("invalid parameter id")
	}
	var primitiveId primitive.ObjectID
	if primitiveId, err = primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	}
	return &primitiveId, nil
}

func deleteLambda(id *primitive.ObjectID) error {
	// query
	result, err := mongodb.Database.
		Collection(lambdaCollection).
		UpdateOne(context.Background(),
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"disabled": true}})
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("document with id not found")
	}
	return nil
}

func getRuntimes() (*[]Runtime, error) {
	// query
	cursor, err := mongodb.Database.
		Collection(lambdaRuntimeCollection).
		Find(context.Background(),
			bson.M{
				"disabled": false,
			})
	// decode
	var result []Runtime
	for cursor.Next(context.Background()) {
		var runtime Runtime
		if err = cursor.Decode(&runtime); err != nil {
			return nil, err
		}
		result = append(result, runtime)
	}
	// no results
	if cursor == nil || result == nil {
		result = []Runtime{}
	}
	// return
	return &result, nil
}

func setupDefaultData() error {
	// insert default query
	lambda := Lambda{Runtime: "golang", Version: "1.22", Disabled: false, RegDate: time.Now(), UpdateDate: time.Now()}
	_, err := mongodb.Database.
		Collection(lambdaDefaultCollection).
		InsertOne(context.Background(), lambda)
	if err != nil {
		return err
	}
	// insert runtimes query
	runtimes := []any{
		Runtime{Runtime: "golang", Version: "1.22", Disabled: false, RegDate: time.Now(), UpdateDate: time.Now()},
		Runtime{Runtime: "node", Version: "20", Disabled: false, RegDate: time.Now(), UpdateDate: time.Now()},
		Runtime{Runtime: "python", Version: "3.12", Disabled: false, RegDate: time.Now(), UpdateDate: time.Now()},
	}
	_, err = mongodb.Database.
		Collection(lambdaRuntimeCollection).
		InsertMany(context.Background(), runtimes)
	if err != nil {
		return err
	}
	return nil
}
