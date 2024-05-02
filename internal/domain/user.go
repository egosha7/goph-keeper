package domain

// User структура для представления пользователя
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Pin      string `json:"pin"`
}

// UserInfo представляет информацию о пользователе.
type UserInfo struct {
	Login string `json:"Login"` // Логин пользователя
}

// CheckPinData представляет данные для проверки пин-кода.
type CheckPinData struct {
	Login string `json:"login"` // Логин пользователя
	Pin   string `json:"pin"`   // Пин-код
}

// CheckPinResponse представляет ответ о результате проверки пин-кода.
type CheckPinResponse struct {
	Valid bool `json:"valid"` // Результат проверки пин-кода
}
