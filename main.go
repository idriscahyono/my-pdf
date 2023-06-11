package main

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/google/uuid"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var fileType string
	var err error

	wd := workDir()
	inputPath := wd + "/input/docx_sample.docx"

	extension := filepath.Ext(inputPath)
	switch extension {
	case ".pdf":
		fileType = "pdf"
	case ".docx":
		fileType = "docx"
	case ".jpg", ".jpeg", ".png":
		fileType = "image"
	default:
		panic("Unknown file type")
	}

	if fileType == "pdf" {
		copyDst := "/output/"
		resCopyFile, err := copyFile(inputPath, copyDst)
		if err != nil {
			panic(fmt.Sprintf("err copyFile %s", err))
		}

		inputPath = *resCopyFile
	}

	if fileType == "image" {
		resImageToPdf, err := imageToPdf(inputPath)
		if err != nil {
			panic(fmt.Sprintf("err imageToPdf %s", err))
		}

		inputPath = *resImageToPdf
	}

	if fileType == "docx" {
		resDocxToPdf, err := docxToPdf(inputPath)
		if err != nil {
			panic(fmt.Sprintf("err docxToPdf %s", err))
		}

		inputPath = *resDocxToPdf
	}

	resWatermark, err := generateQrCode()
	if err != nil {
		panic(fmt.Sprintf("err generateQrCode %s", err))
	}

	err = addWatermarkPdf(inputPath, *resWatermark)
	if err != nil {
		panic(fmt.Sprintf("err addWatermarkPdf %s", err))
	}

	err = overrideMetadata(inputPath)
	if err != nil {
		fmt.Println("err overrideMetadata", err)
	}

	err = pdfThumbnail(inputPath)
	if err != nil {
		fmt.Println("err pdfThumbnail", err)
	}

	fmt.Println("done")
}

func copyFile(inputPath, outputPath string) (*string, error) {

	wd := workDir()
	fileName := filepath.Base(inputPath)
	outputPath = wd + outputPath + fileName

	source, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer source.Close()

	destination, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return nil, err
	}

	return &outputPath, nil
}

func overrideMetadata(inputPath string) error {
	fmt.Println("overrideMetadata")

	properties := []string{
		"Title=My PDF Idris",
		"Author=Idris",
	}
	command := "pdfcpu"
	args := []string{"properties", "add", inputPath}
	args = append(args, properties...)
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func pdfThumbnail(inputPath string) error {
	fmt.Println("pdfThumbnail")

	wd := workDir()
	outputPath := wd + "/output/thumbnail/"

	fileName := filepath.Base(inputPath)
	outputPath = outputPath + strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "-thumb" + ".jpg"

	command := "convert"
	args := []string{"-density", "300", "-trim", inputPath + "[0]", "-resize", "720x360^", "-extent", "720x360", "-define", "jpg:extent=100KB", "-strip", "-trim", "+repage", outputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func docxToPdf(inputPath string) (*string, error) {
	fmt.Println("docxToPdf")

	wd := workDir()
	fileName := filepath.Base(inputPath)
	outputPath := wd + "/output/"

	command := "soffice"
	args := []string{"--headless", "--convert-to", "pdf", "--outdir", outputPath, inputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	output := outputPath + strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".pdf"
	return &output, nil
}

func imageToPdf(inputPath string) (*string, error) {
	fmt.Println("imageToPdf")

	wd := workDir()
	fileName := filepath.Base(inputPath)
	outputPath := wd + "/output/" + strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".pdf"

	command := "convert"
	args := []string{inputPath, "-page", fmt.Sprintf("%dx%d", 595, 842), outputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return &outputPath, nil
}

func generateQrCode() (*string, error) {
	fmt.Println("generateQrCode")

	wd := workDir()
	overlayPath := wd + "/assets/overlay.png"
	outputQrPath := wd + "/output/qrcode/qrcode.png"

	content := fmt.Sprintf("https://privy.id/verify/%s", uuid.New())

	qrCode, err := qr.Encode(content, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	qrCode, err = barcode.Scale(qrCode, 100, 100)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(outputQrPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	png.Encode(file, qrCode)

	res, err := addLogoQrCode(outputQrPath, overlayPath)
	if err != nil {
		return nil, err
	}

	return res, err
}

func addLogoQrCode(inputPath, overlayPath string) (*string, error) {
	fmt.Println("addLogoQrCode")

	command := "convert"
	args := []string{inputPath, "-gravity", "center", overlayPath, "-composite", inputPath}
	cmd := exec.Command(command, args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return &inputPath, err
}

func addWatermarkPdf(inputPath, watermarkImage string) error {
	fmt.Println("addWatermarkImage")

	pos := []string{"sc:1 abs, rot:0, pos:br", "sc:1 abs, rot:0, pos:bl", "sc:1 abs, rot:0, pos:tl", "sc:1 abs, rot:0, pos:tr"}
	command := "pdfcpu"
	for _, p := range pos {
		args := []string{"watermark", "add", "-mode", "image", "--", watermarkImage, p, inputPath, inputPath}
		cmd := exec.Command(command, args...)

		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	}

	return nil
}

func workDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("err workDir %s", err))
	}

	return wd
}
