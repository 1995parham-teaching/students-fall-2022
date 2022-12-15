package student

import (
	"errors"
	"log"

	"github.com/1995parham-teaching/students/internal/model"
	"gorm.io/gorm"
)

type SQLItem struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func (SQLItem) TableName() string {
	return "students"
}

type SQL struct {
	DB *gorm.DB
}

func NewSQL(db *gorm.DB) Student {
	if err := db.AutoMigrate(new(SQLItem)); err != nil {
		log.Fatal(err)
	}

	return SQL{
		DB: db,
	}
}

func (sql SQL) GetAll() ([]model.Student, error) {
	var items []SQLItem

	if err := sql.DB.Find(&items).Error; err != nil {
		return nil, err
	}

	students := make([]model.Student, 0, len(items))

	for _, item := range items {
		students = append(students, model.Student{
			ID:      item.ID,
			Name:    item.Name,
			Courses: nil,
		})
	}

	return students, nil
}

func (sql SQL) Create(s model.Student) error {
	if err := sql.DB.Create(&SQLItem{
		ID:   s.ID,
		Name: s.Name,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (sql SQL) Get(id string) (model.Student, error) {
	var st SQLItem

	if err := sql.DB.Take(&st, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Student{}, ErrStudentNotFound
		}

		return model.Student{}, err
	}

	return model.Student{
		Name:    st.Name,
		ID:      st.ID,
		Courses: nil,
	}, nil
}
