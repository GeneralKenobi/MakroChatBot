package roll

import (
	dg "github.com/bwmarrin/discordgo"
	ct "github.com/generalkenobi/makrochatbot/customtypes"
	"math/rand"
	"strconv"
)

// Default maximum value for random number geenration - i.e. result of random will be from 0 to 100
const defaultMaxRandomValue = 100

// Roll command
// User arguments:
// 1 - if convertible to int and positive it's treated as the boundary for random number generation
func Roll(args *ct.CommandArgs) (*dg.MessageSend, error) {

	// Take the default max value
	max := defaultMaxRandomValue

	// If there is more than 1 argument
	if len(args.UserArgs) > 0 {
		// Try to convert it to an int, if conversion is successful and the number is positive use it as maximum instead
		if conversion, err := strconv.Atoi(args.UserArgs[0]); err == nil && conversion > 0 {
			max = conversion
		}
	}

	// Create a message to send, it has text content only. Text is generated by helper function
	message := &dg.MessageSend{
		Content: rollHelper(max, args.Username),
	}

	return message, nil
}

// rollHelper returns a random integer number from 0 to max (excluding), assigns a reaction message that can be displayed to user.
// rand should be seeded when initializing.
// The random number is generated from 0 to max
func rollHelper(max int, username string) string {

	// Generate a random number (+1 to include the max value)
	randomNumber := rand.Intn(max + 1)

	// Create message for result
	rollMessage := username + " rolled (0 - " + strconv.Itoa(max) + "): " + strconv.Itoa(randomNumber) + "!"

	var reaction string

	// Assign a reaction to the randomed number
	switch {

	case randomNumber < max/10:
		{
			reaction = "I don't care what universe you're from, that's got to hurt!"
		}

	case randomNumber > 19*max/20:
		{
			reaction = "UNLIMITED POWER!"
		}

	case randomNumber > 3*max/4:
		{
			reaction = "A surprise, to be sure, but a welcome one"
		}

	default:
		{
			reaction = "In my experience there's no such thing as luck"
		}
	}

	// Return a concatenation of both strings
	return rollMessage + "\n" + reaction
}
