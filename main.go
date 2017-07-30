package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/rkjdid/util"
	"html/template"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"
)

const CachePath = "cache.json"

var (
	rootPrefix = flag.String("root", "", "root directory for config & stuffs")
	cfgPath    = flag.String("cfg", "", "cfg path, defaults to <root>/config.json")
	htmlRoot   = flag.String("html", "", "path to html templates, defaults to <root>/html/")

	debug = flag.Bool("debug", false, "enable debug mode (e.g. file generation instead of mail)")

	cfg   *Config
	nbRun int
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
		*cfgPath = filepath.Join(*rootPrefix, "config.json")
	}
	if *htmlRoot == "" {
		*htmlRoot = filepath.Join(*rootPrefix, "html")
	}
}

func run() {
	var nbResults int
	for k, q := range cfg.WatchList {
		q.RawUrl = q.BuildURL()

		items, err := q.Run()
		if err != nil {
			log.Printf("error running query: %s\n%s", q.RawUrl, err)
			continue
		}
		if sz := len(items); sz > 0 {
			nbResults += sz
			log.Printf("found %d new items for query #%d \"%s\"", sz, k, q.Search)
		}
	}

	if nbResults == 0 {
		log.Printf("%d: no new results", nbRun)
		return
	}
	log.Printf("%d: found %d new results. sending notification", nbRun, nbResults)

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

	// if debug, just tick once and exit
	if *debug {
		fmt.Fprintf(os.Stdout, wr_html.String())
		fmt.Fprintln(os.Stderr, "debug enabled. exitting now")
		os.Exit(0)
	}

	err = cfg.SMTPConfig.SendMail("lbcbot> "+time.Now().String(), wr_html.String())
	if err != nil {
		log.Fatalf("error sending mail: %s", err)
	}
}

func main() {
	var err error
	err = util.ReadGenericFile(&cfg, *cfgPath)
	cfg.HtmlRoot = *htmlRoot
	if err != nil {
		log.Fatal("in LoadConfig():", err)
	}

	cfg.cache = make(map[string]bool)
	err = util.ReadGenericFile(&cfg.cache, CachePath)
	if err != nil {
		log.Printf("couldn't load cache: %s", err)
	}

	for _, q := range cfg.WatchList {
		q.cfg = cfg
	}
	tick := time.NewTicker(time.Minute * time.Duration(cfg.PollIntervalMin))
	for {
		run()
		nbRun++
		<-tick.C
	}

	os.Exit(0)
}
