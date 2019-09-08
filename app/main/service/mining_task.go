package service

import (
	"battery-anlysis-platform/app/main/dao"
	"battery-anlysis-platform/app/main/model"
	"battery-anlysis-platform/pkg/checker"
	"battery-anlysis-platform/pkg/jtime"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	collectionNameTaskList = "mining_tasks"
	timeout                = time.Second
)

type MiningCreateTaskService struct {
	TaskName     string `json:"taskName" binding:"required"`
	DataComeFrom string `json:"dataComeFrom" binding:"required"`
	StartDate    string `json:"startDate" binding:"required"`
	EndDate      string `json:"endDate" binding:"required"`
	AllData      bool   `json:"allData"` // bool 型不能 required，因为 false 会被判空
}

func (s *MiningCreateTaskService) CreateTask() (*model.MiningTask, error) {
	if _, ok := model.MiningSupportTaskSet[s.TaskName]; !ok {
		return nil, errors.New("参数 TaskName 不合法")
	}
	if _, ok := model.BatteryMysqlNameToTable[s.DataComeFrom]; !ok {
		return nil, errors.New("参数 dataComeFrom 不合法")
	}
	var requestParams string
	if s.AllData {
		requestParams = "所有数据"
	} else {
		if !checker.ReDatetime.MatchString(s.StartDate) {
			return nil, errors.New("参数 startDate 不合法")
		}
		if !checker.ReDatetime.MatchString(s.EndDate) {
			return nil, errors.New("参数 EndDate 不合法")
		}
		requestParams = s.StartDate + " - " + s.EndDate
	}
	data := &model.MiningTask{
		Id:            "123",
		TaskName:      s.TaskName,
		DataComeFrom:  s.DataComeFrom,
		RequestParams: requestParams,
		CreateTime:    jtime.NowStr(),
		TaskStatus:    "执行中",
		Comment:       "",
	}
	return data, nil
}

func GetTaskList() ([]model.MiningTask, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := dao.MongoDB.Collection(collectionNameTaskList)
	filter := bson.M{}                  // 过滤记录
	projection := bson.M{"data": false} // 过滤字段
	cur, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	var records []model.MiningTask
	for cur.Next(ctx) {
		result := model.MiningTask{}
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		records = append(records, result)
	}
	_ = cur.Close(ctx)
	return records, nil
}

func GetTask(id string) (bson.A, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := dao.MongoDB.Collection(collectionNameTaskList)
	filter := bson.M{"_id": id}
	projection := bson.M{"_id": false, "data": true} // 注意 _id 默认会返回，需要手动过滤
	// 注意 bson.E 不能用来映射 mongo 中的 map，
	// 要么使用 bson.D，采用 []bson.E 代表一个字典，其中 bson.E 是 struct，有 key 和 value 字段，
	// 此时，映射出来的子字典也都是 bson.D 类型，
	// 而映射出来的列表是 bson.A 类型，
	// bson.D 在 JSON 序列化时会在最外层加上 []，所以需要序列化的结果不要用，而采用 bson.M；
	// 要么使用 bson.M，采用 map[string]interface{} 代表一个字典，
	// 此时，映射出来的子字典也都是 bson.M 类型，
	// 而映射出来的列表也是 bson.A 类型，
	// 这种方法 JSON 序列化时符合直觉，推荐使用；
	// 若要代表一个列表，类似 Python 中 list，不限定类型，使用 bson.A，即 []interface{}。
	var result bson.M
	err := collection.FindOne(ctx, filter, options.FindOne().
		SetProjection(projection)).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result["data"].(bson.A), nil
}

func DeleteTask(id string) (int64, error) {
	// TODO 终止正在执行的后台任务
	// ...

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := dao.MongoDB.Collection(collectionNameTaskList)
	filter := bson.M{"_id": id}
	ret, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return ret.DeletedCount, nil
}