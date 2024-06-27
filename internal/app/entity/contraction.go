package entity

type Contraction struct {
	ID          int    `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}
