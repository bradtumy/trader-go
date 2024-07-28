package main

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Stock struct {
	ID          int     `json:"id"`
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	TotalShares int     `json:"total_shares"`
}

type Orders struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Symbol   string `json:"symbol"`
	Shares   int    `json:"shares"`
}

type Config struct {
	Server struct {
		Port        string `yaml:"port"`
		Host        string `yaml:"host"`
		Certificate string `yaml:"server_cert"`
		Key         string `yaml:"server_key"`
	} `yaml:"server"`
	Database struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}
