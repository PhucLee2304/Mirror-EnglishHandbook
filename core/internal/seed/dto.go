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
