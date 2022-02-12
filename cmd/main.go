package main

import (
	"log"
	"os"
	"pathconverter"
)

var (
	FILE_PATH string = "./path.yaml"
)

func main() {
	// Читаем файл
	root, err := pathconverter.ReadFile(FILE_PATH)
	if err != nil {
		log.Printf("Failed read file - %s\n", err.Error())
		os.Exit(1)
	}

	// Путь не должен быть пустой
	if !root.CheckPath() {
		log.Println("Path is empty. Please check file path.yaml")
		os.Exit(1)
	}

	// Содержимое директории
	if err := root.DisplayFiles(); err != nil {
		log.Printf("Failed display files - %s\n", err.Error())
		os.Exit(1)
	}

	// Переименовывание файлов
	if err := root.RenameFilesInDirectory(); err != nil {
		log.Printf("Failed rename files - %s\n", err.Error())
		os.Exit(1)
	}

	log.Println("Success")
}
