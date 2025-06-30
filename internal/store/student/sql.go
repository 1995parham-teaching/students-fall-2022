package student

import (
	"context"
	"errors"
	"log"

	"github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/store/course"
	"gorm.io/gorm"
)

type SQLItem struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	Courses []course.SQLItem `gorm:"many2many:students_courses"`
}

func (SQLItem) TableName() string {
	return "students"
}

type SQL struct {
	conn gorm.Interface[SQLItem]
	db   *gorm.DB
}

func NewSQL(db *gorm.DB) Student {
	err := db.AutoMigrate(new(SQLItem))
	if err != nil {
		log.Fatal(err)
	}

	return SQL{
		conn: gorm.G[SQLItem](db),
		db:   db,
	}
}

func (sql SQL) GetAll(ctx context.Context) ([]model.Student, error) {
	items, err := sql.conn.Preload("Courses", nil).Find(ctx)
	if err != nil {
		return nil, err
	}

	students := make([]model.Student, 0, len(items))

	for _, item := range items {
		courses := make([]model.Course, 0, len(item.Courses))

		for _, item := range item.Courses {
			courses = append(courses, model.Course{
				Name: item.Name,
				ID:   item.ID,
			})
		}

		students = append(students, model.Student{
			ID:      item.ID,
			Name:    item.Name,
			Courses: courses,
		})
	}

	return students, nil
}

func (sql SQL) Create(ctx context.Context, s model.Student) error {
	return sql.conn.Create(ctx, &SQLItem{
		ID:      s.ID,
		Name:    s.Name,
		Courses: nil,
	})
}

func (sql SQL) Register(ctx context.Context, sid string, cid string) error {
	c, err := gorm.G[course.SQLItem](sql.db).Where("id = ?", cid).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return course.ErrCourseNotFound
		}

		return err
	}

	s, err := sql.conn.Where("id = ?", sid).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStudentNotFound
		}

		return err
	}

	s.Courses = append(s.Courses, c)

	_, err = sql.conn.Updates(ctx, s)
	if err != nil {
		return err
	}

	return nil
}

func (sql SQL) Get(ctx context.Context, id string) (model.Student, error) {
	// st contains single students repeated multiple times
	// to contains the course information using join.
	//
	// Here joining will remove the n+1 issue which happens
	// with Preload().
	var st []struct {
		ID          string
		Name        string
		CoursesID   *string
		CoursesName *string
	}

	err := sql.db.Table("students").
		Joins("LEFT JOIN `students_courses` ON `students`.`id` = `students_courses`.`sql_item_id`").
		Joins("LEFT JOIN (select id courses_id, name courses_name from `courses`) ON "+
			"`courses_id` = `students_courses`.`course_id`").
		Where("students.id = ?", id).Scan(&st).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Student{}, ErrStudentNotFound
		}

		return model.Student{}, err
	}

	courses := make([]model.Course, 0, len(st))

	for _, course := range st {
		if course.CoursesID != nil {
			courses = append(courses, model.Course{
				Name: *course.CoursesName,
				ID:   *course.CoursesID,
			})
		}
	}

	return model.Student{
		Name:    st[0].Name,
		ID:      st[0].ID,
		Courses: courses,
	}, nil
}
