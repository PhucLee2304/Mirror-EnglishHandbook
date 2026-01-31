package seed

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"core/internal/model"
	"core/pkg/utils/slicex"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func LoadHandbookFile(db *gorm.DB, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var data map[string][]JsonWord
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, words := range data {
			for _, jw := range words {
				if err := insertWord(tx, jw); err != nil {
					fmt.Printf("failed to insert word %s: %v\n", jw.Word, err)
				}
			}
		}
		return nil
	})
}

func insertWord(tx *gorm.DB, jw JsonWord) error {
	var word model.Word
	err := tx.Where("word = ?", jw.Word).First(&word).Error
	if err == nil {
		existingUrls := slicex.JsonToSlice(word.SourceUrls)
		newUrls := addUrl(existingUrls, jw.SourceUrls)
		word.SourceUrls = toJSON(newUrls)
		if err := tx.Save(&word).Error; err != nil {
			return err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		word = model.Word{
			Word:       jw.Word,
			SourceUrls: toJSON(toJSON(jw.SourceUrls)),
		}
		if err := tx.Create(&word).Error; err != nil {
			return err
		}
	} else {
		return err
	}

	for i, jp := range jw.Phonetics {
		if jp.Text == "" && jp.Audio == nil {
			continue
		}

		var count int64
		query := tx.Model(&model.Phonetic{}).
			Where("word_id = ? AND text = ?", word.ID, jp.Text)
		if jp.Audio != nil {
			query = query.Where("audio = ?", jp.Audio)
		} else {
			query = query.Where("audio IS NULL")
		}
		query.Count(&count)

		if count == 0 {
			phonetic := model.Phonetic{
				Text:      jp.Text,
				Audio:     jp.Audio,
				SourceUrl: jp.SourceUrl,
				Order:     i + 1,
				WordID:    word.ID,
			}
			if err := tx.Create(&phonetic).Error; err != nil {
				return err
			}
		}
	}

	for _, jm := range jw.Meanings {
		meaning := model.Meaning{
			PartOfSpeech: jm.PartOfSpeech,
			Synonyms:     toJSON(jm.Synonyms),
			Antonyms:     toJSON(jm.Antonyms),
			WordID:       word.ID,
		}
		if err := tx.Create(&meaning).Error; err != nil {
			return err
		}

		for _, jd := range jm.Definitions {
			if jd.Definition == "" {
				continue
			}

			definition := model.Definition{
				DefinitionText: jd.Definition,
				Example:        jd.Example,
				Antonyms:       toJSON(jd.Antonyms),
				Synonyms:       toJSON(jd.Synonyms),
				MeaningID:      meaning.ID,
			}
			if err := tx.Create(&definition).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func addUrl(a, b []string) []string {
	check := make(map[string]bool)
	var res []string

	for _, val := range a {
		if _, ok := check[val]; !ok {
			check[val] = true
			res = append(res, val)
		}
	}

	for _, val := range b {
		if _, ok := check[val]; !ok {
			check[val] = true
			res = append(res, val)
		}
	}

	return res
}

func toJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON("[]")
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal to JSON: %v", err))
		return datatypes.JSON("[]")
	}
	if string(bytes) == "null" {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(bytes)
}
