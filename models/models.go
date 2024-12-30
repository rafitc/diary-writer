package models

// models to share across project
// to avoid circular imports
type ASSET struct {
	Asset     string `json:"asset"`
	Extension string `json:"extension"`
}

type DiaryEntry struct {
	Content string  `json:"content"`
	Date    string  `json:"date"`
	Asset   []ASSET `json:"asset"`
}
