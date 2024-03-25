package main

import (
	"fmt"
	"runtime"
)

var (
	// Версия приложения (может быть перезаписана во время сборки)
	Version = "dev"
	// Дата сборки приложения (может быть перезаписана во время сборки)
	BuildDate = "unknown"
)

func main() {
	// Вывод информации о версии и дате сборки
	fmt.Println("My CLI App")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
