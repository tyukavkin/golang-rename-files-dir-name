package pathconverter

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

/* Модификация

 */

var (
	TYPES = [2]string{"jpg", "png"}
	// %s %d %s Name Number Shared - %s.%d.%s (Name.1.PNG), %s_%d.%s (Name_1.PNG)
	FILE_NAME_FORMAT = "%s.%d.%s"
)

type RootDirectory struct {
	Path string `yaml:"rootDir"`
}

func ReadFile(filepath string) (*RootDirectory, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var rootDirectory = &RootDirectory{}
	err = yaml.Unmarshal(file, rootDirectory)
	if err != nil {
		return nil, err
	}

	log.Println(rootDirectory.Path)
	return rootDirectory, nil
}

func (rd *RootDirectory) CheckPath() bool {
	return rd.Path != ""
}

func (rd *RootDirectory) DisplayFiles() error {
	checkFile, err := os.Stat(rd.Path)
	if err != nil {
		return err
	}

	if !checkFile.IsDir() {
		return fmt.Errorf("%s - is not directory", rd.Path)
	}

	files, err := os.ReadDir(rd.Path)
	if err != nil {
		return fmt.Errorf("failed path - %s. Cancel with error - %s", rd.Path, err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			log.Printf("Dir  - %s", file.Name())
		} else {
			log.Printf("File - %s", file.Name())
		}
	}
	return nil
}

func (rd *RootDirectory) RenameFilesInDirectory() error {
	checkFile, err := os.Stat(rd.Path)
	if err != nil {
		return err
	}

	if !checkFile.IsDir() {
		return fmt.Errorf("%s - is not directory", rd.Path)
	}

	files, err := os.ReadDir(rd.Path)
	if err != nil {
		return fmt.Errorf("failed path - %s. Cancel with error - %s", rd.Path, err.Error())
	}

	for _, file := range files {
		// Пропускаем если файл
		if !file.IsDir() {
			continue
		}

		log.Printf("In directory - %s", file.Name())
		dirPhoto := filepath.Join(rd.Path, file.Name())

		// Заходим в каждую директорию
		photos, err := os.ReadDir(dirPhoto)
		if err != nil {
			continue
		}

		var idx = 0

		// Перебераем файлы и переименовываем их
		for _, photo := range photos {
			if photo.IsDir() {
				continue
			}
			// Проверяем наличие расширения
			unitName := strings.Split(photo.Name(), ".")
			if len(unitName) <= 1 {
				continue
			}

			shared := unitName[len(unitName)-1]

			if isImage(shared) {
				// Формируем имена файлов
				oldName := filepath.Join(dirPhoto, photo.Name())
				newName := filepath.Join(dirPhoto, fmt.Sprintf(FILE_NAME_FORMAT, file.Name(), idx, shared))
				if oldName == newName {
					idx++
					continue
				}

				log.Printf("Rename file %s to %s", oldName, newName)

				// Берем старый файл
				inputFile, err := os.Open(oldName)
				if err != nil {
					log.Printf("Couldn't open source file: %s", err)
					continue
				}

				// Создаем файл в который будем
				outputFile, err := os.Create(newName)
				if err != nil {
					inputFile.Close()
					log.Printf("Couldn't open dest file: %s", err)
					continue
				}
				idx++

				_, err = io.Copy(outputFile, inputFile)

				// Закрываем исходный файл
				inputFile.Close()
				if err != nil {
					log.Printf("Writing to output file failed: %s", err)
					continue
				}

				// Закрываем новый файл
				outputFile.Close()

				// Удаляем старый файл
				err = os.Remove(oldName)
				if err != nil {
					log.Printf("Failed removing original file: %s", err)
					continue
				}
			}
		}
	}
	return nil
}

func isImage(file string) bool {
	for _, typeFile := range TYPES {
		if typeFile == strings.ToLower(file) {
			return true
		}
	}
	return false
}
