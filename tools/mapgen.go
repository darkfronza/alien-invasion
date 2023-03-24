// Simple tool to generate random maps which could be used for testing the alien-invasion simulator.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func usage() {
	fmt.Printf("Usage: %s <n_cities>\n", os.Args[0])
	fmt.Println("\nGenerates a random map with the number of cities (approximated) provided by the <n-cities> argument.")
	fmt.Println("\nOutput is sent to stdout by default, use > to redirect it to a file.")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	rand.Seed(time.Now().UnixNano())

	nCities, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("[*] ERROR: Invalid number:", os.Args[1])
		fmt.Println()
		usage()
	}

	// keep already generated cities to avoid duplication
	cityDB := make(map[string]bool)

	for nCities > 0 {
		cityName := addRandomCity(cityDB)

		if cityName == "" {
			fmt.Println("Oops! I can't generate a unique city name (perhaps n-cities is too big?)")
			os.Exit(1)
		}

		cityDB[cityName] = true
		fmt.Print(cityName)
		nCities -= 1

		nNeighbors := 1
		if nCities >= 4 {
			nNeighbors = 1 + rand.Intn(4)
		}
		nCities -= nNeighbors
		dirs := []string{"north", "south", "east", "west"}
		var stack []string

		for ; nNeighbors > 0; nNeighbors-- {
			neighbor := addRandomCity(cityDB)
			if neighbor == "" {
				fmt.Println("Oops! I can't generate a unique city name (perhaps n-cities is too big?)")
				os.Exit(1)
			}
			x := rand.Intn(len(dirs))
			dir := dirs[x]
			dirs = append(dirs[:x], dirs[x+1:]...)
			fmt.Printf(" %s=%s", dir, neighbor)
			// keep track of how to back from the neighbor to the current city
			reverseCityDir := fmt.Sprintf("%s %s=%s", neighbor, reverseDir(dir), cityName)
			stack = append(stack, reverseCityDir)
		}

		fmt.Println()
		for _, reverseCity := range stack {
			fmt.Println(reverseCity)
		}
	}
}

func randomCityName(minLength, maxLength int) string {
	var sb strings.Builder

	length := minLength + rand.Intn(maxLength-minLength+1)

	alphabet := "abcdefghijklmnopqrstuvwxyz"
	al := len(alphabet)

	for i := 0; i < length; i++ {
		k := rand.Intn(al)
		if i == 0 {
			sb.WriteByte(byte(unicode.ToUpper(rune(alphabet[k]))))
		} else {
			sb.WriteByte(alphabet[k])
		}
	}

	return sb.String()
}

func addRandomCity(db map[string]bool) string {
	for i := 0; i < 1000; i++ {
		c := randomCityName(3, 6)
		if _, exists := db[c]; !exists {
			db[c] = true
			return c
		}
	}

	return ""
}

func reverseDir(direction string) string {
	switch direction {
	case "north":
		return "south"
	case "south":
		return "north"
	case "west":
		return "east"
	case "east":
		return "west"
	default:
		fmt.Printf("Unknown direction: [%s] aborting...\n", direction)
		os.Exit(1)
	}
	return ""
}
