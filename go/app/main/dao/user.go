package dao

import (
	"battery-analysis-platform/app/main/consts"
	"battery-analysis-platform/app/main/db"
	"battery-analysis-platform/app/main/model"
	"battery-analysis-platform/pkg/security"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(name, password, comment string) (*model.User, error) {
	user := model.NewUser(name, password, comment)
	err := insertMongoCollection(consts.MongoCollectionUser, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUser(name string) (*model.User, error) {
	var user model.User
	collection := db.Mongo.Collection(consts.MongoCollectionUser)
	filter := bson.M{"name": name}
	projection := bson.M{"_id": false} // 注意 _id 默认会返回，需要手动过滤
	ctx, _ := context.WithTimeout(context.Background(), consts.MongoCtxTimeout)
	err := collection.FindOne(ctx, filter,
		options.FindOne().SetProjection(projection)).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func ListCommonUser() ([]model.User, error) {
	collection := db.Mongo.Collection(consts.MongoCollectionUser)
	filter := bson.M{"type": bson.M{"$ne": consts.UserTypeSuperUser}} // 过滤记录
	projection := bson.M{"_id": false}                                // 过滤字段
	ctx := context.TODO()
	cur, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	// 为了使其找不到时返回空列表，而不是 nil
	users := make([]model.User, 0)
	for cur.Next(ctx) {
		result := model.User{}
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		users = append(users, result)
	}
	_ = cur.Close(ctx)
	return users, nil
}

func SaveUserLoginTimeAndCount(user *model.User) error {
	collection := db.Mongo.Collection(consts.MongoCollectionUser)
	filter := bson.M{"name": user.Name} // 过滤记录
	ctx, _ := context.WithTimeout(context.Background(), consts.MongoCtxTimeout)
	update := bson.M{"$set": bson.M{
		"lastLoginTime": user.LastLoginTime,
		"loginCount":    user.LoginCount,
	}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func SaveUserChange(user *model.User) error {
	collection := db.Mongo.Collection(consts.MongoCollectionUser)
	filter := bson.M{"name": user.Name} // 过滤记录
	ctx, _ := context.WithTimeout(context.Background(), consts.MongoCtxTimeout)
	update := bson.M{"$set": bson.M{
		"comment": user.Comment,
		"status":  user.Status,
	}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func ChangeUserPassword(userName, password string) error {
	collection := db.Mongo.Collection(consts.MongoCollectionUser)
	filter := bson.M{"name": userName} // 过滤记录
	ctx, _ := context.WithTimeout(context.Background(), consts.MongoCtxTimeout)
	s, err := security.GeneratePasswordHash(password)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{
		"password": s,
	}}
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// ---------------------------cache---------------------------

func SaveUserToCache(user *model.User) error {
	// 存储 JSON 序列化的数据
	jd, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return db.Redis.Set(consts.RedisKeyPrefix+user.Name, jd, consts.RedisKeyExpiration).Err()
}

func DeleteUserFromCache(name string) error {
	return db.Redis.Del(consts.RedisKeyPrefix + name).Err()
}

func GetUserFromCache(name string) (*model.User, error) {
	val, err := db.Redis.Get(consts.RedisKeyPrefix + name).Bytes()
	if err != nil {
		return nil, err
	}
	user := model.User{}
	err = json.Unmarshal(val, &user)
	if err != nil {
		return nil, err
	}
	// 刷新 key 的过期时间
	db.Redis.Expire(consts.RedisKeyPrefix+name, consts.RedisKeyExpiration)
	return &user, nil
}
