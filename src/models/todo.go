package models

import (
	"gorm.io/gorm"
	_ "time"
)

type Todo struct {
	gorm.Model
	Name        string
	Description string
	IsDone      int
}
