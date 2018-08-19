package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/pascalw/go-alfred"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		debug       = kingpin.Flag("debug", "Show debugging output").Bool()
		repos       = kingpin.Command("repos", "Show repositories from cache")
		reposFilter = repos.Arg("filter", "Fuzzy match the full repository name").String()
		login       = kingpin.Command("login", "Login to github via oauth")
		update      = kingpin.Command("update", "Updates repositories from Github")
	)

	cmd := kingpin.Parse()

	if *debug == false {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}

	switch cmd {
	case repos.FullCommand():
		reposCommand(*reposFilter)
	case login.FullCommand():
		if err := loginCommand(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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

func backgroundUpdate() error {
	cmd := exec.Command(os.Args[0], "update")
	if err := cmd.Start(); err != nil {
		return err
	}
	log.Printf("Background pid %#v", cmd.Process.Pid)
	return nil
}
