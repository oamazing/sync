package config

type Config struct {
	BasePath string `yaml:"base_path"`
	Storage  string `yaml:"storage"`
	Qiniu    Qiniu  `yaml:"qiniu"`
	Ali      Ali    `yaml:"ali"`
}

type Qiniu struct {
	Ak     string `yaml:"ak"`
	Sk     string `yaml:"sk"`
	Bucket string `yaml:"bucket"`
	Url    string `yaml:"url"`
	Zone   string `yaml:"zone"`
}

type Ali struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}
