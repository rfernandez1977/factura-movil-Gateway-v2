package models

type EmailEnviado struct {
	ID      string `json:"id"`
	Para    string `json:"para"`
	Asunto  string `json:"asunto"`
	Mensaje string `json:"mensaje"`
}
