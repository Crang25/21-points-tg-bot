package logger

import (
	"log"
	"os"
)

var (
	// Путь, по которому лежит файл с логами
	path = "Path_to_file"
	// NewFile файл, в который будут записываться логи
	NewFile, _ = os.OpenFile(path, os.O_RDWR, 0755)
	// LogFile ...
	LogFile = log.New(NewFile, "", 0)
)

// CheckErr проверяет на ошибку, если она есть записываем в заранее созданные файл с логами
func CheckErr(err error) {
	if err != nil {
		LogFile.Fatalf("failed to: %v\n", err)
	}
}
