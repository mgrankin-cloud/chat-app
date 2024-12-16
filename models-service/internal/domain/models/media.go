package models

type Media struct {
	ID 		 int64
	Data     []byte
	FileName string
	MimeType string
}