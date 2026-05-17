package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/baigel/lms/main-service/internal/model"
	"github.com/baigel/lms/main-service/internal/repository"
	"github.com/baigel/lms/main-service/pkg/apperror"
	"github.com/baigel/lms/main-service/pkg/storage"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

const maxAttachmentNameLength = 255

type AttachmentUpload struct {
	LessonID    uint
	Name        string
	Reader      io.Reader
	Size        int64
	ContentType string
}

type AttachmentDownload struct {
	Attachment  *model.Attachment
	File        *storage.DownloadResult
	ContentName string
}

type AttachmentService interface {
	Upload(ctx context.Context, input AttachmentUpload) (*model.Attachment, error)
	Download(ctx context.Context, id uint, roles []string) (*AttachmentDownload, error)
}

type attachmentService struct {
	attachments repository.AttachmentRepository
	lessons     repository.LessonRepository
	storage     storage.FileStorage
}

func NewAttachmentService(
	attachments repository.AttachmentRepository,
	lessons repository.LessonRepository,
	storage storage.FileStorage,
) AttachmentService {
	return &attachmentService{
		attachments: attachments,
		lessons:     lessons,
		storage:     storage,
	}
}

func (s *attachmentService) Upload(ctx context.Context, input AttachmentUpload) (*model.Attachment, error) {
	if input.LessonID == 0 {
		return nil, apperror.BadRequest("lesson_id is required", nil)
	}
	if input.Reader == nil || input.Size <= 0 {
		return nil, apperror.BadRequest("file is required", nil)
	}

	fileName := sanitizeFileName(input.Name)
	if fileName == "" {
		return nil, apperror.BadRequest("file name is required", nil)
	}

	if _, err := s.lessons.GetByID(input.LessonID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("lesson not found")
		}
		return nil, apperror.Internal("failed to get lesson", err)
	}

	objectName := fmt.Sprintf("lessons/%d/%s-%s", input.LessonID, ksuid.New().String(), fileName)
	fileURL, err := s.storage.Upload(ctx, storage.UploadInput{
		ObjectName:  objectName,
		Reader:      input.Reader,
		Size:        input.Size,
		ContentType: input.ContentType,
	})
	if err != nil {
		return nil, apperror.Internal("failed to upload attachment", err)
	}

	attachment := &model.Attachment{
		Name:     fileName,
		URL:      fileURL,
		LessonID: input.LessonID,
	}
	if err := s.attachments.Create(attachment); err != nil {
		return nil, apperror.Internal("failed to save attachment", err)
	}

	return attachment, nil
}

func (s *attachmentService) Download(ctx context.Context, id uint, roles []string) (*AttachmentDownload, error) {
	attachment, err := s.attachments.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("attachment not found")
		}
		return nil, apperror.Internal("failed to get attachment", err)
	}

	if err := s.ensureLessonAccess(attachment.LessonID, roles); err != nil {
		return nil, err
	}

	file, err := s.storage.Download(ctx, attachment.URL)
	if err != nil {
		return nil, apperror.Internal("failed to download attachment", err)
	}

	return &AttachmentDownload{
		Attachment:  attachment,
		File:        file,
		ContentName: attachment.Name,
	}, nil
}

func (s *attachmentService) ensureLessonAccess(lessonID uint, roles []string) error {
	if _, err := s.lessons.GetByID(lessonID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("lesson not found")
		}
		return apperror.Internal("failed to get lesson", err)
	}

	if len(roles) == 0 {
		return apperror.Forbidden("lesson access required")
	}

	return nil
}

func sanitizeFileName(name string) string {
	name = strings.TrimSpace(filepath.Base(name))
	if name == "." || name == string(filepath.Separator) {
		return ""
	}
	if len(name) > maxAttachmentNameLength {
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		maxBase := maxAttachmentNameLength - len(ext)
		if maxBase < 1 {
			return name[:maxAttachmentNameLength]
		}
		return base[:maxBase] + ext
	}
	return name
}
