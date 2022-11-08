package types

type Config struct {
	Input    string `yaml:"input"`
	Debug    bool   `yaml:"debug"`
	Database string `yaml:"database"`
	Ent      *Ent   `yaml:"ent"`
	Gorm     *Gorm  `yaml:"gorm"`
}

type Ent struct {
	Output      string   `yaml:"output"`
	Package     string   `yaml:"package"`
	Graphql     *Graphql `yaml:"graphql"`
	Privacy     bool     `yaml:"privacy"`
	PrivacyNode bool     `yaml:"privacy_node"`
	Auth        bool     `yaml:"auth"`
	Echo        bool     `yaml:"echo"`
}

type Gorm struct {
	Output     string `yaml:"output"`
	Package    string `yaml:"package"`
	Fiber      bool   `yaml:"fiber"`
	Auth       bool   `yaml:"auth"`
	Encryption bool   `yaml:"encryption"`
	FileUpload bool   `yaml:"file_uplaod"`
	Socket     bool   `yaml:"socket"`
	GormModel  bool   `yaml:"gorm_model"`
	Swagger    bool   `yaml:"swagger"`
}

type Graphql struct {
	Subscription    bool `yaml:"subscription"`
	FileUpload      bool `yaml:"file_upload"`
	RelayConnection bool `yaml:"relay_connection"`
}
