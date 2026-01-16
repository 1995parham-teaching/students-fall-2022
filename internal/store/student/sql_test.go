package student_test

import (
	"context"
	"errors"
	"testing"

	"github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/store/course"
	"github.com/1995parham-teaching/students/internal/store/student"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{ //nolint:exhaustruct
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	return db
}

func TestSQL_Get_NotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	_, err := store.Get(ctx, "99999999")
	if err == nil {
		t.Error("expected error when getting non-existing student, got nil")
	}

	if !errors.Is(err, student.ErrStudentNotFound) {
		t.Errorf("expected ErrStudentNotFound, got %v", err)
	}
}

func TestSQL_Get_ExistingStudent(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	expected := model.Student{
		ID:      "12345678",
		Name:    "Parham Alvani",
		Courses: nil,
	}

	err := store.Create(ctx, expected)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	got, err := store.Get(ctx, expected.ID)
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}

	if got.ID != expected.ID {
		t.Errorf("expected ID %s, got %s", expected.ID, got.ID)
	}

	if got.Name != expected.Name {
		t.Errorf("expected Name %s, got %s", expected.Name, got.Name)
	}

	if len(got.Courses) != 0 {
		t.Errorf("expected 0 courses, got %d", len(got.Courses))
	}
}

func TestSQL_Get_StudentWithCourses(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	studentStore := student.NewSQL(db)
	courseStore := course.NewSQL(db)
	ctx := context.Background()

	// Create courses first
	courses := []model.Course{
		{ID: "10101010", Name: "Internet Engineering"},
		{ID: "20202020", Name: "Database Design"},
		{ID: "30303030", Name: "Operating Systems"},
	}

	for _, c := range courses {
		err := courseStore.Create(ctx, c)
		if err != nil {
			t.Fatalf("failed to create course %s: %v", c.Name, err)
		}
	}

	// Create student
	st := model.Student{
		ID:      "12345678",
		Name:    "Parham Alvani",
		Courses: nil,
	}

	err := studentStore.Create(ctx, st)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	// Register student for courses
	for _, c := range courses {
		if err := studentStore.Register(ctx, st.ID, c.ID); err != nil {
			t.Fatalf("failed to register student for course %s: %v", c.Name, err)
		}
	}

	// Get student and verify courses
	got, err := studentStore.Get(ctx, st.ID)
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}

	if len(got.Courses) != len(courses) {
		t.Errorf("expected %d courses, got %d", len(courses), len(got.Courses))
	}

	verifyCourses(t, got.Courses, courses)
}

func verifyCourses(t *testing.T, got []model.Course, expected []model.Course) {
	t.Helper()

	courseMap := make(map[string]model.Course)
	for _, c := range got {
		courseMap[c.ID] = c
	}

	for _, exp := range expected {
		c, ok := courseMap[exp.ID]
		if !ok {
			t.Errorf("course %s not found in student courses", exp.ID)

			continue
		}

		if c.Name != exp.Name {
			t.Errorf("course name mismatch: expected %s, got %s", exp.Name, c.Name)
		}
	}
}

func TestSQL_Create_Success(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	st := model.Student{
		ID:      "12345678",
		Name:    "Parham Alvani",
		Courses: nil,
	}

	err := store.Create(ctx, st)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	// Verify student was created
	got, err := store.Get(ctx, st.ID)
	if err != nil {
		t.Fatalf("failed to get created student: %v", err)
	}

	if got.ID != st.ID || got.Name != st.Name {
		t.Errorf("created student mismatch: expected %+v, got %+v", st, got)
	}
}

func TestSQL_Create_DuplicateID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	st := model.Student{
		ID:      "12345678",
		Name:    "Parham Alvani",
		Courses: nil,
	}

	err := store.Create(ctx, st)
	if err != nil {
		t.Fatalf("failed to create first student: %v", err)
	}

	// Try to create another student with the same ID
	duplicate := model.Student{
		ID:      "12345678",
		Name:    "Another Person",
		Courses: nil,
	}

	err = store.Create(ctx, duplicate)
	if err == nil {
		t.Error("expected error when creating duplicate student, got nil")
	}
}

func TestSQL_GetAll_Empty(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	students, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all students: %v", err)
	}

	if len(students) != 0 {
		t.Errorf("expected 0 students, got %d", len(students))
	}
}

func TestSQL_GetAll_MultipleStudents(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	expected := []model.Student{
		{ID: "11111111", Name: "Student One", Courses: nil},
		{ID: "22222222", Name: "Student Two", Courses: nil},
		{ID: "33333333", Name: "Student Three", Courses: nil},
	}

	for _, st := range expected {
		if err := store.Create(ctx, st); err != nil {
			t.Fatalf("failed to create student %s: %v", st.Name, err)
		}
	}

	got, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all students: %v", err)
	}

	if len(got) != len(expected) {
		t.Errorf("expected %d students, got %d", len(expected), len(got))
	}

	verifyStudents(t, got, expected)
}

func verifyStudents(t *testing.T, got []model.Student, expected []model.Student) {
	t.Helper()

	studentMap := make(map[string]model.Student)
	for _, st := range got {
		studentMap[st.ID] = st
	}

	for _, exp := range expected {
		st, ok := studentMap[exp.ID]
		if !ok {
			t.Errorf("student %s not found", exp.ID)

			continue
		}

		if st.Name != exp.Name {
			t.Errorf("student name mismatch: expected %s, got %s", exp.Name, st.Name)
		}
	}
}

func TestSQL_GetAll_WithCourses(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	studentStore := student.NewSQL(db)
	courseStore := course.NewSQL(db)
	ctx := context.Background()

	// Create a course
	c := model.Course{ID: "10101010", Name: "Internet Engineering"}

	err := courseStore.Create(ctx, c)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	// Create students
	students := []model.Student{
		{ID: "11111111", Name: "Student One", Courses: nil},
		{ID: "22222222", Name: "Student Two", Courses: nil},
	}

	for _, st := range students {
		if err := studentStore.Create(ctx, st); err != nil {
			t.Fatalf("failed to create student %s: %v", st.Name, err)
		}
	}

	// Register first student for the course
	err = studentStore.Register(ctx, students[0].ID, c.ID)
	if err != nil {
		t.Fatalf("failed to register student: %v", err)
	}

	// Get all and verify
	got, err := studentStore.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all students: %v", err)
	}

	verifyStudentCourses(t, got, students)
}

func verifyStudentCourses(t *testing.T, got []model.Student, students []model.Student) {
	t.Helper()

	for _, st := range got {
		switch st.ID {
		case students[0].ID:
			if len(st.Courses) != 1 {
				t.Errorf("expected 1 course for student %s, got %d", st.ID, len(st.Courses))
			}
		case students[1].ID:
			if len(st.Courses) != 0 {
				t.Errorf("expected 0 courses for student %s, got %d", st.ID, len(st.Courses))
			}
		}
	}
}

func TestSQL_Register_Success(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	studentStore := student.NewSQL(db)
	courseStore := course.NewSQL(db)
	ctx := context.Background()

	// Create course and student
	c := model.Course{ID: "10101010", Name: "Internet Engineering"}

	err := courseStore.Create(ctx, c)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	st := model.Student{ID: "12345678", Name: "Parham Alvani", Courses: nil}

	err = studentStore.Create(ctx, st)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	// Register
	err = studentStore.Register(ctx, st.ID, c.ID)
	if err != nil {
		t.Fatalf("failed to register student for course: %v", err)
	}

	// Verify registration
	got, err := studentStore.Get(ctx, st.ID)
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}

	if len(got.Courses) != 1 {
		t.Fatalf("expected 1 course, got %d", len(got.Courses))
	}

	if got.Courses[0].ID != c.ID {
		t.Errorf("expected course ID %s, got %s", c.ID, got.Courses[0].ID)
	}
}

func TestSQL_Register_StudentNotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	studentStore := student.NewSQL(db)
	courseStore := course.NewSQL(db)
	ctx := context.Background()

	// Create course only
	c := model.Course{ID: "10101010", Name: "Internet Engineering"}

	err := courseStore.Create(ctx, c)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	// Try to register non-existing student
	err = studentStore.Register(ctx, "99999999", c.ID)
	if err == nil {
		t.Error("expected error when registering non-existing student, got nil")
	}

	if !errors.Is(err, student.ErrStudentNotFound) {
		t.Errorf("expected ErrStudentNotFound, got %v", err)
	}
}

func TestSQL_Register_CourseNotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	studentStore := student.NewSQL(db)
	ctx := context.Background()

	// Create student only
	st := model.Student{ID: "12345678", Name: "Parham Alvani", Courses: nil}

	err := studentStore.Create(ctx, st)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	// Try to register for non-existing course
	err = studentStore.Register(ctx, st.ID, "99999999")
	if err == nil {
		t.Error("expected error when registering for non-existing course, got nil")
	}

	if !errors.Is(err, course.ErrCourseNotFound) {
		t.Errorf("expected ErrCourseNotFound, got %v", err)
	}
}

func TestSQL_Register_MultipleCourses(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	studentStore := student.NewSQL(db)
	courseStore := course.NewSQL(db)
	ctx := context.Background()

	// Create multiple courses
	courses := []model.Course{
		{ID: "10101010", Name: "Internet Engineering"},
		{ID: "20202020", Name: "Database Design"},
	}

	for _, c := range courses {
		if err := courseStore.Create(ctx, c); err != nil {
			t.Fatalf("failed to create course %s: %v", c.Name, err)
		}
	}

	// Create student
	st := model.Student{ID: "12345678", Name: "Parham Alvani", Courses: nil}

	if err := studentStore.Create(ctx, st); err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	// Register for all courses
	for _, c := range courses {
		if err := studentStore.Register(ctx, st.ID, c.ID); err != nil {
			t.Fatalf("failed to register for course %s: %v", c.Name, err)
		}
	}

	// Verify all registrations
	got, err := studentStore.Get(ctx, st.ID)
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}

	if len(got.Courses) != len(courses) {
		t.Errorf("expected %d courses, got %d", len(courses), len(got.Courses))
	}
}

func TestSQL_Get_EmptyDatabase(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	// This specifically tests the panic fix - getting from empty database
	_, err := store.Get(ctx, "12345678")
	if !errors.Is(err, student.ErrStudentNotFound) {
		t.Errorf("expected ErrStudentNotFound, got %v", err)
	}
}

func TestSQL_Get_WrongID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := student.NewSQL(db)
	ctx := context.Background()

	// Create a student
	st := model.Student{ID: "12345678", Name: "Parham Alvani", Courses: nil}

	err := store.Create(ctx, st)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	// Try to get with different ID - this tests the fix
	_, err = store.Get(ctx, "99999999")
	if !errors.Is(err, student.ErrStudentNotFound) {
		t.Errorf("expected ErrStudentNotFound, got %v", err)
	}
}
