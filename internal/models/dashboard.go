package models

type SummaryResponse struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type MonthlyTrendResponse struct {
	Month int     `json:"month"`
	Total float64 `json:"total"`
}
