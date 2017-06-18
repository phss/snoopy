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

func parseConfigOptions() Config {
	showBody := flag.Bool("showBody", false, "shows the request and response bodies")
	port := flag.Int("port", 8080, "proxy server port")
	url := flag.String("url", "http://www.example.com", "proxied url")
	file := flag.String("file", "", "configuration file")
	flag.Parse()

	if *file != "" {
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
	} else {
		return Config{
			ShowBody: *showBody,
			ProxyConfigs: []ProxyConfig{
				ProxyConfig{
					Port: *port,
					Url:  *url,
				},
			},
		}
	}
}
