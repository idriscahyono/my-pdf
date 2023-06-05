package main

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	wd, _ := workDir()
	inputPath := wd + "/input/my-sample.pdf"
	outputPath := wd + "/output/my-sample-output.pdf"
	watermarkPath := wd + "/assets/qrcode.png"

	resWatermark, err := generateQrCode(watermarkPath)
	if err != nil {
		fmt.Println("err generateQrCode", err)
	}

	err = addWatermarkImage(inputPath, outputPath, resWatermark)
	if err != nil {
		fmt.Println("err addWatermarkImage", err)
	}

	fmt.Println("done")
}

func generateQrCode(watermarkPath string) (string, error) {
	fmt.Println("generateQrCode")
	qrCode, err := qr.Encode("https://www.idriscahyono.com/", qr.M, qr.Auto)
	if err != nil {
		return "", err
	}

	qrCode, err = barcode.Scale(qrCode, 100, 100)
	if err != nil {
		return "", err
	}

	file, err := os.Create(watermarkPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	png.Encode(file, qrCode)

	return watermarkPath, err
}

func addWatermarkImage(inputPath, outputPath, watermarkImage string) error {
	fmt.Println("addWatermarkImage")
	command := "pdfcpu"
	args := []string{"watermark", "add", "-mode", "image", "--", watermarkImage, "sc:1 abs, rot:0, pos:br", inputPath, outputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func workDir() (string, error) {
	fmt.Println("getWorkDir")
	wd, err := os.Getwd()
	if err != nil {
		return wd, err
	}

	return wd, err
}
