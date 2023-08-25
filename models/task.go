package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex" json:"name"`
}
