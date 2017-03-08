package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/user"
	"path"
	"time"
)

var (
	rootPrefix = flag.String("root", "", "root directory for config & stuffs")
	cfgPath    = flag.String("cfg", "", "cfg path, defaults to <root>/config.json")
	htmlRoot   = flag.String("html", "", "path to html templates, defaults to <root>/html/")

	debug = flag.Bool("debug", false, "enable debug mode (e.g. file generation instead of mail)")

	cfg *Config
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
	var err error
	cfg, err = LoadConfigFile(*cfgPath)
	if err != nil {
		log.Fatal("in LoadConfig():", err)
	}
	if *htmlRoot != "" {
		cfg.HtmlRoot = *htmlRoot
	}

	for _, q := range cfg.WatchList {
		q.RawUrl = q.BuildURL()
		log.Println(q.RawUrl)
		items, err := q.Run()
		if err != nil {
			log.Printf("error running query: %s\n%s", q.RawUrl, err)
			continue
		}
		log.Printf("found %d items\n", len(items))
	}

	// generate html content
	wr_html := bytes.NewBuffer([]byte{})
	tpl, err := template.ParseGlob(path.Join(cfg.HtmlRoot, "*html"))
	if err != nil {
		log.Fatalf("couldn't parse html folder: %s", err)
	}

	err = tpl.Execute(wr_html, cfg)
	if err != nil {
		log.Fatalf("error in tpl.Execute: %s", err)
	}

	if *debug {
		fmt.Fprintf(os.Stdout, wr_html.String())
		os.Exit(0)
	}

	err = cfg.SMTPConfig.SendMail("lbcbot> "+time.Now().String(), wr_html.String())
	if err != nil {
		log.Fatalf("error sending mail: %s", err)
	}
	os.Exit(0)
}
