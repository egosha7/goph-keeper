package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
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

// UserInfo представляет информацию о пользователе.
type UserInfo struct {
	Login string `json:"login"` // Логин пользователя
}

// PinData содержит информацию о пользователе и пин-коде.
type PinData struct {
	Login string `json:"login"` // Логин пользователя
	Pin   string `json:"pin"`   // Пин-код пользователя
}

// PassData содержит информацию о пользователе и названии пароля.
type PassData struct {
	Login    string `json:"login"`    // Логин пользователя
	PassName string `json:"passName"` // Название пароля
}

// CardData содержит информацию о пользователе и названии карты.
type CardData struct {
	Login    string `json:"login"`    // Логин пользователя
	CardName string `json:"cardName"` // Название карты
}

// CardInfo содержит информацию о банковской карте.
type CardInfo struct {
	Number     string `json:"number"`     // Номер карты
	ExpiryDate string `json:"expiryDate"` // Срок действия карты
	CVV        string `json:"cvv"`        // CVV карты
}

// PasswordData содержит информацию о новом пароле.
type PasswordData struct {
	Login    string `json:"login"`    // Логин пользователя
	PassName string `json:"passName"` // Название нового пароля
	Password string `json:"password"` // Новый пароль
}

// NewCardData содержит информацию о новой карте.
type NewCardData struct {
	Login          string `json:"login"`          // Логин пользователя
	CardName       string `json:"cardName"`       // Название новой карты
	NumberCard     string `json:"numberCard"`     // Номер новой карты
	ExpiryDateCard string `json:"expiryDateCard"` // Срок новой карты
	CvvCard        string `json:"CvvCard"`        // Секретный код новой карты
}

// Функция с которой начинается работа программы
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
		return
	}
}

// showMenu отображает главное меню приложения и обрабатывает выбор пользователя.
// Параметр login представляет логин текущего пользователя.
func showMenu(login string) {
	for {
		// Выводим меню
		fmt.Println("\nМеню:")
		fmt.Println("1. Просмотреть данные")
		fmt.Println("2. Внести новый пароль")
		fmt.Println("3. Внести новую карту")
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
				viewCardsName(login)
			default:
				fmt.Println("Некорректный подпункт меню")
			}
		case '2':
			fmt.Println("Внести новый пароль")
			addNewPassword(login)
		case '3':
			fmt.Println("Внести новую карту")
			addNewCard(login)
		case '0':
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Некорректный пункт меню")
		}
	}
}

// addNewCard запрашивает у пользователя информацию о новой банковской карте и отправляет ее на сервер для добавления.
// Параметр login представляет логин текущего пользователя.
func addNewCard(login string) {
	// Получаем информацию о новой карте от пользователя
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите название карты: ")
	cardName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	fmt.Print("Введите номер карты: ")
	numberCard, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	fmt.Print("Введите дату срока карты (Например: 03/24): ")
	expiryDateCard, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	fmt.Print("Введите секретный номер карты (CVV/CVC): ")
	CvvCard, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	// Отправляем данные новой карты на сервер
	success, err := AddNewCard(login, cardName, numberCard, expiryDateCard, CvvCard)
	if err != nil {
		fmt.Println("Ошибка при добавлении карты:", err)
		return
	}

	if success {
		fmt.Println("Новая карта успешно добавлена!")
	} else {
		fmt.Println("Не удалось добавить новую карту.")
	}
}

func AddNewCard(login, cardName, numberCard, expiryDateCard, CvvCard string) (bool, error) {
	// Создаем JSON-объект с информацией о пользователе и новом пароле
	passwordData := NewCardData{
		Login:          login,
		CardName:       cardName,
		NumberCard:     numberCard,
		ExpiryDateCard: expiryDateCard,
		CvvCard:        CvvCard,
	}

	// Преобразуем информацию в JSON
	jsonData, err := json.Marshal(passwordData)
	if err != nil {
		return false, fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем POST-запрос на сервер для добавления нового пароля
	resp, err := http.Post("http://localhost:8080/card/add", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ошибка: сервер вернул статус %s", resp.Status)
	}

	return true, nil
}

func addNewPassword(login string) {
	// Получаем информацию о новом пароле от пользователя
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите название нового пароля: ")
	passName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	fmt.Print("Введите новый пароль: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}

	// Отправляем новый пароль на сервер
	success, err := AddPassword(login, passName, password)
	if err != nil {
		fmt.Println("Ошибка при добавлении нового пароля:", err)
		return
	}

	if success {
		fmt.Println("Новый пароль успешно добавлен!")
	} else {
		fmt.Println("Не удалось добавить новый пароль.")
	}
}

// AddPassword отправляет запрос на сервер для добавления нового пароля.
// Параметр login представляет логин пользователя.
// Параметр name представляет название нового пароля.
// Параметр password представляет сам пароль.
// Функция возвращает true, если операция добавления прошла успешно, иначе возвращает ошибку.
func AddPassword(login, name, password string) (bool, error) {
	// Создаем JSON-объект с информацией о пользователе и новом пароле
	passwordData := PasswordData{
		Login:    login,
		PassName: name,
		Password: password,
	}

	// Преобразуем информацию в JSON
	jsonData, err := json.Marshal(passwordData)
	if err != nil {
		return false, fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем POST-запрос на сервер для добавления нового пароля
	resp, err := http.Post("http://localhost:8080/password/add", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ошибка: сервер вернул статус %s", resp.Status)
	}

	return true, nil
}

// viewCardsName отправляет запрос на сервер для получения списка названий карт пользователя и выводит их на экран.
// Параметр login представляет логин пользователя.
func viewCardsName(login string) {
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
	resp, err := http.Post("http://localhost:8080/card/namelist", "application/json", bytes.NewBuffer(jsonData))
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

	// Создаем срез строк для хранения списка карт
	var namecards []string
	err = json.NewDecoder(resp.Body).Decode(&namecards)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	for {
		// Выводим список карт
		fmt.Println("\nСписок кард:")
		for i, cards := range namecards {
			fmt.Printf("%d. %s\n", i+1, cards)
		}
		fmt.Println("0. Вернуться назад")

		// Получаем выбор пользователя
		fmt.Print("Выберите номер карты: ")
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
		} else if choice < 1 || choice > len(namecards) {
			fmt.Println("Некорректный номер карты")
			continue
		}

		// Получаем пин-код от пользователя
		fmt.Print("Введите пин-код: ")
		var pinCode string
		_, err = fmt.Scanln(&pinCode)
		if err != nil {
			fmt.Println("Ошибка при чтении ввода пин-кода:", err)
			continue
		}

		// Проверяем пин-код на сервере
		valid, err := checkPinCode(login, pinCode)
		if err != nil {
			fmt.Println("Ошибка при проверке пин-кода:", err)
			continue
		}

		// Если пин-код валиден, выводим данные о карте
		if valid {
			// Выводим выбранную карту
			selectedNameCard := namecards[choice-1]
			cardNumber, cardExpiry, cardCVV, _ := GetCard(login, selectedNameCard)
			fmt.Printf("Данные от карты '%s':\n", selectedNameCard)
			fmt.Printf("Номер карты: %s\n", cardNumber)
			fmt.Printf("Срок действия: %s\n", cardExpiry)
			fmt.Printf("CVV: %s\n", cardCVV)
			return
		} else {
			fmt.Println("Неверный пин-код")
		}
	}
}

// GetCard отправляет запрос на сервер для получения данных о карте пользователя.
// Параметр login представляет логин пользователя.
// Параметр selectedNameCard представляет название выбранной карты.
// Возвращает номер карты, срок действия и CVV карты, а также ошибку, если таковая возникла.
func GetCard(login, selectedNameCard string) (string, string, string, error) {
	// Создаем JSON-объект с информацией о пользователе и названии карты
	cardData := CardData{
		Login:    login,
		CardName: selectedNameCard,
	}

	// Преобразуем информацию о пользователе и названии карты в JSON
	jsonData, err := json.Marshal(cardData)
	if err != nil {
		return "", "", "", fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем POST-запрос на сервер для получения данных о карте
	resp, err := http.Post("http://localhost:8080/card/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("ошибка: сервер вернул статус %s", resp.Status)
	}

	// Декодируем ответ в структуру CardInfo
	var cardInfo CardInfo
	if err := json.NewDecoder(resp.Body).Decode(&cardInfo); err != nil {
		return "", "", "", fmt.Errorf("ошибка при декодировании JSON: %v", err)
	}

	// Возвращаем информацию о карте
	return cardInfo.Number, cardInfo.ExpiryDate, cardInfo.CVV, nil
}

// viewPasswordsName отправляет запрос на сервер для получения списка названий паролей пользователя и их отображения.
// Параметр login представляет логин пользователя.
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
	var namepasswords []string
	err = json.NewDecoder(resp.Body).Decode(&namepasswords)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	for {
		// Выводим список паролей
		fmt.Println("\nСписок паролей:")
		for i, password := range namepasswords {
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
		} else if choice < 1 || choice > len(namepasswords) {
			fmt.Println("Некорректный номер пароля")
			continue
		}

		// Получаем пин-код от пользователя
		fmt.Print("Введите пин-код: ")
		var pinCode string
		_, err = fmt.Scanln(&pinCode)
		if err != nil {
			fmt.Println("Ошибка при чтении ввода пин-кода:", err)
			continue
		}

		// Проверяем пин-код на сервере
		valid, err := checkPinCode(login, pinCode)
		if err != nil {
			fmt.Println("Ошибка при проверке пин-кода:", err)
			continue
		}

		// Если пин-код валиден, выводим пароль
		if valid {
			// Выводим выбранный пароль
			selectedNamePassword := namepasswords[choice-1]
			fmt.Printf("Выбор: %s\n", selectedNamePassword)
			pass, _ := GetPassword(login, selectedNamePassword)
			fmt.Printf("Выбранный пароль: %s\n", pass)
			return
		} else {
			fmt.Println("Неверный пин-код")
		}
	}
}

// GetPassword отправляет запрос на сервер для получения пароля пользователя по выбранному названию.
// Параметр login представляет логин пользователя.
// Параметр selectedNamePassword представляет название выбранного пароля.
// Возвращает пароль и ошибку, если таковая возникла.
func GetPassword(login, selectedNamePassword string) (string, error) {
	// Создаем JSON-объект с информацией о пользователе и пин-коде
	passData := PassData{
		Login:    login,
		PassName: selectedNamePassword,
	}

	// Преобразуем информацию о пользователе и пин-коде в JSON
	jsonData, err := json.Marshal(passData)
	if err != nil {
		return "", fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем POST-запрос на сервер для проверки пин-кода
	resp, err := http.Post("http://localhost:8080/password/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка: сервер вернул статус %s", resp.Status)
	}

	// Читаем ответ и возвращаем полученный пароль
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении тела ответа: %v", err)
	}

	return string(body), nil
}

// checkPinCode отправляет запрос на сервер для проверки пин-кода пользователя.
// Параметр login представляет логин пользователя.
// Параметр pinCode представляет введенный пин-код.
// Возвращает true, если пин-код верен, и ошибку, если таковая возникла.
func checkPinCode(login, pinCode string) (bool, error) {
	// Создаем JSON-объект с информацией о пользователе и пин-коде
	pinData := PinData{
		Login: login,
		Pin:   pinCode,
	}

	// Преобразуем информацию о пользователе и пин-коде в JSON
	jsonData, err := json.Marshal(pinData)
	if err != nil {
		return false, fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем POST-запрос на сервер для проверки пин-кода
	resp, err := http.Post("http://localhost:8080/pincheck", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return false, fmt.Errorf("неверный пин-код")
		}
		return false, fmt.Errorf("Oшибка: сервер вернул статус %s", resp.Status)
	}

	return true, nil
}

// registerUser регистрирует нового пользователя.
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

	fmt.Println("Введите ваш новый пин-код:")
	pin := getUserInput()

	// Создаем данные для отправки в формате JSON
	data := map[string]string{
		"login":    email,
		"password": password2,
		"pin":      pin,
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
