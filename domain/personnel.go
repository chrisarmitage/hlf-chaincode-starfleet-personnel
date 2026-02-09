package domain

type Personnel struct {
	PersonnelID string `json:"personnelID"`
	Name        string `json:"name"`
	Rank        string `json:"rank"`
	Campus      string `json:"campus"`
	Status      string `json:"status"`
}

type Training struct {
	RecordID     string `json:"recordID"`
	PersonnelID  string `json:"personnelID"`
	Campus       string `json:"campus"`
	TrainingCode string `json:"trainingCode"`
	CompletedAt  string `json:"completedAt"`
	IssuedBy     string `json:"issuedBy"`
	Status       string `json:"status"`
}

const (
	PersonnelRankCadet = "Cadet"

	PersonnelStatusActive = "active"

	TrainingStatusCompleted = "completed"
)
