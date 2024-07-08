package inout

//go:generate easyjson -disallow_unknown_fields -all ./output.go

type OutputURL struct {
	Result string `json:"result"`
}

type ExternalOutput struct {
	ExternalID string `json:"correlation_id"`
	ShortURL   string `json:"short_url"`
}

//easyjson:json
type ExternalBatchOutput []ExternalOutput
