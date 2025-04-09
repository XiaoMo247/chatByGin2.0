package models

type ChatData struct {
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Content  string   `json:"content"`
	Time     string   `json:"time"`
	UserList []string `json:"user_list"`
}
