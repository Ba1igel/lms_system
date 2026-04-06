package repository

import (
	"github.com/baigel/lms/main-service/internal/model"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(lesson *model.Lesson) error
	GetByID(id uint) (*model.Lesson, error)
	GetAllByChapterID(chapterID uint) ([]model.Lesson, error)
	Update(lesson *model.Lesson) error
	Delete(id uint) error
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) Create(lesson *model.Lesson) error {
	return r.db.Create(lesson).Error
}

func (r *lessonRepository) GetByID(id uint) (*model.Lesson, error) {
	var lesson model.Lesson
	err := r.db.First(&lesson, id).Error
	return &lesson, err
}

func (r *lessonRepository) GetAllByChapterID(chapterID uint) ([]model.Lesson, error) {
	var lessons []model.Lesson
	err := r.db.Where("chapter_id = ?", chapterID).Order("\"order\" asc").Find(&lessons).Error
	return lessons, err
}

func (r *lessonRepository) Update(lesson *model.Lesson) error {
	return r.db.Save(lesson).Error
}

func (r *lessonRepository) Delete(id uint) error {
	return r.db.Delete(&model.Lesson{}, id).Error
}
