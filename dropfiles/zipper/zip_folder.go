package zipper

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CompressFolder(folderPath string) error {
	fs, err := os.Open(folderPath)
	if err != nil {
		return fmt.Errorf("open error: %w", err)
	}
	defer fs.Close()

	fi, err := fs.Stat()
	if err != nil {
		return fmt.Errorf("stat error: %w", err)
	}

	if !fi.IsDir() {
		return errors.New("path is not folder")
	}

	zipPath := fmt.Sprintf("%s.zip", folderPath)

	var buffer bytes.Buffer
	err = Compress(&buffer, folderPath)
	if err != nil {
		return fmt.Errorf("zip error: %w", err)
	}

	err = ioutil.WriteFile(zipPath, buffer.Bytes(), os.ModePerm.Perm())
	if err != nil {
		return fmt.Errorf("write zip file error: %w", err)
	}

	return nil

}

func Compress(writer io.Writer, target string) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	isDir, err := isDirectory(target)
	if err != nil {
		return err
	}

	if isDir {
		if err := addZipFiles(zipWriter, target, ""); err != nil {
			return err
		}
	} else {
		fileName := filepath.Base(target)
		if err := addZipFile(zipWriter, target, fileName); err != nil {
			return err
		}
	}

	return nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("stat error: %w", err)
	}
	return fileInfo.IsDir(), nil
}

func addZipFiles(writer *zip.Writer, basePath, pathInZip string) error {
	fileInfoArray, err := ioutil.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("stat error: %w", err)
	}

	basePath = complementPath(basePath)
	pathInZip = complementPath(pathInZip)

	for _, fileInfo := range fileInfoArray {
		newBasePath := basePath + fileInfo.Name()
		newPathInZip := pathInZip + fileInfo.Name()

		if fileInfo.IsDir() {
			if err = addDirectory(writer, newBasePath); err != nil {
				return err
			}

			newBasePath = newBasePath + string(os.PathSeparator)
			newPathInZip = newPathInZip + string(os.PathSeparator)
			if err = addZipFiles(writer, newBasePath, newPathInZip); err != nil {
				return err
			}
		} else {
			if err = addZipFile(writer, newBasePath, newPathInZip); err != nil {
				return err
			}
		}
	}

	return nil
}

func addZipFile(writer *zip.Writer, targetFilePath, pathInZip string) error {
	data, err := ioutil.ReadFile(targetFilePath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	fileInfo, err := os.Lstat(targetFilePath)
	if err != nil {
		return fmt.Errorf("lstat error: %w", err)
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return fmt.Errorf("fileinfoheader: %w", err)
	}

	header.Name = pathInZip
	header.Method = zip.Store
	w, err := writer.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create zip header: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write zip: %w", err)
	}
	return nil
}

func addDirectory(writer *zip.Writer, basePath string) error {
	fileInfo, err := os.Lstat(basePath)
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	if _, err = writer.CreateHeader(header); err != nil {
		return err
	}
	return nil
}

func complementPath(path string) string {
	l := len(path)
	if l == 0 {
		return path
	}

	lastChar := path[l-1 : l]
	if lastChar == "/" || lastChar == "\\" {
		return path
	} else {
		return path + string(os.PathSeparator)
	}
}
