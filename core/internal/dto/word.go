package dto

type WordUri struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type DefinitionBase struct {
	ID         uint     `json:"id"`
	Definition string   `json:"definition"`
	Example    *string  `json:"example,omitempty"`
	Synonyms   []string `json:"synonyms,omitempty"`
	Antonyms   []string `json:"antonyms,omitempty"`
}

type MeaningBase struct {
	ID           uint             `json:"id"`
	PartOfSpeech string           `json:"partOfSpeech"`
	Synonyms     []string         `json:"synonyms,omitempty"`
	Antonyms     []string         `json:"antonyms,omitempty"`
	Definitions  []DefinitionBase `json:"definitions"`
}

type PhoneticBase struct {
	ID        uint    `json:"id"`
	Text      string  `json:"text"`
	Audio     *string `json:"audio,omitempty"`
	SourceUrl *string `json:"sourceUrl,omitempty"`
}

type WordBase struct {
	ID         uint           `json:"id"`
	Word       string         `json:"word"`
	Phonetics  []PhoneticBase `json:"phonetics"`
	Meanings   []MeaningBase  `json:"meanings"`
	SourceUrls []string       `json:"sourceUrls,omitempty"`
}

type GetWordResponse struct {
	WordBase
}

type GetWordsQuery struct {
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int    `form:"offset" binding:"omitempty,min=0"`
	Type   string `form:"type" binding:"omitempty"`
	Query  string `form:"query" binding:"omitempty,max=100"`
}

type GetWordsResponse struct {
	Words []WordResponse `json:"words"`
	Total int64          `json:"total"`
}

type WordResponse struct {
	ID           uint   `json:"id"`
	Word         string `json:"word"`
	Phonetics    string `json:"phonetic"`
	Audio        string `json:"audio,omitempty"`
	PartOfSpeech string `json:"part_of_speech"`
}
