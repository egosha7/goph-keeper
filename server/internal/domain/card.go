package domain

// CardData содержит информацию о карте.
type CardData struct {
	Login    string `json:"login"`    // Логин пользователя
	CardName string `json:"cardName"` // Название карты
}

// NewCardData содержит информацию о новой карте.
type NewCardData struct {
	Login          string `json:"login"`          // Логин пользователя
	CardName       string `json:"cardName"`       // Название новой карты
	NumberCard     string `json:"numberCard"`     // Номер новой карты
	ExpiryDateCard string `json:"expiryDateCard"` // Срок новой карты
	CvvCard        string `json:"CvvCard"`        // Секретный код новой карты
}

// CardInfo содержит информацию о банковской карте.
type CardInfo struct {
	Number     string `json:"number"`     // Номер карты
	ExpiryDate string `json:"expiryDate"` // Срок действия карты
	CVV        string `json:"cvv"`        // CVV карты
}
