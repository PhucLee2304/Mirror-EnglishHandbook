package dto

type GetBooksQuery struct {
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int    `form:"offset" binding:"omitempty,min=0"`
	Query  string `form:"query" binding:"omitempty,max=100"`
}

type BookBase struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type GetBooksResponse struct {
	Books []BookBase `json:"books"`
	Total int64      `json:"total"`
}

type BookUri struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type GetLessonsQuery struct {
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int    `form:"offset" binding:"omitempty,min=0"`
	Query  string `form:"query" binding:"omitempty,max=100"`
}

type LessonBase struct {
	BookBase
	ID          uint   `json:"id"`
	Title       string `json:"lessonTitle"`
	Description string `json:"description"`
	IsVideo     bool   `json:"isVideo"`
	AudioURL    string `json:"audioUrl"`
}

type GetLessonsResponse struct {
	Lessons []LessonBase `json:"lessons"`
	Total   int64        `json:"total"`
}

type LessonUri struct {
	BookID   uint `uri:"id" binding:"required,min=1"`
	LessonID uint `uri:"lessonID" binding:"required,min=1"`
}

//type GetQuestionsQuery struct {
//	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
//	Offset int    `form:"offset" binding:"omitempty,min=0"`
//	Query  string `form:"query" binding:"omitempty,max=100"`
//}

type QuestionBase struct {
	ID        uint    `json:"id"`
	Content   string  `json:"content"`
	TimeStart float64 `json:"timeStart"`
	TimeEnd   float64 `json:"timeEnd"`
	Order     int     `json:"order"`
}

type GetQuestionsResponse struct {
	LessonBase
	Questions []QuestionBase `json:"questions"`
	Duration  float64        `json:"duration"`
}
