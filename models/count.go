package models

import (
	"gorm.io/gorm"
)

type Count struct {
	gorm.Model
	Number int `json:"count" param:"count"`
}