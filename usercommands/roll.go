package usercommands

import "math/rand"

// Roll returns a random integer number from 0 to 99, assigns a reaction message that can be displayed to user
func Roll() (int, string) {

	randomNumber := rand.Intn(100)
	var reaction string

	switch {

	case randomNumber < 25:
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

	return randomNumber, reaction
}
