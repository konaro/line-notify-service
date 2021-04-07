package data

type Security struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Default  bool   `json:"default"`
}
