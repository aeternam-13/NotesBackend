package models

type Note struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Color     int    `json:"color"`
}
