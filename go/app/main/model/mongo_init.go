package model

import (
	"battery-analysis-platform/app/main/db"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	mongoCtxTimeout = time.Second * 5
)

const (
	mongoCollectionUser          = "user"
	mongoCollectionYuTongVehicle = "yutong_vehicle"
	mongoCollectionBeiQiVehicle  = "beiqi_vehicle"
	mongoCollectionMiningTask    = "mining_task"
	mongoCollectionDlTask        = "deeplearning_task"
)

// 确保创建 mongo 索引
func createMongoCollectionIdx(name string, model mongo.IndexModel) error {
	collection := db.Mongo.Collection(name)
	ctx, _ := context.WithTimeout(context.Background(), mongoCtxTimeout)
	_, err := collection.Indexes().CreateOne(
		ctx,
		model,
	)
	return err
}

// 在 collection 中插入一条记录
func insertMongoCollection(collectionName string, item interface{}) error {
	collection := db.Mongo.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), mongoCtxTimeout)
	_, err := collection.InsertOne(ctx, item)
	return err
}

func init() {
	// user
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	if err := createMongoCollectionIdx(mongoCollectionUser, indexModel); err != nil {
		panic(err)
	}
	indexModel = mongo.IndexModel{
		Keys: bson.M{
			"type": 1,
		},
		Options: options.Index().SetUnique(false),
	}
	if err := createMongoCollectionIdx(mongoCollectionUser, indexModel); err != nil {
		panic(err)
	}

	// yutong_vehicle
	indexModel = mongo.IndexModel{
		Keys: bson.M{
			"时间": 1,
		},
		Options: options.Index().SetUnique(false),
	}
	if err := createMongoCollectionIdx(mongoCollectionYuTongVehicle, indexModel); err != nil {
		panic(err)
	}
	indexModel = mongo.IndexModel{
		Keys: bson.M{
			"状态号": 1,
		},
		Options: options.Index().SetUnique(false),
	}
	if err := createMongoCollectionIdx(mongoCollectionYuTongVehicle, indexModel); err != nil {
		panic(err)
	}

	// beiqi_vehicle
	indexModel = mongo.IndexModel{
		Keys: bson.M{
			"时间": 1,
		},
		Options: options.Index().SetUnique(false),
	}
	if err := createMongoCollectionIdx(mongoCollectionBeiQiVehicle, indexModel); err != nil {
		panic(err)
	}
	indexModel = mongo.IndexModel{
		Keys: bson.M{
			"状态号": 1,
		},
		Options: options.Index().SetUnique(false),
	}
	if err := createMongoCollectionIdx(mongoCollectionBeiQiVehicle, indexModel); err != nil {
		panic(err)
	}

	// task
	indexModel = mongo.IndexModel{
		Keys: bson.M{
			"taskId": 1,
		},
		Options: options.Index().SetUnique(false),
	}
	if err := createMongoCollectionIdx(mongoCollectionMiningTask, indexModel); err != nil {
		panic(err)
	}
	if err := createMongoCollectionIdx(mongoCollectionDlTask, indexModel); err != nil {
		panic(err)
	}
}
