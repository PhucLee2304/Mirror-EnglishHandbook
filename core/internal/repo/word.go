package repo

import (
	"context"
	"core/internal/dto"
	"core/internal/model"
	"core/pkg/utils/slicex"
	"strings"

	"gorm.io/gorm"
)

type WordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) WithTx(tx *gorm.DB) *WordRepository {
	return &WordRepository{db: tx}
}

func (r *WordRepository) GetByID(ctx context.Context, id uint) (*dto.WordBase, error) {
	var word model.Word
	err := r.db.WithContext(ctx).
		Preload("Phonetics").
		Preload("Meanings.Definitions").
		First(&word, id).Error
	if err != nil {
		return nil, err
	}

	response := &dto.WordBase{
		ID:         word.ID,
		Word:       word.Word,
		SourceUrls: slicex.JsonToSlice(word.SourceUrls),
		Phonetics:  make([]dto.PhoneticBase, 0),
		Meanings:   make([]dto.MeaningBase, 0),
	}
	for _, p := range word.Phonetics {
		response.Phonetics = append(response.Phonetics, dto.PhoneticBase{
			ID:        p.ID,
			Text:      p.Text,
			Audio:     p.Audio,
			SourceUrl: p.SourceUrl,
		})
	}
	for _, m := range word.Meanings {
		meaningDTO := dto.MeaningBase{
			ID:           m.ID,
			PartOfSpeech: m.PartOfSpeech,
			Synonyms:     slicex.JsonToSlice(m.Synonyms),
			Antonyms:     slicex.JsonToSlice(m.Antonyms),
			Definitions:  make([]dto.DefinitionBase, 0),
		}

		for _, d := range m.Definitions {
			meaningDTO.Definitions = append(meaningDTO.Definitions, dto.DefinitionBase{
				ID:         d.ID,
				Definition: d.DefinitionText,
				Example:    d.Example,
				Synonyms:   slicex.JsonToSlice(d.Synonyms),
				Antonyms:   slicex.JsonToSlice(d.Antonyms),
			})
		}
		response.Meanings = append(response.Meanings, meaningDTO)
	}

	return response, nil
}

func (r *WordRepository) GetList(ctx context.Context, query dto.GetWordsQuery) (*dto.GetWordsResponse, error) {
	var words []model.Word
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Word{})

	if query.Query != "" {
		db = db.Where("word ILIKE ?", query.Query+"%")
	}

	if query.Type != "" {
		db = db.Joins("JOIN meanings ON meanings.word_id = words.id").
			Where("meanings.part_of_speech = ?", query.Type).
			Group("words.id")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	err := db.Preload("Phonetics").
		Preload("Meanings").
		Limit(query.Limit).
		Offset(query.Offset).
		Order("word ASC").
		Find(&words).Error

	if err != nil {
		return nil, err
	}

	response := make([]dto.WordResponse, 0)
	for _, word := range words {
		item := dto.WordResponse{
			ID:   word.ID,
			Word: word.Word,
		}

		if len(word.Phonetics) > 0 {
			item.Phonetics = word.Phonetics[0].Text
			if word.Phonetics[0].Audio != nil {
				item.Audio = *word.Phonetics[0].Audio
			}
		}

		uniquePartOfSpeech := make(map[string]bool)
		var partOfSpeechList []string
		for _, meaning := range word.Meanings {
			if !uniquePartOfSpeech[meaning.PartOfSpeech] {
				uniquePartOfSpeech[meaning.PartOfSpeech] = true
				partOfSpeechList = append(partOfSpeechList, meaning.PartOfSpeech)
			}
		}
		if len(partOfSpeechList) > 0 {
			item.PartOfSpeech = strings.Join(partOfSpeechList, ", ")
		}
		response = append(response, item)
	}

	return &dto.GetWordsResponse{
		Words: response,
		Total: total,
	}, nil
}
