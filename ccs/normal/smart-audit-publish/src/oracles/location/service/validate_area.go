package service

import (
	"errors"
	"strconv"
)

// 通过经纬度定义地理位置
type Position struct {
	// 维度
	Lat float64 `json:"lat"`

	// 经度
	Lon float64 `json:"lon"`
}

// 定义用于验证是否在某个范围内的结构
type ValidateArea struct {
	// 验证范围的中心地理位置
	Center Position `json:"center"`

	// 范围半径，以米为单位
	Radius float64 `json:"radius"`
}

// 从输入参数获取范围信息验证规则
func FromStrings(args []string) (*ValidateArea, error) {
	if len(args) < 3 {
		return nil, errors.New("参数不足")
	}

	position, err := PositionFromString(args)
	if err != nil {
		return nil, err
	}

	radius, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, err
	}

	return &ValidateArea{
		Center: *position,
		Radius: radius,
	}, nil
}

// 从输入参数获取位置信息
func PositionFromString(args []string) (*Position, error) {
	if len(args) < 2 {
		return nil, errors.New("参数不足")
	}

	lat, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, err
	}

	return &Position{
		Lat: lat,
		Lon: lon,
	}, nil
}
