package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"core/internal/model"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func LoadBookFile(db *gorm.DB, path string, minioClient *minio.Client, baseDir string, bucketName string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var data JsonBook
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	ctx := context.Background()
	exists, errBucket := minioClient.BucketExists(ctx, bucketName)
	if errBucket == nil && !exists {
		_ = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		log.Printf("Bucket created successfully: %s\n", bucketName)
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

				if lesson.AudioURL != "" {
					localFilePath := filepath.Join(baseDir, lesson.AudioURL)
					objectName := lesson.AudioURL

					if _, err := os.Stat(localFilePath); err == nil {
						log.Printf("Uploading %s to MinIO...\n", objectName)
						_, errUpload := minioClient.FPutObject(ctx, bucketName, objectName, localFilePath, minio.PutObjectOptions{
							ContentType: "video/mp4",
						})
						if errUpload != nil {
							log.Printf("Error uploading %s to MinIO: %v\n", objectName, errUpload)
						} else {
							log.Printf("Successfully uploaded %s to MinIO.\n", objectName)
						}
					} else {
						log.Printf("Failed to upload %s to MinIO: %v\n", objectName, err)
					}
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
