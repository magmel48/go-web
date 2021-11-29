package shortener

import (
	"bufio"
	"os"
	"strings"
)

type FileBackup struct {
	file *os.File
}

func NewFileBackup(filePath string) FileBackup {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	return FileBackup{file: file}
}

func (fb FileBackup) Append(record string) {
	_, err := fb.file.Write([]byte(record))
	if err != nil {
		panic(err)
	}
}

func (fb FileBackup) ReadAll() map[string]string {
	result := make(map[string]string)

	scanner := bufio.NewScanner(fb.file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ln := strings.Split(scanner.Text(), linksDelimiter)

		if len(ln) > 1 {
			result[ln[0]] = ln[1]
		}
	}

	return result
}
