package inout

//go:generate easyjson -disallow_unknown_fields -all ./input.go

type InputURL struct {
	URL string `json:"url"`
}

type ExternalInput struct {
	ExternalID  string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

//easyjson:json
type ExternalBatchInput []ExternalInput
