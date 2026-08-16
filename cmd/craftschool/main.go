package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"crafttraining/internal/training"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	admin := training.Actor{ID: "admin-demo", Role: training.RoleAdmin}
	teacher := training.Actor{ID: "teacher-demo", Role: training.RoleTeacher}
	fixedTime := time.Date(2026, time.August, 16, 9, 30, 0, 0, time.UTC)
	service := training.NewService(training.NewMemoryRepository(), func() time.Time { return fixedTime })

	for _, category := range []training.Category{
		{ID: "woodwork", Name: "木工"},
		{ID: "pottery", Name: "陶艺"},
		{ID: "leatherwork", Name: "皮具"},
		{ID: "metalwork", Name: "金工"},
	} {
		if err := service.CreateCategory(admin, category); err != nil {
			return fmt.Errorf("创建分类: %w", err)
		}
	}

	course := training.Course{ID: "leather-card-holder", Title: "手缝植鞣革卡包", CategoryID: "leatherwork"}
	if err := service.CreateCourse(teacher, course); err != nil {
		return fmt.Errorf("创建课程: %w", err)
	}
	if err := service.EditCourse(teacher, course.ID, []string{"植鞣革", "蜡线"}, []string{"菱斩", "木槌", "针"}, []string{"佩戴护目镜", "刀具及时入套"}); err != nil {
		return fmt.Errorf("编辑课程: %w", err)
	}
	if err := service.SubmitForReview(teacher, course.ID); err != nil {
		return fmt.Errorf("提交审核: %w", err)
	}
	if err := service.ApproveCourse(admin, course.ID); err != nil {
		return fmt.Errorf("审核课程: %w", err)
	}
	if err := service.PublishCourse(teacher, course.ID); err != nil {
		return fmt.Errorf("公开课程: %w", err)
	}
	if err := service.ArchiveCourse(admin, course.ID); err != nil {
		return fmt.Errorf("归档课程: %w", err)
	}

	stored, err := service.Course(course.ID)
	if err != nil {
		return fmt.Errorf("读取课程: %w", err)
	}
	categories := service.Categories()
	names := make([]string, 0, len(categories))
	for _, category := range categories {
		names = append(names, category.Name)
	}
	slices.Sort(names)
	fmt.Printf("分类: %v\n", names)
	fmt.Printf("课程: %s\n状态: %s\n", stored.Title, stored.Status)
	return nil
}
