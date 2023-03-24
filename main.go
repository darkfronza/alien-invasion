package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/darkfronza/alien-invasion/simulation"
	"github.com/darkfronza/alien-invasion/worldmap"
)

func main() {
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
