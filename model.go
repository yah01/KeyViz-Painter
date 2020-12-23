package main

type Bitmap struct {
	Id          int64 `gorm:"primaryKey"`
	Description string
}

type PixelData struct {
	Data []byte
}
