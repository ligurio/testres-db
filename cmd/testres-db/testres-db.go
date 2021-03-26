package main

import (
	"flag"
	"fmt"
	"github.com/ligurio/testres-db/backends"
	"github.com/ligurio/testres-db/db"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var options struct {
	configName  string
	dbName      string
	projectName string
	buildsNumber int
}

func init() {
	flag.StringVar(&options.configName, "config", "testres-db.yaml", "config")
	flag.StringVar(&options.dbName, "db", "testres.sqlite", "database file")
	flag.StringVar(&options.projectName, "project", "", "project name")
	flag.IntVar(&options.buildsNumber, "limit", -1, "limit a number of builds")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Import test results to an SQLite DB.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}
}

type Project struct {
	Name     string
	Backends []backends.Backend
}

type Config struct {
	Projects []Project `projects`
}

func (c *Config) getConf(configName *string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(*configName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	flag.Parse()

	log.Println("Using database:", options.dbName)
	log.Println("Using config:", options.configName)

	db := &db.DB{}
	db.SetPath(options.dbName)

	var initDb bool = false
	if _, err := os.Stat(options.dbName); os.IsNotExist(err) {
		initDb = true
	}

	err := db.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if initDb {
		err = db.Init()
		if err != nil {
			log.Fatal(err)
		}
	}

	var c Config
	_, err = c.getConf(&options.configName)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range c.Projects {
		log.Println("======> Processing project", p.Name)
		if options.projectName != "" {
			if options.projectName == p.Name {
				for _, b := range p.Backends {
					results, err := b.GetTestResults(options.buildsNumber)
					if err != nil {
						log.Println(err)
						continue
					}
					if results == nil {
						continue
					}
					err = db.AddResults(p.Name, results)
					if err != nil {
						log.Println(err)
					}
				}
			}
		} else {
			for _, b := range p.Backends {
				results, err := b.GetTestResults(options.buildsNumber)
				if err != nil {
					log.Println(err)
					continue
				}
				if results == nil {
					continue
				}
				err = db.AddResults(p.Name, results)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
