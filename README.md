<img src="https://static.vecteezy.com/system/resources/previews/002/914/945/original/beautiful-of-ufo-invasion-in-the-city-multicolored-background-free-vector.jpg"/>

# alien-invasion
An alien invasion simulator written in Go.

## Readme Table of Contents
<!-- TOC -->

- [Alien Invasion](#alien-invasion)
- [Readme Table of Contents](#readme-table-of-contents)
- [About the project](#about-the-project)
- [Installation and usage](#installation-and-usage)
  - [Installing directly with go](#installing-directly-with-go)
  - [Cloning and building from sources](#cloning-and-buulding-from-sources)
  - [Running tests](#running-tests)
  - [Running the map generator tool](#running-the-map-generator-tool)
  - [Usage](#usage)
- [Map file format](#map-file-format)
- [Assumptions made by the simulator](#assumptions-made-by-the-simulator)
- [Simulation step](#simulation-step)
- [Simulation outcomes](#simulation-outcomes)


<!-- /TOC -->

## About the project

The **alien-invasion** project attempts to simulate an alien invasion to our beloved world.

The simulator expects users to provide an input file representing a map with cities and directions to other cities, along with the number of potential alien invaders.

A tool for [generating random maps](#running-the-map-generator-tool) is included in the project for helping testing the simulator.

The simulator then simulates an alien invasion based on many [assumptions](#assumptions-made-by-the-simulator), which can randomly produce different outcomes.

## Installation and usage

Currently, the project has no pre-compiled binaries available, it assumes that the Go toolchain is installed on the target machine, please refer to this link [Go installation](https://go.dev/doc/install) and follow the instructions provided.

### Installing directly with go

Run the following command to get the latest `alien-invasion` binary installed on your pc:

```
$ go install github.com/darkfronza/alien-invasion@latest
```

### Cloning and building from sources
```
$ git clone github.com/darkfronza/alien-invasion@latest
$ cd alien-invasion
$ make
```

A binary named `alien-invasion` must be generated at the source folder.

### Running tests
```
$ make test
```

### Running the map generator tool

First, compile the *mapgen* tool.
```
$ make tools
```

Then, you can generate a sample *map.txt* worldmap file with approximately 100 cities by running the following command:
```
$ ./mapgen 100 > map.txt
```

This generated map file can then be used for experimenting with the simulator.

### Usage

The `alien-invasion` program expects two arguments, a map file and a number of aliens, in order:
* *map-file*: A file that lists cities and directions to other cities, please refer to [Map file format](#map-file-format).
* *n-aliens*: The number of aliens attempting to invade the world

Sample usage:

```
$ ./alien-invasion worldmap.txt 10
```

## Map file format

The map file has a simple file layout:

```
CityName east=CityToTheEast west=CityToTheWest south=CityToTheSouth north=CityToTheNorth
...
CityNameN [direction=AnotherCityN]
```

* Each line represents a city along with optional directions to another cities.
* All the fields are separated by a single space.
* The first field is a city name.
* The remainder fields are optional and have the following format: direction=cityName, where direction must be one of: north,south,east,west. They represent directions that lead to another cities from the original one.


## Assumptions made by the simulator

* The aliens are randomly distributed to the cities available.
* The aliens are identified by a string "Alien#ID" where ID is a number between 1 and *n-aliens* argument provided, inclusive.
* The number of aliens can exceed the number of cities, in this cased some cities will contain more than one alien during the simulation setup.
* All the aliens have a shared CPU clock provided by the mothership, that means all aliens move simultaneously during each [simulation step](#simulation-step).
* After each [simulation step](#simulation-step), if more than one alien ended up landing in the same city, all of them will be mutually destroyed, the city is also destroyed in the process and no alien can move to this city again during the simulation.
* Aliens that were dropped by the mothership at the same city during setup don't count for self/city destrucion during simulation steps, those that moved do count.
* The simulation stops on any of the following conditions:
  * All aliens were mutually destroyed
  * All aliens are stuck, which happens when all of them have moved at least 10,000 times.
 
 ## Simulation step
 
 The simulation revolves around the concept of a simulation step, during which the following events takes place:
 
 * All the aliens are scanned and for each of them check if the city they are currently in has connections to another cities (north, south, east, west).
 * If there is connection to another city ou cities, enumerate those not yet destroyed and chose one of them at random.
 * Move each aliens to the next random city and increment the number of steps taken by each of them.
 * If more than one alien moved to the same target city, mutually destroy all of them along with the city and connections to the city. The list of destroyed aliens is displayed in the screen during the step.
 * If all the aliens are stuck or all of them are destroyed, end the simulation.


## Simulation outcomes

When the simulation ends, the following are the list of possible outcomes produced by the simulator:

* All aliens were destroyed and all the cities were destroyed. In this case, goodbye world and nothing else to show about it
* All aliens were destroyed, some or none of the cities were destroyed. In this case, what remained from the world is displayed at the end of the simulation using the same map format for input.
* Some or no aliens were destroyed. In this case, the simulation ends when all of them get stuck, what remained from the world is displayed at the screen.


