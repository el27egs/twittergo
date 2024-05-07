package models

type Settings struct {
	DbHost     string `json:"DbHostname"`
	DbName     string `json:"DbName"`
	DbUsername string `json:"DbUsername"`
	DbPassword string `json:"DbPassword"`
	JwtSeed    string `json:"JwtSeed"`
}
