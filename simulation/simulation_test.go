package simulation

import (
	"strings"
	"testing"

	"github.com/darkfronza/alien-invasion/worldmap"
	"github.com/stretchr/testify/require"
)

// A world with two cities cannot survive an attack by 4 aliens
func TestSimulationWorldFullyDestroyed(t *testing.T) {
	input := `
Foo	north=Bar
Bar south=Foo`

	wm, err := worldmap.Load(strings.NewReader(input))
	require.Nil(t, err)
	require.False(t, wm.IsDestroyed())

	sim, err := New(wm, 4)
	require.Nil(t, err)

	for sim.Step() {
	}

	require.True(t, wm.IsDestroyed())
}

// A 4 cities world must survive to an alien invasion with just 2 aliens.
func TestSimulationWorldResists(t *testing.T) {
	input := "Foo north=Bar west=Jamaica south=Peru"

	wm, err := worldmap.Load(strings.NewReader(input))
	require.Nil(t, err)
	require.False(t, wm.IsDestroyed())

	sim, err := New(wm, 2)
	require.Nil(t, err)

	for sim.Step() {
	}

	require.False(t, wm.IsDestroyed())
}
