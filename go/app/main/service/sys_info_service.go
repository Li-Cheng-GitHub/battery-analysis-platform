package service

import (
	"battery-analysis-platform/app/main/model"
	"battery-analysis-platform/pkg/jd"
)

type SysInfoShowService struct {
}

func (s *SysInfoShowService) Do() (*jd.Response, error) {
	data, err := model.NewSysInfo()
	if err != nil {
		return nil, err
	}
	return jd.Build(jd.SUCCESS, "", data), nil
}
