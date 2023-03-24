package worldmap

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Table data where each element has a sample map data and an associated expected
	// number of cities to be loaded by the worldmap module.
	testCases := []struct {
		mapData          string
		expectedNoCities int
	}{
		{
			`Foo`, 1,
		},
		{
			`Foo north=Bar`, 2,
		},
		{
			`
Foo	north=Bar
Bar south=Foo
			`, 2,
		},
		{
			`
		Foo	north=Bar west=Moon
		Bar south=Foo
					`, 3,
		},
		{
			`
		Foo	north=Bar west=Moon, east=Silverado, south=Mars
					`, 5,
		},
		{
			`
		Foo
		Bar
		0Baz
					`, 3,
		},
		{
			``, 0,
		},
		{
			`
		Foo	north=Bar west=Moon, east=Silverado, south=Mars
		Bar south=Foo west=Venus
		Invalid souP=Oops
		ALmostValid south=Some north=Data westx=error
					`, 6,
		},
	}

	for _, td := range testCases {
		src := strings.NewReader(td.mapData)
		wm, err := Load(src)
		require.Nil(t, err)
		require.Equal(t, td.expectedNoCities, len(wm.GetAllCityNames()))
	}
}

func TestGetCity(t *testing.T) {
	input := `
Foo	north=Bar west=Moon east=Silverado south=Mars
Bar south=Foo west=Venus`

	wm, err := Load(strings.NewReader(input))
	require.Nil(t, err)
	require.NotNil(t, wm.GetCity("Foo"))
	require.NotNil(t, wm.GetCity("Bar"))
	require.NotNil(t, wm.GetCity("Moon"))
	require.NotNil(t, wm.GetCity("Silverado"))
	require.NotNil(t, wm.GetCity("Mars"))
	require.NotNil(t, wm.GetCity("Venus"))
	require.Nil(t, wm.GetCity("NonExistent"))
}

func TestDestroyCity(t *testing.T) {
	input := `
Foo	north=Bar west=Moon east=Silverado south=Mars
Bar south=Foo west=Venus`

	wm, err := Load(strings.NewReader(input))
	require.Nil(t, err)

	require.False(t, wm.IsDestroyed())

	cities := []string{"Foo", "Bar", "Moon", "Silverado", "Mars", "Venus"}
	total := len(cities)
	for _, city := range cities {
		require.NotNil(t, wm.GetCity(city))
		wm.DestroyCity(city)
		require.Nil(t, wm.GetCity(city))
		total -= 1
	}

	require.True(t, wm.IsDestroyed())
}
