package training_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"crafttraining/internal/training"
)

var (
	admin   = training.Actor{ID: "admin-1", Role: training.RoleAdmin}
	teacher = training.Actor{ID: "teacher-1", Role: training.RoleTeacher}
)

func TestCourseBusinessLifecycle(t *testing.T) {
	service := newFixture(t)

	err := service.CreateCourse(teacher, training.Course{
		ID: "course-1", Title: "手缝植鞣革卡包", CategoryID: "leatherwork",
	})
	if err != nil {
		t.Fatalf("create course: %v", err)
	}

	err = service.EditCourse(teacher, "course-1",
		[]string{"植鞣革", "蜡线"},
		[]string{"菱斩", "木槌", "针"},
		[]string{"敲击时佩戴护目镜", "刀具离手后收入刀套"},
	)
	if err != nil {
		t.Fatalf("edit course: %v", err)
	}
	if err = service.SubmitForReview(teacher, "course-1"); err != nil {
		t.Fatalf("submit course: %v", err)
	}
	if err = service.ApproveCourse(admin, "course-1"); err != nil {
		t.Fatalf("approve course: %v", err)
	}
	if err = service.PublishCourse(teacher, "course-1"); err != nil {
		t.Fatalf("publish course: %v", err)
	}
	if err = service.ArchiveCourse(admin, "course-1"); err != nil {
		t.Fatalf("archive course: %v", err)
	}

	course, err := service.Course("course-1")
	if err != nil {
		t.Fatalf("load course: %v", err)
	}
	if course.Status != training.StatusArchived {
		t.Fatalf("status = %q, want %q", course.Status, training.StatusArchived)
	}
	if course.ReviewedBy != admin.ID || course.ReviewedAt.IsZero() || course.PublishedAt.IsZero() || course.ArchivedAt.IsZero() {
		t.Fatalf("lifecycle timestamps or reviewer missing: %+v", course)
	}
}

func TestCategoryAdministrationAndCourseEditing(t *testing.T) {
	service := newFixture(t)

	categories := service.Categories()
	names := make([]string, 0, len(categories))
	for _, category := range categories {
		names = append(names, category.Name)
	}
	slices.Sort(names)
	if !slices.Equal(names, []string{"木工", "皮具", "金工", "陶艺"}) {
		t.Fatalf("category names = %v", names)
	}

	if err := service.RenameCategory(teacher, "metalwork", "锻造"); !errors.Is(err, training.ErrForbidden) {
		t.Fatalf("teacher rename error = %v", err)
	}
	if err := service.RenameCategory(admin, "metalwork", "金属工艺"); err != nil {
		t.Fatalf("admin rename: %v", err)
	}
	if err := service.SetCategoryActive(admin, "pottery", false); err != nil {
		t.Fatalf("deactivate category: %v", err)
	}
	if err := service.CreateCourse(teacher, training.Course{ID: "course-2", Title: "拉坯入门", CategoryID: "pottery"}); !errors.Is(err, training.ErrInvalidInput) {
		t.Fatalf("inactive category create error = %v", err)
	}

	if err := service.CreateCourse(teacher, training.Course{ID: "course-3", Title: "榫卯基础", CategoryID: "woodwork"}); err != nil {
		t.Fatalf("create editable course: %v", err)
	}
	if err := service.SubmitForReview(teacher, "course-3"); !errors.Is(err, training.ErrInvalidInput) {
		t.Fatalf("incomplete submit error = %v", err)
	}
	if err := service.EditCourse(training.Actor{ID: "teacher-2", Role: training.RoleTeacher}, "course-3", []string{"木料"}, []string{"凿子"}, []string{"夹紧工件"}); !errors.Is(err, training.ErrForbidden) {
		t.Fatalf("other teacher edit error = %v", err)
	}
}

func TestArchivedCourseCannotBePublished(t *testing.T) {
	service := newFixture(t)
	createArchivedCourse(t, service)

	err := service.PublishCourse(teacher, "archived-course")
	if !errors.Is(err, training.ErrInvalidTransition) {
		t.Errorf("publish archived course error = %v, want %v", err, training.ErrInvalidTransition)
	}
	course, loadErr := service.Course("archived-course")
	if loadErr != nil {
		t.Fatalf("load archived course: %v", loadErr)
	}
	if course.Status != training.StatusArchived {
		t.Errorf("status = %q, want %q", course.Status, training.StatusArchived)
	}
}

func newFixture(t *testing.T) *training.Service {
	t.Helper()
	fixedTime := time.Date(2026, time.August, 16, 9, 30, 0, 0, time.UTC)
	service := training.NewService(training.NewMemoryRepository(), func() time.Time { return fixedTime })
	categories := []training.Category{
		{ID: "woodwork", Name: "木工"},
		{ID: "pottery", Name: "陶艺"},
		{ID: "leatherwork", Name: "皮具"},
		{ID: "metalwork", Name: "金工"},
	}
	for _, category := range categories {
		if err := service.CreateCategory(admin, category); err != nil {
			t.Fatalf("create category %q: %v", category.ID, err)
		}
	}
	return service
}

func createArchivedCourse(t *testing.T, service *training.Service) {
	t.Helper()
	err := service.CreateCourse(teacher, training.Course{ID: "archived-course", Title: "停课课程", CategoryID: "woodwork"})
	if err != nil {
		t.Fatalf("create course: %v", err)
	}
	if err = service.EditCourse(teacher, "archived-course", []string{"榉木"}, []string{"手锯"}, []string{"固定材料"}); err != nil {
		t.Fatalf("edit course: %v", err)
	}
	if err = service.SubmitForReview(teacher, "archived-course"); err != nil {
		t.Fatalf("submit course: %v", err)
	}
	if err = service.ApproveCourse(admin, "archived-course"); err != nil {
		t.Fatalf("approve course: %v", err)
	}
	if err = service.PublishCourse(teacher, "archived-course"); err != nil {
		t.Fatalf("publish course: %v", err)
	}
	if err = service.ArchiveCourse(admin, "archived-course"); err != nil {
		t.Fatalf("archive course: %v", err)
	}
}
