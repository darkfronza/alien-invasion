package worldmap

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

type Direction int

// Enumerate directions
const (
	Invalid     Direction = -1
	North                 = 0
	South                 = 1
	West                  = 2
	East                  = 3
	NDirections           = 4
)

// Map direction contants above to their string representation
var dirNames = []string{
	"north", "south", "west", "east",
}

// City represents a possible target city for invasion by the aliens
// We map the city's neighbors as an array of pointers to another City, where
// each position in the array (North=0, South=1, ...) is the neighbor City.
// A nil pointer indicates that the city has no connection to another city under that direction.
type City struct {
	Name string

	// Each direction leads to another city if directions[dir] != nil
	directions [4]*City

	// Is the city destroyed?
	destroyed bool

	// Is the city a neighbor one? (only found as neighbor of other cities in the input map)
	// e.g, for the sample map file below:
	// Foo north=Moon west=Silverado
	// Moon south=Foo
	// Silverado would have isNeighbor=true, as it's only found as a neighboir reference.
	isNeighbor bool
}

// IsDestroyed checks if the city was destroyed
func (c *City) IsDestroyed() bool {
	return c.destroyed
}

// RandomNeighborCity gets a random neighbor city, if available.
// Returns a pointer to the random neighbor city, or nil if none is available.
func (c *City) RandomNeighborCity() *City {
	var dirs []int

	for i := 0; i < NDirections; i++ {
		if c.directions[i] != nil && !c.directions[i].destroyed {
			dirs = append(dirs, i)
		}
	}

	if len(dirs) > 0 {
		k := rand.Intn(len(dirs))
		return c.directions[dirs[k]]
	}

	return nil
}

// WorlMap represents all the possible cities in the world, as loaded by the map file.
type WorldMap struct {
	// Map city name to it's City object
	cities map[string]*City
	// Keep an separate array of city names to facilitate processing
	cityNames []string
}

// Load loads the map from the input stream provided.
func Load(r io.Reader) (*WorldMap, error) {
	worldMap := &WorldMap{
		cities: make(map[string]*City),
	}

	ln := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		ln++
		// Skip empty lines
		if line == "" {
			continue
		}

		tokens := strings.Fields(line)
		cityName := tokens[0]
		if city, ok := worldMap.cities[cityName]; ok && !city.isNeighbor {
			// Already added city
			fmt.Printf("worldMap.Load(): found duplicated city! name=%s, line=%d", cityName, ln)
			continue
		}

		// Create a new city
		city := &City{
			Name: cityName,
		}
		foundError := false

		var neighborCities []*City
		// Parse directions to another cities.
		for i := 1; i < len(tokens); i++ {
			dirData := strings.Split(tokens[i], "=")
			if len(dirData) != 2 {
				fmt.Printf("worldMap.Load(): Invalid direction format! found:%s, expected:dirname=city, at line=%d\n", tokens[i], ln)
				foundError = true
				continue
			}
			direction, derr := parseDirection(dirData[0])
			if derr != nil {
				fmt.Printf("worldMap.Load(): Invalid direction:%s, at line=%d\n", dirData[0], ln)
				foundError = true
				continue
			}
			neighborCity := &City{
				Name:       dirData[1],
				isNeighbor: true,
			}
			// Add the neighbor city to the listed city's direction
			city.directions[direction] = neighborCity
			neighborCities = append(neighborCities, neighborCity)
		}
		if !foundError {
			for _, neighborCity := range neighborCities {
				// Check if the neighbor city was already add to the worlmap
				// If not, add it as a neighbor city.
				if _, ok := worldMap.cities[neighborCity.Name]; !ok {
					worldMap.cityNames = append(worldMap.cityNames, neighborCity.Name)
					worldMap.cities[neighborCity.Name] = neighborCity
				}
			}
			if _, ok := worldMap.cities[cityName]; !ok {
				worldMap.cityNames = append(worldMap.cityNames, cityName)
			}
			// Add city to the worldmap, if a neighbor one was added before, replace it as a non-neighbor
			worldMap.cities[cityName] = city
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return worldMap, nil
}

// GetCity returns a city from the worldmap by its name.
func (m *WorldMap) GetCity(name string) *City {
	city, ok := m.cities[name]
	if !ok || city.destroyed {
		return nil
	}
	return city
}

// GetAllCityNames return all loaded cities' names.
func (m *WorldMap) GetAllCityNames() []string {
	return m.cityNames
}

// DestroyCity marks a city as destroyed
func (m *WorldMap) DestroyCity(name string) {
	if city, ok := m.cities[name]; ok {
		city.destroyed = true
	}
}

// Print prints the current world map to the stdout, in the same format loaded by the input map file.
func (m WorldMap) Print() {
	for _, city := range m.cities {
		if city.destroyed || city.isNeighbor {
			continue
		}

		fmt.Printf("%s", city.Name)

		for i, neighbor := range city.directions {
			if neighbor != nil {
				// check that the neighbor city was not destroyed
				if city, ok := m.cities[neighbor.Name]; ok && !city.destroyed {
					fmt.Printf(" %s=%s", directionToStr(Direction(i)), neighbor.Name)
				}
			}
		}
		fmt.Println()
	}
}

// IsDestroyed tests whether the entire world is destroyed.
func (m *WorldMap) IsDestroyed() bool {
	for _, city := range m.cities {
		if !city.destroyed {
			// Found one city not yet destroyed, great!
			return false
		}
	}

	return true
}

// parseDirection returns a constant associated with its string representation.
func parseDirection(dir string) (Direction, error) {
	switch dir {
	case "north":
		return North, nil
	case "south":
		return South, nil
	case "east":
		return East, nil
	case "west":
		return West, nil
	default:
		return Invalid, fmt.Errorf("Invalid direction: %v", dir)
	}
}

// directionToStr returns the string representation for the given direction (North, East, ...)
func directionToStr(dir Direction) string {
	if dir > Invalid && dir < NDirections {
		return dirNames[dir]
	}
	return "Unknown"
}
