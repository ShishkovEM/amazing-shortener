package models

// Link Структура записи информации о гиперссылках
type Link struct {
	ID       int    `json:"id"`       // Идентификатор гиперссылки
	Original string `json:"original"` // Исходная (длинная) ссылка
	Short    string `json:"short"`    // Короткая ссылка
	UserID   uint64 `json:"userID"`
}

type User struct {
	ID      uint64
	Urls    []Link
	Session *Session
}

type Session struct {
	ID        string
	Signature string
}
