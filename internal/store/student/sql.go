package student

import (
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

	if err := sql.DB.Model(new(SQLItem)).
		Select("Name", "ID", "Courses.ID", "Courses.Name").
		Joins("LEFT JOIN `students_courses` ON `students`.`id` = `students_courses`.`sql_item_id`").
		Joins("LEFT JOIN `courses` ON `courses`.`id` = `students_courses`.`course_id`").
		Find(&items).Error; err != nil {
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

func (sql SQL) Create(s model.Student) error {
	return sql.DB.Create(&SQLItem{
		ID:      s.ID,
		Name:    s.Name,
		Courses: nil,
	}).Error
}

func (sql SQL) Register(sid string, cid string) error {
	var c course.SQLItem

	if err := sql.DB.First(&c, cid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return course.ErrCourseNotFound
		}

		return err
	}

	var s SQLItem

	if err := sql.DB.Model(new(SQLItem)).First(&s, sid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStudentNotFound
		}

		return err
	}

	s.Courses = append(s.Courses, c)

	return sql.DB.Save(&s).Error
}

func (sql SQL) Get(id string) (model.Student, error) {
	var st []struct {
		ID          string
		Name        string
		CoursesID   string
		CoursesName string
	}

	if err := sql.DB.Table("students").
		Joins("LEFT JOIN `students_courses` ON `students`.`id` = `students_courses`.`sql_item_id`").
		Joins("LEFT JOIN (select id courses_id, name courses_name from `courses`) ON `courses_id` = `students_courses`.`course_id`").
		Where("students.id = ?", id).Scan(&st).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Student{}, ErrStudentNotFound
		}

		return model.Student{}, err
	}

	log.Println(st)

	courses := make([]model.Course, 0, len(st))
	for _, course := range st {
		courses = append(courses, model.Course{
			Name: course.CoursesName,
			ID:   course.CoursesID,
		})
	}

	return model.Student{
		Name:    st[0].Name,
		ID:      st[0].ID,
		Courses: courses,
	}, nil
}
