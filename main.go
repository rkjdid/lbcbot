package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

var (
	rootPrefix = flag.String("root", "", "root directory for config & stuffs")
	cfgPath    = flag.String("cfg", "", "cfg path, defaults to <root>/config.json")
	htmlRoot   = flag.String("html", "", "path to html templates, defaults to <root>/html/")
)

func init() {
	flag.Parse()

	if *rootPrefix == "" || *rootPrefix == "." {
		cwd, err := os.Executable()
		if err != nil {
			log.Println("couldn't retreive working dir", err)
			usr, err := user.Current()
			if err != nil {
				log.Fatal("couldn't retreive system user", err)
			}
			*rootPrefix = usr.HomeDir
		} else {
			*rootPrefix = path.Dir(cwd)
		}
	}

	if *cfgPath == "" {
		*cfgPath = path.Join(*rootPrefix, "config.json")
	}
}

func main() {
	cfg, err := LoadConfigFile(*cfgPath)
	if err != nil {
		log.Fatal("in LoadConfig():", err)
	}

	for _, q := range cfg.WatchList {
		items, err := q.Run()
		if err != nil {
			fmt.Printf("error running query: %s\n%s", q.RawUrl, err)
			continue
		}

		fmt.Printf("found %d items for %s:\n", len(items), q)
		for k, item := range items {
			if k == 5 {
				fmt.Printf("  ... %d more\n", len(items[k:]))
				break
			}
			fmt.Printf("  - %#v\n", item)
		}
		fmt.Println()
	}
}
