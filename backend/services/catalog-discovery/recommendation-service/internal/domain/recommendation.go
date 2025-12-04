package domain

type RecommendedItem struct {
	ProductID string
	Score     float64
	Reason    string
}

type Interaction struct {
	UserID          string
	ProductID       string
	InteractionType string
}
