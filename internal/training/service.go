package training

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrAlreadyExists     = errors.New("already exists")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidTransition = errors.New("invalid course transition")
)

type Clock func() time.Time

type Service struct {
	repository *MemoryRepository
	now        Clock
}

func NewService(repository *MemoryRepository, now Clock) *Service {
	if repository == nil {
		repository = NewMemoryRepository()
	}
	if now == nil {
		now = time.Now
	}
	return &Service{repository: repository, now: now}
}

func (service *Service) CreateCategory(actor Actor, category Category) error {
	if err := requireRole(actor, RoleAdmin); err != nil {
		return err
	}
	category.ID = strings.TrimSpace(category.ID)
	category.Name = strings.TrimSpace(category.Name)
	if category.ID == "" || category.Name == "" {
		return ErrInvalidInput
	}
	category.Active = true
	return service.repository.CreateCategory(category)
}

func (service *Service) RenameCategory(actor Actor, categoryID, name string) error {
	if err := requireRole(actor, RoleAdmin); err != nil {
		return err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidInput
	}
	return service.repository.UpdateCategory(categoryID, func(category *Category) error {
		category.Name = name
		return nil
	})
}

func (service *Service) SetCategoryActive(actor Actor, categoryID string, active bool) error {
	if err := requireRole(actor, RoleAdmin); err != nil {
		return err
	}
	return service.repository.UpdateCategory(categoryID, func(category *Category) error {
		category.Active = active
		return nil
	})
}

func (service *Service) Categories() []Category {
	return service.repository.Categories()
}

func (service *Service) CreateCourse(actor Actor, course Course) error {
	if err := requireRole(actor, RoleTeacher); err != nil {
		return err
	}
	course.ID = strings.TrimSpace(course.ID)
	course.Title = strings.TrimSpace(course.Title)
	course.CategoryID = strings.TrimSpace(course.CategoryID)
	if course.ID == "" || course.Title == "" || course.CategoryID == "" || actor.ID == "" {
		return ErrInvalidInput
	}
	category, err := service.repository.Category(course.CategoryID)
	if err != nil {
		return fmt.Errorf("category: %w", err)
	}
	if !category.Active {
		return ErrInvalidInput
	}
	course.TeacherID = actor.ID
	course.Status = StatusDraft
	course.Materials = nil
	course.Tools = nil
	course.SafetyTips = nil
	course.ReviewedBy = ""
	course.ReviewedAt = time.Time{}
	course.PublishedAt = time.Time{}
	course.ArchivedAt = time.Time{}
	return service.repository.CreateCourse(course)
}

func (service *Service) EditCourse(actor Actor, courseID string, materials, tools, safetyTips []string) error {
	if err := requireRole(actor, RoleTeacher); err != nil {
		return err
	}
	cleanMaterials := cleanItems(materials)
	cleanTools := cleanItems(tools)
	cleanSafetyTips := cleanItems(safetyTips)
	return service.repository.UpdateCourse(courseID, func(course *Course) error {
		if course.TeacherID != actor.ID {
			return ErrForbidden
		}
		if course.Status != StatusDraft {
			return ErrInvalidTransition
		}
		course.Materials = cleanMaterials
		course.Tools = cleanTools
		course.SafetyTips = cleanSafetyTips
		return nil
	})
}

func (service *Service) SubmitForReview(actor Actor, courseID string) error {
	if err := requireRole(actor, RoleTeacher); err != nil {
		return err
	}
	stored, err := service.repository.Course(courseID)
	if err != nil {
		return err
	}
	category, err := service.repository.Category(stored.CategoryID)
	if err != nil || !category.Active {
		return ErrInvalidInput
	}
	return service.repository.UpdateCourse(courseID, func(course *Course) error {
		if course.TeacherID != actor.ID {
			return ErrForbidden
		}
		if course.Status != StatusDraft {
			return ErrInvalidTransition
		}
		if len(course.Materials) == 0 || len(course.Tools) == 0 || len(course.SafetyTips) == 0 {
			return ErrInvalidInput
		}
		course.Status = StatusPendingReview
		return nil
	})
}

func (service *Service) ApproveCourse(actor Actor, courseID string) error {
	if err := requireRole(actor, RoleAdmin); err != nil {
		return err
	}
	return service.repository.UpdateCourse(courseID, func(course *Course) error {
		if course.Status != StatusPendingReview {
			return ErrInvalidTransition
		}
		course.Status = StatusApproved
		course.ReviewedBy = actor.ID
		course.ReviewedAt = service.now()
		return nil
	})
}

func (service *Service) PublishCourse(actor Actor, courseID string) error {
	if err := requireRole(actor, RoleTeacher); err != nil {
		return err
	}
	return service.repository.UpdateCourse(courseID, func(course *Course) error {
		if course.TeacherID != actor.ID {
			return ErrForbidden
		}
		if course.ReviewedAt.IsZero() {
			return ErrInvalidTransition
		}
		course.Status = StatusPublished
		course.PublishedAt = service.now()
		return nil
	})
}

func (service *Service) ArchiveCourse(actor Actor, courseID string) error {
	if err := requireRole(actor, RoleAdmin); err != nil {
		return err
	}
	return service.repository.UpdateCourse(courseID, func(course *Course) error {
		if course.Status != StatusPublished {
			return ErrInvalidTransition
		}
		course.Status = StatusArchived
		course.ArchivedAt = service.now()
		return nil
	})
}

func (service *Service) Course(courseID string) (Course, error) {
	return service.repository.Course(courseID)
}

func requireRole(actor Actor, role Role) error {
	if actor.ID == "" || actor.Role != role {
		return ErrForbidden
	}
	return nil
}

func cleanItems(items []string) []string {
	cleaned := make([]string, 0, len(items))
	for _, item := range items {
		if value := strings.TrimSpace(item); value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return cleaned
}
