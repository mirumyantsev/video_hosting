package main

import "github.com/mirumyantsev/video_hosting/pkg/config"

type NonCatVideo struct {
	Id             int
	CodeMP         string
	StartDatetime  string
	DurationRecord int
}

type Repo struct {
	cfg *config.Config
}
