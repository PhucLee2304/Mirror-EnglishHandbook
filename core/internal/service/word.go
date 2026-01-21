package service

import (
	"context"
	"core/config"
	"core/internal/dto"
	"core/internal/repo"
	"core/pkg/response"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type WordService struct {
	cfg      *config.Config
	wordRepo *repo.WordRepository
}

func NewWordService(cfg *config.Config, wordRepo *repo.WordRepository) *WordService {
	return &WordService{
		cfg:      cfg,
		wordRepo: wordRepo,
	}
}

func (s *WordService) GetByID(ctx context.Context, id uint) (*dto.WordBase, *response.HTTPError) {
	dto, err := s.wordRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &response.HTTPError{
				StatusCode: http.StatusNotFound,
				Error: response.ErrorResponse{
					Error:  response.MessageCodeNotFound,
					Detail: err.Error(),
				},
			}
		}
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGetWords,
				Detail: err.Error(),
			},
		}
	}
	return dto, nil
}

func (s *WordService) GetList(ctx context.Context, query dto.GetWordsQuery) (*dto.GetWordsResponse, *response.HTTPError) {
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	result, err := s.wordRepo.GetList(ctx, query)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGetWords,
				Detail: err.Error(),
			},
		}
	}

	return result, nil
}
