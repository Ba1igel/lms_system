package repository

import (
	"github.com/baigel/lms/main-service/internal/model"
	"gorm.io/gorm"
)

type ChapterRepository interface {
	Create(chapter *model.Chapter) error
	GetByID(id uint) (*model.Chapter, error)
	GetAllByCourseID(courseID uint) ([]model.Chapter, error)
	Update(chapter *model.Chapter) error
	Delete(id uint) error
}

type chapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) ChapterRepository {
	return &chapterRepository{db: db}
}

func (r *chapterRepository) Create(chapter *model.Chapter) error {
	return r.db.Create(chapter).Error
}

func (r *chapterRepository) GetByID(id uint) (*model.Chapter, error) {
	var chapter model.Chapter
	err := r.db.First(&chapter, id).Error
	return &chapter, err
}

func (r *chapterRepository) GetAllByCourseID(courseID uint) ([]model.Chapter, error) {
	var chapters []model.Chapter
	err := r.db.Where("course_id = ?", courseID).Order("\"order\" asc").Find(&chapters).Error
	return chapters, err
}

func (r *chapterRepository) Update(chapter *model.Chapter) error {
	return r.db.Save(chapter).Error
}

func (r *chapterRepository) Delete(id uint) error {
	return r.db.Delete(&model.Chapter{}, id).Error
}
