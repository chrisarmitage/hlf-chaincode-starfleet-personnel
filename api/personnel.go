package api

type Personnel struct {
	PersonnelID string `json:"personnelID"`
	Name        string `json:"name"`
	Rank        string `json:"rank"`
	Campus      string `json:"campus"`
	Status      string `json:"status"`
}
