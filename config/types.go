package config

type Config struct {
	Input  string `json:"input"`
	Output string `json:"output"`

	Ent  *Ent `json:"ent"`
	Gorm *Ent `json:"gorm"`
}

type Ent struct {
	Package string `json:"package"`
	Graphql bool   `json:"graphql"`
	Echo    bool   `json:"echo"`
	Auth    bool   `json:"auth"`
	Privacy bool   `json:"privacy"`
}

type Gorm struct {
}
