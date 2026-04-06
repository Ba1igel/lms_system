package repository

import (
	"github.com/baigel/lms/main-service/internal/model"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *model.Course) error
	GetByID(id uint) (*model.Course, error)
	GetAll() ([]model.Course, error)
	Update(course *model.Course) error
	Delete(id uint) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) GetByID(id uint) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, id).Error
	return &course, err
}

func (r *courseRepository) GetAll() ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Find(&courses).Error
	return courses, err
}

func (r *courseRepository) Update(course *model.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uint) error {
	return r.db.Delete(&model.Course{}, id).Error
}
