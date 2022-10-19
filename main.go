package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/briandowns/spinner"
)

func main() {
	parser := argparse.NewParser("wallsw", "Switches between wallpapers")

	directory := parser.String("d", "directory", &argparse.Options{Required: true, Help: "The directory to look for wallpapers in"})
	startAt := parser.Int("", "start-at", &argparse.Options{Required: false, Help: "The number of files into the directory to start at"})
	filter := parser.String("f", "filter", &argparse.Options{Required: false, Help: "A string to filter files by"})
	random := parser.Flag("", "random", &argparse.Options{Required: false, Help: "Should it fetch a random wallpaper?"})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatal(err)
	}

	if *random == false {
		catalog(*directory, *startAt, *filter)	

		return
	}

	randomWallpaper(*directory, *filter) 
}

func catalog(dir string, startAt int, filter string) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal("There was an error fetching wallpapers from that path!")
	}

	n := 0
	for _, file := range entries {
		if strings.Contains(file.Name(), filter) {
			entries[n] = file
			n++
		}
	}

	entries = entries[:n]

	for _, file := range entries[startAt:] {
		if file.IsDir() {
			continue
		}

		s := spinner.New(spinner.CharSets[39], 100 * time.Millisecond)

		s.Start()

		s.Suffix = fmt.Sprintf(" Switching wallpaper to %s...", file.Name())

		home, err := os.UserHomeDir()

		if err != nil {
			log.Fatal("Could not find your home directory!")
		}

		cmd := exec.Command(
			"python3", 
			fmt.Sprintf("%s/wallpaper/switch.py", home),
			fmt.Sprintf("%s/%s", dir, file.Name(),
		))

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		if waitErr := cmd.Wait(); waitErr != nil {
			log.Fatal("There was an error switching the wallpaper!")
		}

		s.Stop()

		fmt.Println("Do you want to switch to the next wallpaper? (yes/no)")

		reader := bufio.NewReader(os.Stdin)

		text, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("yes", text) == 0 {
			continue
		}
		
		return
	}
}

func randomWallpaper(dir string, filter string) {
	entries, err := os.ReadDir(dir)

	n := 0
	for _, file := range entries {
		if strings.Contains(file.Name(), filter) {
			entries[n] = file	
			n++
		}
	}

	entries = entries[:n]

	if err != nil {
		log.Fatal("There was an error fetching wallpapers from that path!")						
	}

	wallpaper := getWallpaper(entries)

	s := spinner.New(spinner.CharSets[43], 100 * time.Millisecond)

	s.Start()

	s.Suffix = fmt.Sprintf(" Randomizing wallpaper...")
	s.FinalMSG = fmt.Sprintf("Wallpaper randomized to %s.\n", wallpaper.Name())

	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatal("could not find your home directory!")
	}

	cmd := exec.Command(
		"python3",
		fmt.Sprintf("%s/wallpaper/switch.py", home),
		fmt.Sprintf("%s/%s", dir, wallpaper.Name()),
	)

	if err := cmd.Start(); err != nil {
		log.Fatal("There was an error randomizing the wallpaper!")
	}

	if waitErr := cmd.Wait(); waitErr != nil {
		log.Fatal("There was an error switching the wallpaper!")
	}

	s.Stop()
}

func getWallpaper(entries []os.DirEntry) os.DirEntry {
	rand.Seed(time.Now().UnixNano())
	wallpaper := entries[rand.Intn(len(entries))]	

	if wallpaper.IsDir() {
		wallpaper = getWallpaper(entries)
	}

	return wallpaper
}
