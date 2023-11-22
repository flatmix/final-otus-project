package usecase

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const folderName = "migrations"

func Create(name string) error {
	file, err := createEmptyFile(name)
	if err != nil {
		return err
	}
	_, err = file.Write(getTemplate())
	if err != nil {
		return err
	}
	return nil
}

func getTemplate() []byte {
	template :=
		`--migrate:UP
--write your sql for migration...


--migrate:DOWN
--write your sql for rollback migration...

`
	return []byte(template)
}

func createMigrationFolder() bool {
	err := os.MkdirAll(folderName, 0755)
	return err == nil
}

func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func continueIfDuplicates(name string) bool {
	var isDuplucate bool
	var duplName string

	files, _ := os.ReadDir(folderName)
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

func createEmptyFile(name string) (*os.File, error) {
	ok := createMigrationFolder()
	if !ok {
		return nil, errors.New("Error open migration folder!")
	}
	lowerName := toSnakeCase(name)
	continueIfDuplicates := continueIfDuplicates(lowerName)
	if !continueIfDuplicates {
		return nil, errors.New("")
	}
	timestamp := time.Now().Format("2006_01_02_150405")
	filename := fmt.Sprintf("%s/%s_%s.sql", folderName, timestamp, lowerName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil

}
