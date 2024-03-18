package usecase

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
)

var (
	ErrOpenMigrationFile = errors.New("error open migration folder")
	ErrStopOfDuplicate   = errors.New("stop of duplicate")
)

func Create(name string) error {
	file, err := createEmptyFile(name)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	_, err = file.Write(getTemplate())
	if err != nil {
		return fmt.Errorf("file write template: %w", err)
	}

	fmt.Printf("Create migration file: %s \n", file.Name())
	return nil
}

func getTemplate() []byte {
	template := migrationFileTemplate
	return []byte(template)
}

// создаём папку.
func createMigrationFolder() bool {
	err := os.MkdirAll(config.FolderName, 0o755)
	return err == nil
}

// приводим стиль названия папки к snake_case.
func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// проверяем на дублирование названия миграции, в случае совпадения просим подтвердить создание.
func continueIfDuplicates(name string) bool {
	var isDuplucate bool
	var duplName string

	files, _ := os.ReadDir(config.FolderName)
	for _, file := range files {
		if strings.Contains(file.Name(), name) {
			isDuplucate = true
			duplName = file.Name()
		}
	}

	if isDuplucate {
		var done string
		fmt.Printf("Мы нашли дубль миграции '%s' продолжить? Y/n \n", duplName)
		fmt.Scanf("%s\n", &done)
		return done == "Y"
	}
	return true
}

// создаем пустой файл.
func createEmptyFile(name string) (*os.File, error) {
	ok := createMigrationFolder()
	if !ok {
		return nil, ErrOpenMigrationFile
	}
	lowerName := toSnakeCase(name)
	continueIfDuplicates := continueIfDuplicates(lowerName)
	if !continueIfDuplicates {
		fmt.Println("Stop of duplicate!")
		return nil, ErrStopOfDuplicate
	}
	timestamp := time.Now().Format("2006_01_02_150405")
	filename := fmt.Sprintf("%s/%s_%s.sql", config.FolderName, timestamp, lowerName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
