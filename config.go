package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Target string `yaml:"target"`
	Proxy  string `yaml:"proxy"`
	Listen struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
	}
	Redis *struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
		DB   int   `yaml:"db"`
		TLS  int64  `yaml:"tls"`
	}
}

var config = Config{}

func init() {
	configFile, err := os.Open("odoh.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
}