package model

import (
	"gorm.io/gorm"
)
type AccessLog struct {
	gorm.Model
	Method	string
	Root 	string
	Path 	string
	ClientIp string
	AccessTime string
}