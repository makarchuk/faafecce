package faafecce

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"gocv.io/x/gocv"
	// "gocv.io/x/gocv"
)

//Middler is a function for finding mirroring axis
type Middler func(image.Image) (int, error)

//Middle returns x-axis middle of the picture given the picture
func Middle(img image.Image) (int, error) {
	rect := img.Bounds()
	return (rect.Max.X - rect.Min.X) / 2, nil
}

//Face returns x-axis middle of the face found on the picture
func Face(img image.Image) (int, error) {
	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()
	if !classifier.Load("data/haarcascade_frontalface_default.xml") {
		return 0, fmt.Errorf("can't read classifier file")
	}
	b := bytes.NewBuffer([]byte{})
	jpeg.Encode(b, img, &jpeg.Options{Quality: 100})
	mat, err := gocv.IMDecode(b.Bytes(), gocv.IMReadColor)
	if err != nil {
		return 0, err
	}
	rects := classifier.DetectMultiScale(mat)
	if len(rects) == 0 {
		return 0, fmt.Errorf("no faces found on this picture")
	}
	main := mainFace(rects)
	return (main.Max.X + main.Min.X) / 2, nil
}

//Just choose the bigest face there is
func mainFace(rects []image.Rectangle) image.Rectangle {
	maxArea := 0
	var current *image.Rectangle
	for _, rect := range rects {
		area := (rect.Max.Y - rect.Min.Y) * (rect.Max.X - rect.Min.X)
		if area > maxArea {
			current = &rect
		}
	}
	return *current
}

func loadImage(file string) (image.Image, error) {
	infile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer infile.Close()
	img, _, err := image.Decode(infile)
	return img, err
}

//Transform actually does the job of transforming the picture
func Transform(mid Middler, file string, name string, outfile string) error {
	img, err := loadImage(file)
	out, err := os.Create(outfile)
	defer out.Close()
	if err != nil {
		return fmt.Errorf("can't open file for writing")
	}
	middle, err := mid(img)
	if err != nil {
		return fmt.Errorf("Problem with middle detection %v", err)
	}
	dst := mirroredImage(img, middle)
	jpeg.Encode(out, dst, &jpeg.Options{Quality: 100})
	return nil
}

func mirroredImage(img image.Image, middle int) *image.NRGBA {
	maxX := img.Bounds().Max.X
	fmt.Println(img.Bounds())
	rect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{maxX * 2, img.Bounds().Max.Y},
	}
	dst := image.NewNRGBA(rect)
	fmt.Println(dst.Bounds())
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			color := img.At(x, y)
			if x == middle {
				dst.Set(2*middle+maxX-x, y, color)
				dst.Set(x, y, color)
			} else if x < middle {
				dst.Set(x, y, color)
				dst.Set(2*middle-x, y, color)
			} else {
				fromRight := maxX - x
				dst.Set(2*middle+fromRight, y, color)
				dst.Set(2*maxX-fromRight, y, color)
			}
		}
	}
	return dst
}
