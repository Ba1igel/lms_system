package repository

import (
	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/pkg/pagination"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *model.Course) error
	GetByID(id uint) (*model.Course, error)
	GetAll(p pagination.Params) ([]model.Course, int64, error)
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
	if err := r.db.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) GetAll(p pagination.Params) ([]model.Course, int64, error) {
	var courses []model.Course
	var total int64

	if err := r.db.Model(&model.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(p.Limit).Offset(p.Offset()).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *courseRepository) Update(course *model.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uint) error {
	return r.db.Delete(&model.Course{}, id).Error
}
