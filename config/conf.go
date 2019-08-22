package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type AmongConfig struct {
	ETCDAddr []string `yaml:"etcdaddr"`
}

func NewAmongConfig(path string) *AmongConfig {
	var filepath string
	if path != "" {
		filepath = path
	} else {
		filepath = "./among.yaml"
	}
	fmt.Println("filepath is ", filepath)
	confFile, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("open %s err %s\n", filepath, err)
		return nil
	}
	defer confFile.Close()
	conf := new(AmongConfig)
	decoder := yaml.NewDecoder(confFile)
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Printf("decode %s err %s\n", filepath, err)
		return nil
	}
	fmt.Printf("%v\n", conf.ETCDAddr)
	return conf
}
