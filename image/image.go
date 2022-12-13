package main

import (
	"bytes"
	"gitee.com/kayi-cloud/gdimse/lib/go-dicom"
	"github.com/linxlib/logs"
	_ "github.com/linxlib/openjpeg"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {

	path := "/Users/linx/Downloads/I00000010082.dcm"
	ds, e := dicom.ReadDataSetFromFile(path, dicom.ReadOptions{})
	if e != nil {
		logs.Error(e)
		return
	}
	//logs.Infoln(ds)
	//return
	//pi, err := ds.GetPixelDataInfo()
	//pi, err := ds.GetPixelDataInfo()
	pi, _, err := ds.ParsePixelImage(0, true)
	if err != nil {
		logs.Error(err)
		return
	}
	//level := int16(400)
	//width := int16(380)
	frames := pi.Frames
	if len(frames) == 0 {
		panic("No images found")
	} else if len(frames) > 1 {
		log.Println("Many images found, displaying only first element")
	}

	frame := frames[0].EncapsulatedData.Data
	//dicomjp2,_ := os.Create("dicom.jp2")
	os.WriteFile("dicom.jp2", frame, os.ModePerm)
	logs.Info("解码Dicom jp2数据")

	in1 := bytes.NewReader(frame)

	img1, _, err := image.Decode(in1)
	pngFile, _ := os.Create("dicom.png")
	err = png.Encode(pngFile, img1)
	if err != nil {
		logs.Error(err)
		return
	}

	//bs, _ := os.ReadFile("file1.jp2")
	//f, _ := os.Open("/Users/linx/Downloads/file1.jp2")
	//img, err := jpeg2000.Decode(f)
	//if err != nil {
	//	logs.Error(err)
	//	return
	//}
	//w, _ := os.Create("1.jpg")
	//err = jpeg.Encode(w, img, nil)
	//if err != nil {
	//	logs.Error(err)
	//	return
	//}
	logs.Info("解码jp2文件")
	in, err := os.Open("jp2000.jp2")
	if err != nil {
		logs.Error(err)
		return
	}
	img, inFmt, err := image.Decode(in)

	defer in.Close()
	logs.Printf("Decoded %s: %dx%d %s\n", in.Name(), img.Bounds().Dx(), img.Bounds().Dy(), inFmt)
	f, _ := os.Create("1.png")
	err = png.Encode(f, img)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Error("解码jp2文件成功，转为png成功，请查看1.png")
}
