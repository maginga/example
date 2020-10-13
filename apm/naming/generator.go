package naming

import (
	"fmt"
	"math/rand"
)

var (
	// NOUNS ...
	NOUNS = []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning",
		"snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest",
		"hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly",
		"bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass",
		"haze", "mountain", "night", "pond", "darkness", "snowflake", "silence",
		"sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower",
		"wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "canoe", "kayak",
		"frost", "voice", "paper", "frog", "smoke", "star", "dog", "space", "mario", "nike",
		"venus", "jupiter", "earth", "mars", "saturn", "uranus", "neptune", "pluto", "mercury"}
)

// Generator ...
type Generator interface {
	Generate() string
}

// NameGenerator ...
type NameGenerator struct {
	random *rand.Rand
}

// Generate ...
func (rn *NameGenerator) Generate() string {
	randomNoun := NOUNS[rn.random.Intn(len(NOUNS))]
	randomName := fmt.Sprintf("%v", randomNoun)
	return randomName
}

// NewNameGenerator ...
func NewNameGenerator(seed int64) Generator {
	nameGenerator := &NameGenerator{
		random: rand.New(rand.New(rand.NewSource(99))),
	}
	nameGenerator.random.Seed(seed)

	return nameGenerator
}
