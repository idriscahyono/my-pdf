package main

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"os"
	"os/exec"
)

func main() {
	wd, _ := workDir()
	//inputPath := wd + "/input/my-sample.pdf"
	//outputPath := wd + "/output/my-sample-output.pdf"
	//watermarkPath := wd + "/assets/qrcode.png"
	//
	//resWatermark, err := generateQrCode(watermarkPath)
	//if err != nil {
	//	fmt.Println("err generateQrCode", err)
	//}
	//
	//err = addWatermarkImage(inputPath, outputPath, resWatermark)
	//if err != nil {
	//	fmt.Println("err addWatermarkImage", err)
	//}

	//inputPath := wd + "/input/logo-privy.png"
	//outputPath := wd + "/output/output.pdf"
	//err := imageToPdf(inputPath, outputPath)
	//if err != nil {
	//	fmt.Println("err imageToPdf", err)
	//}

	//inputPath := wd + "/input/sample_doc.docx"
	//outputPath := wd + "/output/"
	//err := docxToPdf(inputPath, outputPath)
	//if err != nil {
	//	fmt.Println("err docxToPdf", err)
	//}

	inputPath := wd + "/output/sample_doc.pdf"
	outputPath := wd + "/output/thumbnail/tm.jpg"
	err := pdfThumbnail(inputPath, outputPath)
	if err != nil {
		fmt.Println("err pdfThumbnail", err)
	}

	fmt.Println("done")
}

func pdfThumbnail(inputPath, outputPath string) error {
	fmt.Println("pdfThumbnail")
	command := "pdfcpu"
	args := []string{"extract", "images", "-verbose=false", "-outdir", outputPath, inputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("line 60", err)
		return err
	}

	return nil
}

func docxToPdf(inputPath, outputPath string) error {
	fmt.Println("docxToPdf")
	//absPath, err := filepath.Abs(inputPath)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//_ = absPath

	command := "soffice"
	args := []string{"--headless", "--convert-to", "pdf", "--outdir", outputPath, inputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func imageToPdf(inputPath, outputPath string) error {
	fmt.Println("imageToPdf")
	command := "convert"
	args := []string{inputPath, outputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
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
