package course

import (
	"log"

	"github.com/1995parham-teaching/students/internal/model"
	"gorm.io/gorm"
)

type SQLItem struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func (SQLItem) TableName() string {
	return "courses"
}

type SQL struct {
	DB *gorm.DB
}

func NewSQL(db *gorm.DB) Course {
	if err := db.AutoMigrate(new(SQLItem)); err != nil {
		log.Fatal(err)
	}

	return SQL{
		DB: db,
	}
}

func (sql SQL) GetAll() ([]model.Course, error) {
	return nil, nil
}

func (sql SQL) Create(s model.Course) error {
	return nil
}

func (sql SQL) Get(id string) (model.Course, error) {
	return model.Course{}, nil
}
