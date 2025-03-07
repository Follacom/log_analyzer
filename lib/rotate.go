package lib

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func setFileMetadata(filePath string) error {
	dir := filepath.Dir(filePath)
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), ext)

	metadata := map[string]string{
		"month": time.Now().Format("01"),
	}

	data, _ := json.Marshal(metadata)

	metaFile := filepath.Join(dir, fmt.Sprintf("%s_meta.json", base))
	if _, err := os.Stat(metaFile); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0644); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return os.WriteFile(metaFile, data, 0644)
}

func getFileMetadata(filePath string) (map[string]string, error) {
	dir := filepath.Dir(filePath)
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), ext)

	data, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("%s_meta.json", base)))
	if err != nil {
		return nil, err
	}

	var metadata = make(map[string]string, 0)
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func archiveLogs(dir, subdir, base, month string) error {
	archiveFile := filepath.Join(dir, subdir, fmt.Sprintf("%s_%s.zip", base, month))
	if _, err := os.Stat(archiveFile); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(archiveFile), 0644); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	archive, err := os.Create(archiveFile)
	if err != nil {
		LogError(err)
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	files, err := filepath.Glob(filepath.Join(dir, base+"_*.*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if file == archiveFile {
			continue
		}
		fileToZip, err := os.OpenFile(file, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		zipEntry, err := zipWriter.Create(filepath.Base(file))
		if err != nil {
			return err
		}

		if _, err := io.Copy(zipEntry, fileToZip); err != nil {
			return err
		}
		if err = fileToZip.Close(); err != nil {
			return err
		}

		err = os.Remove(file) // Remove old log after archiving
		if err != nil {
			return err
		}
	}
	return nil
}

func RotateFile(filePath string) string {
	dir := filepath.Dir(filePath)
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), ext)

	shouldRotate := viper.GetViper().GetBool("rotate")
	if shouldRotate {
		backupTimestamp := time.Now().Format("02")

		backupFile := filepath.Join(dir, fmt.Sprintf("%s_%s%s", base, backupTimestamp, ext))
		return backupFile
	}

	return filePath
}

func RotateDatabase() (*gorm.DB, error) {
	filePath := viper.GetViper().GetString("database.url")
	backupFile := RotateFile(filePath)

	shouldRotate := viper.GetViper().GetBool("rotate")
	if shouldRotate {
		metadata, err := getFileMetadata(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				setFileMetadata(filePath)
			} else {
				LogError(err)
			}
		}

		if metadata["month"] != "" && metadata["month"] != time.Now().Format("01") {
			// archive all files
			dir := filepath.Dir(filePath)
			ext := filepath.Ext(filePath)
			base := strings.TrimSuffix(filepath.Base(filePath), ext)

			month, err := time.Parse("01", metadata["month"])
			if err != nil {
				LogError(err)
			}
			err = archiveLogs(dir, "archives", base, month.Month().String())
			if err != nil {
				LogError(err)
			}

			setFileMetadata(filePath)
		}
	}

	return gorm.Open(sqlite.Open(backupFile), &gorm.Config{
		SkipDefaultTransaction: true,
	})
}
