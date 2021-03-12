package model

type DBUrl struct {
	Id       int64  `json:"id"`
	LongURL  string `json:"longURL`
	ShortURL string `json:"shortURL"`
}
