package models

// Card - игральная карта
type Card struct {
	// Стоимость карты
	Worth int
	// Сама карта
	Value string
}

// PlayerDeck - карта игроков
type PlayerDeck struct {
	Deck       []string
	PointsQuan int
	IsWin      bool
}
