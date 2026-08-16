package training

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleTeacher Role = "teacher"
)

type Actor struct {
	ID   string
	Role Role
}

type Category struct {
	ID     string
	Name   string
	Active bool
}

type CourseStatus string

const (
	StatusDraft         CourseStatus = "draft"
	StatusPendingReview CourseStatus = "pending_review"
	StatusApproved      CourseStatus = "approved"
	StatusPublished     CourseStatus = "published"
	StatusArchived      CourseStatus = "archived"
)

type Course struct {
	ID          string
	Title       string
	CategoryID  string
	TeacherID   string
	Materials   []string
	Tools       []string
	SafetyTips  []string
	Status      CourseStatus
	ReviewedBy  string
	ReviewedAt  time.Time
	PublishedAt time.Time
	ArchivedAt  time.Time
}

func cloneCourse(course Course) Course {
	course.Materials = append([]string(nil), course.Materials...)
	course.Tools = append([]string(nil), course.Tools...)
	course.SafetyTips = append([]string(nil), course.SafetyTips...)
	return course
}
