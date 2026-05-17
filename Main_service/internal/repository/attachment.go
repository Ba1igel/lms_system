package repository

import (
	"github.com/baigel/lms/main-service/internal/model"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(attachment *model.Attachment) error
	GetByID(id uint) (*model.Attachment, error)
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(attachment *model.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) GetByID(id uint) (*model.Attachment, error) {
	var attachment model.Attachment
	if err := r.db.First(&attachment, id).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}
