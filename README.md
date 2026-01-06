# English Dictionary Service

## Description
The **Dictionary Service** is the authoritative source of truth for all lexical data within the English Handbook ecosystem. It abstracts the complexity of linguistic databases and provides a clean, fast API for querying word information.

## Responsibilities
- **Word Lookup**: Retrieval of definitions, parts of speech, and usage examples.
- **Phonetics**: IPA transcriptions and audio references.
- **Thesaurus**: Synonyms, antonyms, and related terms.

## API Overview
This service exposes RESTful endpoints primarily consumed by internal services and the API Gateway.

### Key Endpoints
- `GET /api/v1/words/{word}`: functional details for a specific lexical entry.
- `GET /api/v1/search?q={query}`: Fuzzy search for auto-completion and suggestion.

## Local Development

```bash
# Run the service locally on port 8081 (default)
go run main.go
```

Ensure the database environment variables are set correctly before starting the service.
