package roll

import (
	dgo "github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
)

// Default maximum value for random number geenration - i.e. result of random will be from 0 to 100
const defaultMaxRandomValue = 100

// Roll command
// Custom arguments:
// 1 - if convertible to int and positive it's treated as the boundary for random number generation
func Roll(session *dgo.Session, args []string) {

	if len(args) == 0 {
		// Log error - for some reason there isn't even the guaranteed userID argument
	}

	// Take the default max value
	max := defaultMaxRandomValue

	// If there is more than 1 argument
	if len(args) > 1 {
		// Try to convert it to an int, if conversion is successful and the number is positive use it as maximum instead
		if conversion, err := strconv.Atoi(args[1]); err == nil && conversion > 0 {
			max = conversion
		}
	}
}

// rollHelper returns a random integer number from 0 to max (excluding), assigns a reaction message that can be displayed to user.
// rand should be seeded when initializing.
// The random number is generated from 0 to max
func rollHelper(max int, userID string) []string {

	// Generate a random number (+1 to include the max value)
	randomNumber := rand.Intn(max + 1)

	rollMessage := userID + " rolled (0 - " + string(max) + "): " + string(randomNumber) + "!"

	var reaction string

	// Assign a reaction to the randomed number
	switch {

	case randomNumber < 10:
		{
			reaction = "I don't care what universe you're from, that's got to hurt!"
		}

	case randomNumber > 95:
		{
			reaction = "UNLIMITED POWER!"
		}

	case randomNumber > 75:
		{
			reaction = "A surprise, to be sure, but a welcome one"
		}

	default:
		{
			reaction = "In my experience there's no such thing as luck"
		}
	}

	return []string{rollMessage, reaction}
}