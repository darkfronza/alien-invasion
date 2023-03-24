// Package simulation provides a way for simulating an alien invastion

package simulation

import (
	"fmt"
	"math/rand"

	"github.com/darkfronza/alien-invasion/worldmap"
)

// Alien identification, for now a simple integer.
type alienID int

// An Alien has an ID and a city where he's currently trying to dominate
type alien struct {
	id          alienID
	currentCity string
	nSteps      int // Number of successful movements performed by this beautiful Alien
}

// AlienInvasionSimul keeps worldwide alien invasion state
type AlienInvasionSimul struct {
	aliens   map[alienID]*alien
	worldMap *worldmap.WorldMap
}

// New constructs a new alien simulation
func New(m *worldmap.WorldMap, nAliens int) (*AlienInvasionSimul, error) {
	if nAliens == 0 {
		return nil, fmt.Errorf("No aliens, no simulation")
	}

	// Shuffle cities to randomly distribute the aliens
	shuffledCities := make([]string, len(m.GetAllCityNames()))
	copy(shuffledCities, m.GetAllCityNames())
	rand.Shuffle(len(shuffledCities), func(i, j int) {
		shuffledCities[i], shuffledCities[j] = shuffledCities[j], shuffledCities[i]
	})

	// Alien mothership would like to drop 1 alien per city
	aliensPerCity := 1
	if nAliens > len(shuffledCities) {
		// Not enough cities, mothership needs to drop more than 1 alien per city
		aliensPerCity = nAliens / len(shuffledCities)
	}

	simul := &AlienInvasionSimul{
		aliens:   make(map[alienID]*alien),
		worldMap: m,
	}

	// Drop aliens into random cities
	currentCity := 0
	currAlienID := 0
	for nAliens > 0 && currentCity < len(shuffledCities) {
		targetCity := m.GetCity(shuffledCities[currentCity])
		for i := 0; i < aliensPerCity; i++ {
			simul.aliens[alienID(currAlienID+1)] = &alien{
				id:          alienID(currAlienID + 1),
				currentCity: targetCity.Name,
			}
			currAlienID++
		}
		currentCity++
		nAliens -= aliensPerCity
	}

	// Drop remaining aliens on random cities (if #Aliens > #Cities)
	for nAliens > 0 {
		nAliens--
		k := rand.Intn(len(shuffledCities))
		targetCity := m.GetCity(shuffledCities[k])
		simul.aliens[alienID(currAlienID+1)] = &alien{
			id:          alienID(currAlienID + 1),
			currentCity: targetCity.Name,
		}
		currAlienID++
	}

	return simul, nil
}

// Step advances the alien invasion
// Each alien moves from its current city to a random neighbor one (if available)
// All aliens move at the same time/step (because they have a chip with a shared cpu clock controlled by the mothership)
// Step returns false when either condition below takes place:
//   - 1: All aliens were destroyed
//   - 2: Each alien has moved at least 10000 times
func (s *AlienInvasionSimul) Step() bool {
	// If at least one alien didn't move at least 10000 steps, then at least one alien is not stuck.
	allAliensAreStuck := true

	if len(s.aliens) == 0 {
		// All aliens are dead
		fmt.Println("All aliens were mutually destroyed!")
		return false
	}

	// Map cities to the aliens that have recently moved to it
	targetCities := make(map[string][]*alien)

	// Try to move all aliens to a random neighbor city
	for _, alien := range s.aliens {
		city := s.worldMap.GetCity(alien.currentCity)
		if city == nil || city.IsDestroyed() {
			delete(s.aliens, alien.id)
			continue
		}

		if targetCity := city.RandomNeighborCity(); targetCity != nil {
			alien.currentCity = targetCity.Name
			alien.nSteps++
			if alien.nSteps < 10000 {
				allAliensAreStuck = false
			}
			// Add the alien to the target city
			targetCities[targetCity.Name] = append(targetCities[targetCity.Name], alien)
		}
	}

	// Destroy the cities where two or more aliens have landed.
	for city, aliens := range targetCities {
		if len(aliens) > 1 {
			fmt.Printf("%s has been destroyed due to an alien conflict between", city)
			for i, a := range aliens {
				if i > 0 {
					fmt.Print(" and")
				}
				fmt.Printf(" Alien#%d", a.id)
				// Destroy the alien
				delete(s.aliens, a.id)
			}
			fmt.Println()

			// Destroy the city
			s.worldMap.DestroyCity(city)
		}
	}

	if allAliensAreStuck && len(s.aliens) > 0 {
		fmt.Printf("All remaining aliens are trapped! /%d\n", len(s.aliens))
	}

	return !allAliensAreStuck
}
