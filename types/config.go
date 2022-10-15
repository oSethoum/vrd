package types

type Config struct {
	Input    string `yaml:"input"`
	Debug    bool   `yaml:"debug"`
	Database string `yaml:"database"`
	Ent      Ent    `yaml:"ent"`
}

type Ent struct {
	Output      string  `yaml:"output"`
	Package     string  `yaml:"package"`
	Graphql     Graphql `yaml:"graphql"`
	Privacy     bool    `yaml:"privacy"`
	PrivacyNode bool    `yaml:"-"`
	Auth        bool    `yaml:"auth"`
}

type Graphql struct {
	Subscription    bool `yaml:"subscription"`
	FileUpload      bool `yaml:"file_upload"`
	RelayConnection bool `yaml:"relay_connection"`
}
