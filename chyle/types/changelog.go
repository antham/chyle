package types

// Changelog represents a changelog entry with datas extracted and metadatas
type Changelog struct {
	Datas     []map[string]any `json:"datas"`
	Metadatas map[string]any   `json:"metadatas"`
}
