package filesystem

import (
	"errors"
	"fmt"
	"os"
)

const (
	NotFound  = -1
	errorCode = 0
)

func ReadFile(data []byte, fileName string) (int, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return 0, fmt.Errorf("opening the file %s has error %s", fileName, err.Error())
	}

	defer file.Close()

	n, err := file.Read(data)

	if n == len(data) {
		return errorCode, fmt.Errorf("reading the file %s in buffer is overflow", fileName)
	}

	if err != nil {
		return errorCode, fmt.Errorf("reading the file %s has error %s", fileName, err.Error())
	}

	return n, nil
}

func ReadAll(directory string, fileNames []string, files [][]byte) ([][]byte, error) {

	err := ForEach(directory, fileNames, func(file *os.File) error {
		stat, err := file.Stat()

		if err != nil {
			return errors.New("could not get stats of file")
		}

		buf := make([]byte, stat.Size())

		_, err = file.Read(buf)

		if err != nil {
			return err
		}

		files = append(files, buf)

		return nil
	})

	return files, err
}

func ForEach(directory string, fileNames []string, action func(file *os.File) error) error {
	errorStr := "files with error: "
	var hasErrorFiles bool

	for _, name := range fileNames {

		file, err := os.OpenFile(directory+name, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			hasErrorFiles = true
			errorStr += name + ", "
		}

		err = action(file)
		file.Close()

		if err != nil {
			hasErrorFiles = true
			errorStr += name + ", "
		}
	}

	if hasErrorFiles {
		return errors.New(errorStr)
	}

	return nil
}

func MakeFileNames(directory string) ([]string, error) {
	files, err := os.ReadDir(directory)

	if err != nil {
		return []string(nil), fmt.Errorf("could not read directory %s with error", directory)
	}

	fileNames := make([]string, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

func FindLastFile(directory string) (string, error) {
	names, err := MakeFileNames(directory)

	if err != nil {
		return "", errors.New("could not find last file name")
	}

	if len(names) == 0 {
		return "", nil
	}

	return names[len(names)-1], nil
}

func FindFile(names []string, target string) int {
	left := 0
	right := len(names) - 1

	for left < right {
		mid := (left + right) / 2

		if names[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}

	if names[left] != target {
		return NotFound
	}

	return left
}

func WriteFiles(dir string, fileNames []string, fileData [][]byte) error {
	if len(fileNames) != len(fileData) {
		return fmt.Errorf("could not write %d file names with %d file data", len(fileNames), len(fileData))
	}

	i := 0
	err := ForEach(dir, fileNames, func(file *os.File) error {

		data := fileData[i]
		i++

		n, err := file.Write(data)

		if err != nil {
			return fmt.Errorf("written only %d bytes", n)
		}

		return nil
	})

	return err
}
