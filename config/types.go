package config

type Config struct {
	Input  string `json:"input"`
	Output string `json:"output"`

	Ent *Ent `json:"ent"`
}

type Ent struct {
	Package     string `json:"package"`
	Graphql     bool   `json:"graphql"`
	Auth        bool   `json:"auth"`
	Privacy     bool   `json:"privacy"`
	PrivacyNode bool   `json:"-"`
	FileUpload  bool   `json:"fileUpload"`
	Debug       bool   `json:"debug"`
	Database    string `json:"database"` //sqlite3, mysql, postgres
}

type Gorm struct {
	Package    string `json:"package"`
	Auth       bool   `json:"auth"`
	Graphql    bool   `json:"graphql"`
	FileUpload bool   `json:"fileUpload"`
	Debug      bool   `json:"debug"`
	Database   string `json:"database"` //sqlite3, mysql, postgres
}

type Prisma struct {
	Database string `json:"database"` //sqlite3, mysql, postgres
}
