package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ShowBody     bool          `yaml:"showBody"`
	ProxyConfigs []ProxyConfig `yaml:"proxyConfigs"`
}

type ProxyConfig struct {
	Port int    `yaml:"port"`
	Url  string `yaml:"url"`
}

func parseOptions() Config {
	showBody := flag.Bool("showBody", false, "shows the request and response bodies")
	port := flag.Int("port", 8080, "proxy port")
	url := flag.String("url", "http://www.example.com", "url")
	file := flag.String("file", "none", "config file")

	flag.Parse()

	singleProxyConfig := ProxyConfig{
		Port: *port,
		Url:  *url,
	}

	if *file != "none" {
		yamlFile, err := ioutil.ReadFile(*file)
		if err != nil {
			panic(err)
		}
		var config Config
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
		return config
	}

	return Config{
		ShowBody:     *showBody,
		ProxyConfigs: []ProxyConfig{singleProxyConfig},
	}
}
