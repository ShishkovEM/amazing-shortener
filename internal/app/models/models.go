package models

// Link Структура записи информации о гиперссылках
type Link struct {
	ID        int    `json:"id"`        // Идентификатор гиперссылки
	Original  string `json:"original"`  // Исходная (длинная) ссылка
	Short     string `json:"short"`     // Короткая ссылка
	UserID    uint32 `json:"userID"`    // Индентификатор пользователя, создавшего ссылку
	IsDeleted bool   `json:"isDeleted"` // Пометка об удалении ссылки
}

type OriginalURL struct {
	OriginalURL string
	IsDeleted   bool
}

type User struct {
	ID      uint32
	Urls    []Link
	Session *Session
}

type Session struct {
	ID        string
	Signature string
}
