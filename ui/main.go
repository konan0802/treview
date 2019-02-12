package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/inabajunmr/treview/github/trending"
	treview "github.com/inabajunmr/treview/service"
	"github.com/zserge/lorca"
)

type Condition struct {
	Span    string
	Lang    string
	OnlyNew bool
}

var ui lorca.UI

func main() {
	ui, _ = lorca.New("", "", 1280, 800)

	defer ui.Close()

	err := ui.Bind("load", load)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = ui.Bind("reload", reloadRepositories)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer ln.Close()

	go serveContents(ln)

	err = ui.Load(fmt.Sprintf("http://%s", ln.Addr()))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer ui.Close()

	<-ui.Done()

}

func serveContents(ln net.Listener) {
	err := http.Serve(ln, http.FileServer(FS))
	ui.Eval(`console.log("served");`)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func load() {

	span := trending.GetSpanByString("today")
	repos := treview.GetRepositories(span, treview.GetLangs(""), true)
	bindRepositories(repos)

	bindLangs(trending.FindLangs())
}

func reloadRepositories(cond Condition) {
	langs := []string{cond.Lang}
	span := trending.GetSpanByString(cond.Span)
	repos := treview.GetRepositories(span, langs, cond.OnlyNew)
	bindRepositories(repos)
}

func bindRepositories(repos []trending.Repository) {
	val, _ := json.Marshal(repos)
	ui.Eval("vm.repos = " + string(val[:]))
}

func bindLangs(langs []string) {
	val, _ := json.Marshal(langs)
	ui.Eval("vm.langs = " + string(val[:]))
}