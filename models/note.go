package models

// Note yapısı
type Note struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Content  string `json:"content"`
	FolderID *uint  `json:"folder_id,omitempty"` // uint olarak güncellendi, gorm ile uyumlu
	UserID   uint   `json:"user_id"`             // Kullanıcıya bağlamak için
}

// Folder yapısı
type Folder struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	UserID uint   `json:"user_id"` // Kullanıcıya bağlamak için
}

// API yanıt yapısı
type APIResponse struct {
	Folders []Folder `json:"folders"`
	Notes   []Note   `json:"notes"`
}
