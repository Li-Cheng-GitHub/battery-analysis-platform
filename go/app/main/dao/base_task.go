package dao

import (
	"battery-analysis-platform/app/main/db"
	"go.mongodb.org/mongo-driver/bson"
)

func creatTask(collectionName string, task interface{}) error {
	return insertMongoCollection(collectionName, task)
}

func deleteTask(collectionName string, id string) error {
	collection := db.Mongo.Collection(collectionName)
	filter := bson.M{"taskId": id}
	ctx := newTimeoutCtx()
	_, err := collection.DeleteOne(ctx, filter)
	return err
}
