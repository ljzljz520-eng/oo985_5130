package training

import (
	"errors"
	"sort"
	"sync"
)

var ErrNotFound = errors.New("not found")

type MemoryRepository struct {
	mu         sync.RWMutex
	categories map[string]Category
	courses    map[string]Course
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		categories: make(map[string]Category),
		courses:    make(map[string]Course),
	}
}

func (repository *MemoryRepository) CreateCategory(category Category) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if _, exists := repository.categories[category.ID]; exists {
		return ErrAlreadyExists
	}
	repository.categories[category.ID] = category
	return nil
}

func (repository *MemoryRepository) UpdateCategory(id string, update func(*Category) error) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	category, exists := repository.categories[id]
	if !exists {
		return ErrNotFound
	}
	if err := update(&category); err != nil {
		return err
	}
	repository.categories[id] = category
	return nil
}

func (repository *MemoryRepository) Category(id string) (Category, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	category, exists := repository.categories[id]
	if !exists {
		return Category{}, ErrNotFound
	}
	return category, nil
}

func (repository *MemoryRepository) Categories() []Category {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	categories := make([]Category, 0, len(repository.categories))
	for _, category := range repository.categories {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})
	return categories
}

func (repository *MemoryRepository) CreateCourse(course Course) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if _, exists := repository.courses[course.ID]; exists {
		return ErrAlreadyExists
	}
	repository.courses[course.ID] = cloneCourse(course)
	return nil
}

func (repository *MemoryRepository) UpdateCourse(id string, update func(*Course) error) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	course, exists := repository.courses[id]
	if !exists {
		return ErrNotFound
	}
	course = cloneCourse(course)
	if err := update(&course); err != nil {
		return err
	}
	repository.courses[id] = cloneCourse(course)
	return nil
}

func (repository *MemoryRepository) Course(id string) (Course, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	course, exists := repository.courses[id]
	if !exists {
		return Course{}, ErrNotFound
	}
	return cloneCourse(course), nil
}
