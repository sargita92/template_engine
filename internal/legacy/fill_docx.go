package legacy

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// FillDocxBytes - Preenche um DOCX a partir de bytes e aplica ProcessPlaceholders em document.xml
func FillDocxBytes(inputData []byte, outputPath string, replaces map[string]string) error {
	reader, err := zip.NewReader(bytes.NewReader(inputData), int64(len(inputData)))
	if err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}

		// Processa apenas o documento principal
		if file.Name == "word/document.xml" {
			text := string(content)
			text = ProcessPlaceholders(text, replaces)
			content = []byte(text)
		}

		header := &zip.FileHeader{
			Name:   file.Name,
			Method: zip.Deflate,
		}
		header.SetMode(0666)

		w, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = w.Write(content)
		if err != nil {
			return err
		}
	}

	return nil
}

// FillDocx - Preenche um DOCX a partir de um arquivo no disco
func FillDocx(inputPath, outputPath string, replaces map[string]string) error {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(absInput)
	if err != nil {
		return err
	}

	return FillDocxBytes(data, outputPath, replaces)
}
