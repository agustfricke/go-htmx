package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name    string `json:"name"`
}
