package config

type Config struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output"`

	Ent *Ent `yaml:"ent"`
}

type Ent struct {
	Package     string `yaml:"package"`
	Graphql     bool   `yaml:"graphql"`
	Auth        bool   `yaml:"auth"`
	Privacy     bool   `yaml:"privacy"`
	PrivacyNode bool   `yaml:"-"`
	FileUpload  bool   `yaml:"fileUpload"`
	Debug       bool   `yaml:"debug"`
	Database    string `yaml:"database"` //sqlite3, mysql, postgres
}

type Gorm struct {
	Package    string `yaml:"package"`
	Auth       bool   `yaml:"auth"`
	Graphql    bool   `yaml:"graphql"`
	FileUpload bool   `yaml:"fileUpload"`
	Debug      bool   `yaml:"debug"`
	Database   string `yaml:"database"` //sqlite3, mysql, postgres
}

type Prisma struct {
	Database string `yaml:"database"` //sqlite3, mysql, postgres
}
