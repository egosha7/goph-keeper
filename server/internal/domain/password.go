package domain

// PassData содержит информацию о пароле.
type PassData struct {
	Login    string `json:"login"`    // Логин пользователя
	PassName string `json:"passName"` // Название пароля
}

// PasswordData содержит информацию о новом пароле.
type PasswordData struct {
	Login    string `json:"login"`    // Логин пользователя
	PassName string `json:"passName"` // Название нового пароля
	Password string `json:"password"` // Новый пароль
}
