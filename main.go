package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/darkfronza/alien-invasion/simulation"
	"github.com/darkfronza/alien-invasion/worldmap"
)

var Version string

func main() {
	// TODO: add daemon feature to test systemd on nix
	if len(os.Args) == 2 && os.Args[1] == "daemon" {
		fmt.Printf("daemon mode..... version=%s\n", Version)
		time.Sleep(50 * time.Minute)
		os.Exit(0)
	}

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "version") {
		fmt.Printf("alien-invasion simulator %s\n", Version)
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <map-file> <n-aliens>\n", os.Args[0])
		os.Exit(0)
	}

	nAliens, err := strconv.Atoi(os.Args[2])
	if err != nil || nAliens <= 0 {
		fmt.Printf("Invalid number:%d, <n-aliens> must be > 0\n", nAliens)
		os.Exit(1)
	}

	// Open map file for reading
	fp, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer fp.Close()

	worldMap, err := worldmap.Load(fp)
	if err != nil {
		log.Fatal(err)
	}

	alienInvasionSimul, err := simulation.New(worldMap, nAliens)
	if err != nil {
		log.Fatal(err)
	}

	for alienInvasionSimul.Step() {
	}

	if !worldMap.IsDestroyed() {
		fmt.Println("\nSome cities survived the alien attack, let's celebrate!")
		// Print final world state
		fmt.Println("\nWorld map after the end of alien invasion:")
		fmt.Println("-----------------------------------------------------------------")
		worldMap.Print()
		fmt.Println("-----------------------------------------------------------------")
	} else {
		fmt.Println("\nUnfortunately, the world was totally destroyed by the aliens :(\n\nSee you in the heavens!")
	}
}
