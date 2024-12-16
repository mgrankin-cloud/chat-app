package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
)

func CompressFiles(archivePath string, files []string) error {
	archive, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("не удалось создать архив: %v", err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	for _, file := range files {
		err := addFileToZip(zipWriter, file)
		if err != nil {
			return fmt.Errorf("не удалось добавить файл %s в архив: %v", file, err)
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %v", filePath, err)
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о файле %s: %v", filePath, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("не удалось создать заголовок для файла %s: %v", filePath, err)
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("не удалось создать запись для файла %s в архиве: %v", filePath, err)
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return fmt.Errorf("не удалось скопировать содержимое файла %s в архив: %v", filePath, err)
	}

	return nil
}

func SendArchive(conn *websocket.Conn, userID int64, files []string) error {
	archivePath := fmt.Sprintf("%s/%d.zip", os.Getenv("GOPATH"), userID)

	err := CompressFiles(archivePath, files)
	if err != nil {
		return fmt.Errorf("ошибка создания архива: %v", err)
	}

	archive, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть архив: %v", err)
	}
	defer archive.Close()
	defer os.Remove(archivePath)

	err = conn.WriteMessage(websocket.BinaryMessage, func() []byte {
		data, _ := io.ReadAll(archive)
		return data
	}())
	if err != nil {
		return fmt.Errorf("не удалось отправить архив: %v", err)
	}

	return nil
}
