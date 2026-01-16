package seed

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"core/internal/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func LoadFile(db *gorm.DB, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var data map[string][]JsonWord
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	for _, words := range data {
		// fmt.Printf("Seeding %d words for prefix '%s'\n", len(words), prefix)
		for _, jw := range words {
			if err := insertWord(db, jw); err != nil {
				return fmt.Errorf("failed to insert word '%s': %w\n", jw.Word, err)
			}
		}
	}
	return nil
}

func insertWord(db *gorm.DB, jw JsonWord) error {
	hash, err := hashJsonWord(jw)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		word := model.Word{
			Word:        jw.Word,
			SourceUrls:  toJSON(jw.SourceUrls),
			ContentHash: hash,
		}

		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "word"},
				{Name: "content_hash"},
			},
			DoNothing: true,
		}).Create(&word)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return nil
		}

		for i, jp := range jw.Phonetics {
			phonetic := model.Phonetic{
				Text:      jp.Text,
				Audio:     jp.Audio,
				SourceUrl: jp.SourceUrl,
				Order:     i,
				WordID:    word.ID,
			}
			if err := tx.Create(&phonetic).Error; err != nil {
				return err
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
				definition := model.Definition{
					DefinitionText: jd.DefinitionText,
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
	})
}

func toJSON(v any) datatypes.JSON {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal to JSON: %v", err))
	}
	return datatypes.JSON(bytes)
}

func hashJsonWord(jw JsonWord) (string, error) {
	h := sha256.New()

	write := func(s string) {
		h.Write([]byte(s))
		h.Write([]byte{0})
	}

	write(jw.Word)

	for _, p := range jw.Phonetics {
		write(p.Text)
		if p.Audio != nil {
			write(*p.Audio)
		}
	}

	for _, m := range jw.Meanings {
		write(m.PartOfSpeech)
		for _, d := range m.Definitions {
			write(d.DefinitionText)
			if d.Example != nil {
				write(*d.Example)
			}
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
