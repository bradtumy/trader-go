package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("resources/properties.yml")
	if err != nil {
		fmt.Println("ERROR: I couldn't read the properties file.")
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		fmt.Println("ERROR: I couldn't decode the YAML.")
		processError(err)
	}
}
