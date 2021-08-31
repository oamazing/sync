package config

type Config struct {
	BasePath string `yaml:"base_path"`
	Storage  string `yaml:"storage"`
	Qiniu    Qiniu  `yaml:"qiniu"`
}

type Qiniu struct {
	Ak     string `yaml:"ak"`
	Sk     string `yaml:"sk"`
	Bucket string `yaml:"bucket"`
}
