package seed

type JsonWord struct {
	Word       string         `json:"word"`
	SourceUrls []string       `json:"sourceUrls"`
	Phonetics  []JsonPhonetic `json:"phonetics"`
	Meanings   []JsonMeaning  `json:"meanings"`
}

type JsonPhonetic struct {
	Text      string  `json:"text"`
	Audio     *string `json:"audio"`
	SourceUrl *string `json:"sourceUrl"`
}

type JsonMeaning struct {
	PartOfSpeech string           `json:"partOfSpeech"`
	Synonyms     []string         `json:"synonyms"`
	Antonyms     []string         `json:"antonyms"`
	Definitions  []JsonDefinition `json:"definitions"`
}

type JsonDefinition struct {
	Definition string   `json:"definition"`
	Example    *string  `json:"example"`
	Antonyms   []string `json:"antonyms"`
	Synonyms   []string `json:"synonyms"`
}

type JsonBook struct {
	Title    string        `json:"title"`
	Sessions []JsonSession `json:"sessions"`
}

type JsonSession struct {
	Title   string       `json:"title"`
	Lessons []JsonLesson `json:"lessons"`
}

type JsonLesson struct {
	Title        string           `json:"title"`
	Subtitle     string           `json:"subtitle"`
	Description  string           `json:"description"`
	LessonDetail JsonLessonDetail `json:"lessionDetail"`
}

type JsonLessonDetail struct {
	IsVideo      bool           `json:"isVideo"`
	FullAudioUrl string         `json:"fullAudioUrl"`
	Questions    []JsonQuestion `json:"questions"`
}

type JsonQuestion struct {
	Content   string  `json:"content"`
	TimeStart float64 `json:"timeStart"`
	TimeEnd   float64 `json:"timeEnd"`
}
