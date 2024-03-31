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

type CreateStruct struct {
	cfg  *config.Config
	name string
}

func NewCreateStruct(name string, cfg *config.Config) *CreateStruct {
	return &CreateStruct{cfg: cfg, name: name}
}

func Create(name string, configuration *config.Config) error {
	createStruct := NewCreateStruct(name, configuration)
	file, err := createStruct.createEmptyFile()
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	_, err = file.Write(createStruct.getTemplate())
	if err != nil {
		return fmt.Errorf("file write template: %w", err)
	}

	fmt.Printf("Create migration file: %s \n", file.Name())
	return nil
}

func (cs *CreateStruct) getTemplate() []byte {
	template := migrationFileTemplate
	return []byte(template)
}

// создаём папку.
func (cs *CreateStruct) createMigrationFolder() bool {
	err := os.MkdirAll(cs.cfg.FolderName, 0o755)
	return err == nil
}

// приводим стиль названия папки к snake_case.
func (cs *CreateStruct) toSnakeCase() string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(cs.name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// проверяем на дублирование названия миграции, в случае совпадения просим подтвердить создание.
func (cs *CreateStruct) continueIfDuplicates() bool {
	var isDuplucate bool
	var duplName string

	name := cs.toSnakeCase()

	files, _ := os.ReadDir(cs.cfg.FolderName)
	for _, file := range files {
		if strings.Contains(file.Name(), name) {
			isDuplucate = true
			duplName = file.Name()
		}
	}

	if isDuplucate {
		var done string
		fmt.Printf("Мы нашли дубль миграции '%s' продолжить? Y/n \n", duplName)
		_, err := fmt.Scanf("%s\n", &done)
		if err != nil {
			return false
		}
		return done == "Y"
	}
	return true
}

// создаем пустой файл.
func (cs *CreateStruct) createEmptyFile() (*os.File, error) {
	ok := cs.createMigrationFolder()
	if !ok {
		return nil, ErrOpenMigrationFile
	}

	continueIfDuplicates := cs.continueIfDuplicates()
	if !continueIfDuplicates {
		fmt.Println("Stop of duplicate!")
		return nil, ErrStopOfDuplicate
	}
	timestamp := time.Now().Format("2006_01_02_150405")
	filename := fmt.Sprintf("%s/%s_%s.sql", cs.cfg.FolderName, timestamp, cs.toSnakeCase())
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
