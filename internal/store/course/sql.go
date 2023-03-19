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
	var items []SQLItem

	if err := sql.DB.Model(new(SQLItem)).Preload("Courses").Find(&items).Error; err != nil {
		return nil, err
	}

	courses := make([]model.Course, 0)

	for _, item := range items {
		courses = append(courses, model.Course{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return courses, nil
}

func (sql SQL) Create(s model.Course) error {
	if err := sql.DB.Create(&SQLItem{
		ID:   s.ID,
		Name: s.Name,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (sql SQL) Get(id string) (model.Course, error) {
	var c SQLItem

	if err := sql.DB.First(&c, id).Error; err != nil {
		return model.Course{}, err
	}

	return model.Course{
		ID:   c.ID,
		Name: c.Name,
	}, nil
}
