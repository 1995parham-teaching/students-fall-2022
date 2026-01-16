package course_test

import (
	"context"
	"testing"

	"github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/store/course"
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

func TestSQL_Create_Success(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	c := model.Course{
		ID:   "10101010",
		Name: "Internet Engineering",
	}

	err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	// Verify course was created
	got, err := store.Get(ctx, c.ID)
	if err != nil {
		t.Fatalf("failed to get created course: %v", err)
	}

	if got.ID != c.ID || got.Name != c.Name {
		t.Errorf("created course mismatch: expected %+v, got %+v", c, got)
	}
}

func TestSQL_Create_DuplicateID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	c := model.Course{
		ID:   "10101010",
		Name: "Internet Engineering",
	}

	err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("failed to create first course: %v", err)
	}

	// Try to create another course with the same ID
	duplicate := model.Course{
		ID:   "10101010",
		Name: "Another Course",
	}

	err = store.Create(ctx, duplicate)
	if err == nil {
		t.Error("expected error when creating duplicate course, got nil")
	}
}

func TestSQL_Get_NotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	_, err := store.Get(ctx, "99999999")
	if err == nil {
		t.Error("expected error when getting non-existing course, got nil")
	}
}

func TestSQL_Get_ExistingCourse(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	expected := model.Course{
		ID:   "10101010",
		Name: "Internet Engineering",
	}

	err := store.Create(ctx, expected)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	got, err := store.Get(ctx, expected.ID)
	if err != nil {
		t.Fatalf("failed to get course: %v", err)
	}

	if got.ID != expected.ID {
		t.Errorf("expected ID %s, got %s", expected.ID, got.ID)
	}

	if got.Name != expected.Name {
		t.Errorf("expected Name %s, got %s", expected.Name, got.Name)
	}
}

func TestSQL_GetAll_Empty(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	courses, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all courses: %v", err)
	}

	if len(courses) != 0 {
		t.Errorf("expected 0 courses, got %d", len(courses))
	}
}

func TestSQL_GetAll_MultipleCourses(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	expected := []model.Course{
		{ID: "10101010", Name: "Internet Engineering"},
		{ID: "20202020", Name: "Database Design"},
		{ID: "30303030", Name: "Operating Systems"},
	}

	for _, c := range expected {
		err := store.Create(ctx, c)
		if err != nil {
			t.Fatalf("failed to create course %s: %v", c.Name, err)
		}
	}

	got, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all courses: %v", err)
	}

	if len(got) != len(expected) {
		t.Errorf("expected %d courses, got %d", len(expected), len(got))
	}

	verifyCourses(t, got, expected)
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
			t.Errorf("course %s not found", exp.ID)

			continue
		}

		if c.Name != exp.Name {
			t.Errorf("course name mismatch: expected %s, got %s", exp.Name, c.Name)
		}
	}
}

func TestSQL_Get_EmptyDatabase(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	_, err := store.Get(ctx, "10101010")
	if err == nil {
		t.Error("expected error when getting from empty database, got nil")
	}
}

func TestSQL_Get_WrongID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := course.NewSQL(db)
	ctx := context.Background()

	// Create a course
	c := model.Course{ID: "10101010", Name: "Internet Engineering"}

	err := store.Create(ctx, c)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	// Try to get with different ID
	_, err = store.Get(ctx, "99999999")
	if err == nil {
		t.Error("expected error when getting with wrong ID, got nil")
	}
}
