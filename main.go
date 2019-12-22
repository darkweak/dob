package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

type certificate struct {
	Cert   string   `yaml:"certificate"`
	Domain string   `yaml:"domain"`
	Sans   []string `yaml:"SANs"`
	Key    string   `yaml:"key"`
}

type Config struct {
	Certificates []certificate `yaml:"certificates"`
	Email        string        `yaml:"email"`
}

func (c *Config) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}

func getConfig() string {
	return "./config.yml"
}

func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func watchConfig(w *fsnotify.Watcher) {
	data := readFile(getConfig())
	var config Config
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", config)
	acmeFile := acmeFileBuilder(config.Email)

	for _, item := range config.Certificates {
		item.initWatcher(&config, acmeFile, w)
	}
}

func main() {
	w, err := fsnotify.NewWatcher()
	watchConfig(w)

	if err != nil {
		log.Fatal(err)
	}

	watcherWithCallback(w, func() {watchConfig(w)}, []string{getConfig()})
}

func watcherWithCallback (watcher *fsnotify.Watcher, callback func(), files []string) {
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if (event.Op&fsnotify.Rename == fsnotify.Rename) ||
					(event.Op&fsnotify.Write == fsnotify.Write) {
						callback()
						watcher.Add(event.Name)
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	for _, f := range files {
		err := watcher.Add(f)
		if err != nil {
			return
		}
	}
	<-done
}

func (c *certificate) initWatcher(config *Config, acmeFile *AcmeFile, w *fsnotify.Watcher) {
	acmeFile.addCertificate(c)
	watcherWithCallback(w, func () {acmeFile.addCertificate(c)}, []string{c.Cert, c.Key})
}
