package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var parser = flags.NewParser(&options, flags.Default)

func readFile(cfg *Config) {
	f, err := os.Open(options.Configfile)
	if err != nil {
		checkErr(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		checkErr(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func debug(text string) {

	if options.Verbose {
		log.Println(text)
	}

}
