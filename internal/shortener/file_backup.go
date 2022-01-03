package shortener

import (
	"bufio"
	"os"
	"strings"
)

// OpenFile created for tests purposes to be sure that program opens proper file with proper file name.
type OpenFile func(name string, flag int, perm os.FileMode) (*os.File, error)

// FileBackup backup that stores data on disk.
type FileBackup struct {
	file *os.File
}

// NewFileBackup creates new file backup.
func NewFileBackup(filePath string, openFile OpenFile) (*FileBackup, error) {
	file, err := openFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	backup := FileBackup{file: file}
	return &backup, nil
}

// Append appends a record to existing stored records on disk.
func (fb *FileBackup) Append(record string) error {
	_, err := fb.file.Write([]byte(record))
	return err
}

// ReadAll reads all stored data from disk.
func (fb *FileBackup) ReadAll() map[string]string {
	result := make(map[string]string)

	scanner := bufio.NewScanner(fb.file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ln := strings.Split(scanner.Text(), delimiter)

		if len(ln) > 1 {
			result[ln[0]] = ln[1]
		}
	}

	return result
}
