package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

var (
	// Версия приложения (может быть перезаписана во время сборки)
	Version = "dev"
	// Дата сборки приложения (может быть перезаписана во время сборки)
	BuildDate = "unknown"
)

type UserInfo struct {
	Login string `json:"login"`
}

// SiteList представляет список сайтов
type SiteList struct {
	Sites []string `json:"sites"`
}

func main() {
	// Вывод информации о версии и дате сборки
	fmt.Println(string(colorGreen), name, string(colorReset))
	showStartMenu()

	// Создание канала для сигналов
	sigChan := make(chan os.Signal, 1)
	// Регистрация обработчика сигнала завершения
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Ожидание сигнала
	for {
		select {
		case sig := <-sigChan:
			if sig == syscall.SIGQUIT {
				fmt.Println("Exiting...")
				os.Exit(0)
			}
		}
	}
}

// Функция для вывода меню выбора действия
func showStartMenu() {
	fmt.Println("1. Регистрация")
	fmt.Println("2. Авторизация")
	fmt.Println("3. Версия и дата сборки")
	fmt.Print("Выберите действие: ")

	// Получаем выбор пользователя
	reader := bufio.NewReader(os.Stdin)
	choice, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	// Обрабатываем выбор пользователя
	switch choice {
	case '1':
		fmt.Println("Вы выбрали регистрацию.")
		registerUser()
	case '2':
		fmt.Println("Вы выбрали авторизацию.")
		loginUser()
	case '3':
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println()
		showStartMenu()
	default:
		fmt.Println("Неверный выбор. Пожалуйста, выберите 1 или 2.")
	}
}

func showMenu(login string) {
	for {
		// Выводим меню
		fmt.Println("\nМеню:")
		fmt.Println("1. Просмотреть данные")
		fmt.Println("2. Внести новый пароль")
		fmt.Println("3. Внести новую карту")
		fmt.Println("4. Поменять пин-код")
		fmt.Println("0. Выйти")

		// Получаем выбор пользователя
		fmt.Print("Выберите пункт меню: ")
		reader := bufio.NewReader(os.Stdin)
		choice, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("Ошибка при чтении ввода:", err)
			continue
		}

		// Обрабатываем выбор пользователя
		switch choice {
		case '1':
			fmt.Println("\nПросмотр данных:")
			fmt.Println("1. Пароль")
			fmt.Println("2. Карта")
			subChoice := getUserInputInfo("Выберите подпункт: ")
			switch subChoice {
			case "1":
				viewPasswordsName(login)
			case "2":
				fmt.Println("Ваши данные по картам")
			default:
				fmt.Println("Некорректный подпункт меню")
			}
		case '2':
			fmt.Println("Внести новый пароль")
		case '3':
			fmt.Println("Внести новую карту")
		case '4':
			fmt.Println("Поменять пин-код")
		case '0':
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Некорректный пункт меню")
		}
	}
}

func viewPasswordsName(login string) {
	// Создаем JSON-объект с информацией о пользователе
	userInfo := UserInfo{
		Login: login,
	}

	// Преобразуем информацию о пользователе в JSON
	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	// Отправляем POST-запрос на сервер
	resp, err := http.Post("http://localhost:8080/pass/namelist", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка: сервер вернул статус", resp.Status)
		return
	}

	// Создаем срез строк для хранения списка паролей
	var passwords []string
	err = json.NewDecoder(resp.Body).Decode(&passwords)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	for {
		// Выводим список паролей
		fmt.Println("\nСписок паролей:")
		for i, password := range passwords {
			fmt.Printf("%d. %s\n", i+1, password)
		}
		fmt.Println("0. Вернуться назад")

		// Получаем выбор пользователя
		fmt.Print("Выберите номер пароля: ")
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Ошибка при чтении ввода:", err)
			continue
		}

		// Обрабатываем выбор пользователя
		if choice == 0 {
			fmt.Println("Возвращаемся назад...")
			return
		} else if choice < 1 || choice > len(passwords) {
			fmt.Println("Некорректный номер пароля")
			continue
		}

		// Выводим выбранный пароль
		selectedPassword := passwords[choice-1]
		fmt.Printf("Выбранный пароль: %s\n", selectedPassword)

		// Здесь можно добавить дополнительную логику для работы с выбранным паролем
	}
}

// Функция для регистрации пользователя
func registerUser() {
	fmt.Println("Введите ваш email:")
	email := getUserInput()

	// Проверяем валидность email с помощью регулярного выражения
	if !isValidEmail(email) {
		fmt.Println("Некорректный email. Пожалуйста, попробуйте снова.")
		registerUser()
		return
	}

	fmt.Println("Введите ваш пароль:")
	password1 := getHiddenUserInput()

	fmt.Println("Повторите ваш пароль:")
	password2 := getHiddenUserInput()

	// Проверяем совпадение паролей
	if password1 != password2 {
		fmt.Println("Пароли не совпадают. Пожалуйста, попробуйте снова.")
		registerUser()
		return
	}

	// Здесь вы можете отправить запрос на сервер для регистрации
	fmt.Printf("Вы ввели email: %s и пароль: %s\n", email, password1)

	// Создаем данные для отправки в формате JSON
	data := map[string]string{
		"login":    email,
		"password": password2,
	}

	// Кодируем данные в формат JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка при кодировании данных в JSON:", err)
		return
	}

	// Отправляем POST запрос на сервер
	resp, err := http.Post("http://localhost:8080/auth/registration", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем код статуса ответа
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно зарегистрирован")
		showMenu(email)
	} else if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusConflict {
			fmt.Println("Ошибка при регистрации пользователя: данный логин уже зарегистрирован")
			registerUser()
		}
		fmt.Println("Ошибка при регистрации пользователя:", resp.Status)
		registerUser()
	}

}

// Функция для аутентификации пользователя
func loginUser() {
	fmt.Println("Введите ваш email:")
	email := getUserInput()

	fmt.Println("Введите ваш пароль:")
	password := getHiddenUserInput()

	// Здесь вы можете отправить запрос на сервер для аутентификации
	fmt.Printf("Вы ввели email: %s и пароль: %s\n", email, password)

	// Создаем данные для отправки в формате JSON
	data := map[string]string{
		"login":    email,
		"password": password,
	}

	// Кодируем данные в формат JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка при кодировании данных в JSON:", err)
		return
	}

	// Отправляем POST запрос на сервер
	resp, err := http.Post("http://localhost:8080/auth", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем код статуса ответа
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно авторизован")
		showMenu(email)
	} else if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Ошибка при авторизации пользователя: неправильно введен e-mail или пароль")
			fmt.Println("Попробуйте еще раз или вернитесь в меню")
			fmt.Println() // Переход на следующую строку после ввода пароля
			loginUser()
		}
		fmt.Println("Ошибка при регистрации пользователя:", resp.Status)
		return
	}
}

// Функция для получения ввода пользователя
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		os.Exit(1)
	}
	// Удаляем символ новой строки из ввода
	input = strings.TrimSpace(input)
	return input
}

// Функция для получения ввода пользователя
func getUserInputInfo(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return ""
	}
	// Удаляем символ новой строки из ввода
	input = strings.TrimSpace(input)
	return input
}

// Функция для получения скрытого ввода пароля
func getHiddenUserInput() string {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		os.Exit(1)
	}
	fmt.Println() // Переход на следующую строку после ввода пароля
	return string(bytePassword)
}

// Функция для проверки валидности email
func isValidEmail(email string) bool {
	// Простая проверка с помощью регулярного выражения
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
