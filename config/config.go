package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var conf *Config

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadFile(filepath.Join(dir, "config.yml"))
	if err != nil {
		log.Panic(err)
	}
	if err = yaml.Unmarshal(bs, &conf); err != nil {
		log.Panic(err)
	}
}

func GetConfig() *Config {
	return conf
}

func GetQiniu() *Qiniu {
	return &conf.Qiniu
}
