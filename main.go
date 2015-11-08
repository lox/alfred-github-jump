package main

import (
	"io/ioutil"
	"log"

	"github.com/pascalw/go-alfred"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		debug        = kingpin.Flag("debug", "Show debugging output").Bool()
		repos        = kingpin.Command("repos", "Show repositories from cache")
		reposFilters = repos.Arg("filter", "Filter strings").Strings()
		login        = kingpin.Command("login", "Login to github via oauth")
		update       = kingpin.Command("update", "Updates repositories from Github")
	)

	cmd := kingpin.Parse()

	if *debug == false {
		log.SetOutput(ioutil.Discard)
	}

	switch cmd {
	case repos.FullCommand():
		reposCommand(*reposFilters)
	case login.FullCommand():
		loginCommand()
	case update.FullCommand():
		updateCommand()
	}
}

func alfredError(err error) *alfred.AlfredResponseItem {
	return &alfred.AlfredResponseItem{
		Valid:    false,
		Uid:      "error",
		Title:    "Error Occurred",
		Subtitle: err.Error(),
		Icon:     "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/AlertStopIcon.icns",
	}
}
