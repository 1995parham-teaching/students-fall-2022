package course

import (
	"context"
	"log"

	"gorm.io/gorm"

	"github.com/1995parham-teaching/students/internal/model"
)

type SQLItem struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func (SQLItem) TableName() string {
	return "courses"
}

type SQL struct {
	conn gorm.Interface[SQLItem]
}

func NewSQL(db *gorm.DB) Course {
	err := db.AutoMigrate(new(SQLItem))
	if err != nil {
		log.Fatal(err)
	}

	return SQL{
		conn: gorm.G[SQLItem](db),
	}
}

func (sql SQL) GetAll(ctx context.Context) ([]model.Course, error) {
	items, err := sql.conn.Find(ctx)
	if err != nil {
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

func (sql SQL) Create(ctx context.Context, s model.Course) error {
	return sql.conn.Create(ctx, &SQLItem{
		ID:   s.ID,
		Name: s.Name,
	})
}

func (sql SQL) Get(ctx context.Context, id string) (model.Course, error) {
	c, err := sql.conn.Where("id = ?", id).First(ctx)
	if err != nil {
		return model.Course{}, err
	}

	return model.Course{
		ID:   c.ID,
		Name: c.Name,
	}, nil
}
