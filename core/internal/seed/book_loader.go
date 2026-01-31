package seed

import (
	"encoding/json"
	"fmt"
	"os"

	"core/internal/model"

	"gorm.io/gorm"
)

func LoadBookFile(db *gorm.DB, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var data JsonBook
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		book := model.Book{
			Title: data.Title,
		}

		if err := tx.Create(&book).Error; err != nil {
			return err
		}
		fmt.Printf("Seeding book %s\n", book.Title)

		for _, session := range data.Sessions {
			for _, lessonJson := range session.Lessons {
				lesson := model.Lesson{
					Title:       lessonJson.Title,
					Description: lessonJson.Description,
					IsVideo:     lessonJson.LessonDetail.IsVideo,
					AudioURL:    lessonJson.LessonDetail.FullAudioUrl,
					BookID:      book.ID,
				}
				if err := tx.Create(&lesson).Error; err != nil {
					return err
				}

				var questions []model.Question
				for i, questionJson := range lessonJson.LessonDetail.Questions {
					questions = append(questions, model.Question{
						Content:   questionJson.Content,
						TimeStart: questionJson.TimeStart,
						TimeEnd:   questionJson.TimeEnd,
						Order:     i,
						LessonID:  lesson.ID,
					})
				}
				if len(questions) > 0 {
					if err := tx.Create(&questions).Error; err != nil {
						return err
					}
				}
				fmt.Printf("Adding lesson %s (%d sentences)\n", lesson.Title, len(questions))
			}
		}
		return nil
	})
}
