package service

import (
	"context"
	"core/config"
	"core/internal/dto"
	"core/internal/repo"
	"core/pkg/response"
	"net/http"
)

type BookService struct {
	cfg      *config.Config
	bookRepo *repo.BookRepository
}

func NewBookService(cfg *config.Config, bookRepo *repo.BookRepository) *BookService {
	return &BookService{cfg: cfg, bookRepo: bookRepo}
}

func (s *BookService) GetList(ctx context.Context, query dto.GetBooksQuery) (*dto.GetBooksResponse, *response.HTTPError) {
	if query.Limit <= 0 || query.Limit > 100 {
		query.Limit = 10
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	result, err := s.bookRepo.GetList(ctx, query)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGetBooks,
				Detail: err.Error(),
			},
		}
	}

	return result, nil
}

func (s *BookService) GetByID(ctx context.Context, id uint, query dto.GetLessonsQuery) (*dto.GetLessonsResponse, *response.HTTPError) {
	if query.Limit <= 0 || query.Limit > 100 {
		query.Limit = 10
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	result, err := s.bookRepo.GetByID(ctx, id, query)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGetLessons,
				Detail: err.Error(),
			},
		}
	}

	return result, nil
}

func (s *BookService) GetLessonByID(ctx context.Context, id uint, query dto.GetQuestionsQuery) (*dto.GetQuestionsResponse, *response.HTTPError) {
	if query.Limit <= 0 || query.Limit > 100 {
		query.Limit = 10
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	result, err := s.bookRepo.GetLessonByID(ctx, id, query)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGetQuestions,
				Detail: err.Error(),
			},
		}
	}

	return result, nil
}
