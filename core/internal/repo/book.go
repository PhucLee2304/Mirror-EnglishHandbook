package repo

import (
	"context"
	"core/internal/dto"
	"core/internal/model"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) WithTx(tx *gorm.DB) *BookRepository {
	return &BookRepository{db: tx}
}

func (r *BookRepository) GetList(ctx context.Context, query dto.GetBooksQuery) (*dto.GetBooksResponse, error) {
	var books []model.Book
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Book{})

	if query.Query != "" {
		db = db.Where("book ILIKE ?", query.Query+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	err := db.Limit(query.Limit).Offset(query.Offset).Order("id ASC").Find(&books).Error
	if err != nil {
		return nil, err
	}

	bookResps := make([]dto.BookBase, 0, len(books))
	for _, book := range books {
		bookResps = append(bookResps, dto.BookBase{
			ID:    book.ID,
			Title: book.Title,
		})
	}

	return &dto.GetBooksResponse{
		Books: bookResps,
		Total: total,
	}, nil
}

func (r *BookRepository) GetByID(ctx context.Context, id uint, query dto.GetLessonsQuery) (*dto.GetLessonsResponse, error) {
	var book model.Book
	if err := r.db.WithContext(ctx).First(&book, id).Error; err != nil {
		return nil, err
	}

	var lessons []model.Lesson
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Lesson{}).Where("book_id = ?", id)
	if query.Query != "" {
		db = db.Where("title ILIKE ?", query.Query+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Limit(query.Limit).Offset(query.Offset).Order("id ASC").Find(&lessons).Error; err != nil {
		return nil, err
	}

	resp := make([]dto.LessonBase, 0, len(lessons))
	for _, lesson := range lessons {
		resp = append(resp, dto.LessonBase{
			BookBase: dto.BookBase{
				ID:    book.ID,
				Title: book.Title,
			},
			ID:          lesson.ID,
			Title:       lesson.Title,
			Description: lesson.Description,
			IsVideo:     lesson.IsVideo,
			AudioURL:    lesson.AudioURL,
		})
	}

	return &dto.GetLessonsResponse{
		Lessons: resp,
		Total:   total,
	}, nil
}

func (r *BookRepository) GetLessonByID(ctx context.Context, id uint) (*dto.GetQuestionsResponse, error) {
	var lesson model.Lesson
	if err := r.db.WithContext(ctx).Preload("Book").First(&lesson, id).Error; err != nil {
		return nil, err
	}

	var questions []model.Question
	//var total int64

	//db := r.db.WithContext(ctx).Model(&model.Question{}).Where("lesson_id = ?", id)
	//if query.Query != "" {
	//	db = db.Where("content ILIKE ?", query.Query+"%")
	//}

	//if err := db.Count(&total).Error; err != nil {
	//	return nil, err
	//}

	//if err := db.Limit(query.Limit).Offset(query.Offset).Order("\"order\" ASC").Find(&questions).Error; err != nil {
	//	return nil, err
	//}
	if err := r.db.WithContext(ctx).
		Where("lesson_id = ?", id).
		Order("\"order\" ASC").
		Find(&questions).Error; err != nil {
		return nil, err
	}

	lessonBase := dto.LessonBase{
		BookBase: dto.BookBase{
			ID:    lesson.Book.ID,
			Title: lesson.Book.Title,
		},
		ID:          lesson.ID,
		Title:       lesson.Title,
		Description: lesson.Description,
		IsVideo:     lesson.IsVideo,
		AudioURL:    lesson.AudioURL,
	}

	resp := make([]dto.QuestionBase, 0, len(questions))
	for _, question := range questions {
		resp = append(resp, dto.QuestionBase{
			ID:        question.ID,
			Content:   question.Content,
			TimeStart: question.TimeStart,
			TimeEnd:   question.TimeEnd,
			Order:     question.Order,
		})
	}

	var duration float64
	if len(questions) > 0 {
		lastQuestion := questions[len(questions)-1]
		duration = lastQuestion.TimeEnd
	}

	return &dto.GetQuestionsResponse{
		LessonBase: lessonBase,
		Questions:  resp,
		Duration:   duration,
	}, nil
}
