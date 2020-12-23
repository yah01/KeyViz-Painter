package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	BytesUpperBound = 375 * 1024
)

var (
	db          *gorm.DB
	grayImage   *image.Gray
	imageWidth  int
	imageHeight int

)

func init() {
	dsn := "root@tcp(127.0.0.1:4000)/KeyViz"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	imgFile, err := os.Open("s.jpg")
	if err != nil {
		panic(err)
	}
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	grayImage = image.NewGray(img.Bounds())
	draw.Draw(grayImage, img.Bounds(), img, image.ZP, draw.Src)

	imageWidth, imageHeight = grayImage.Rect.Dx(), grayImage.Rect.Dy()
}

func rowTable(r int) string {
	return fmt.Sprintf("pixel%d", imageHeight-r-1)
}

func createPixels() {
	for r := 0; r < imageHeight; r++ {
		sqlCmd := fmt.Sprintf("create table %s (\n    data longblob null\n)", rowTable(r))
		if err := db.Exec(sqlCmd).Error; err != nil {
			panic(err)
		}
	}
}

func drawColumnAt(col int) {
	// var err error
	for r := 0; r < imageHeight; r++ {
		fmt.Printf("%v-th row drawing\n", r)
		grayScale := int(grayImage.GrayAt(col, r).Y)
		bytesLen := grayScale * BytesUpperBound / 255

		pixelData := PixelData{
			Data: make([]byte, bytesLen),
		}
		if err := db.Table(rowTable(imageHeight-r-1)).Create(&pixelData).Error; err != nil {
			fmt.Printf("err=%+v", err)
			panic(err)
		}
	}
}

func main() {
	// createPixels()
	for c := 0; c < imageWidth; c++ {
		timer := time.NewTimer(time.Minute)
		drawColumnAt(c)
		fmt.Printf("%v-th column drawing done\n", c)
		<-timer.C
	}
}
