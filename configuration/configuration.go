package configuration

import (
	"encoding/json"
	"os"

	ct "github.com/generalkenobi/makrochatbot/customtypes"
)

// Name of the config file
var configFileName = "config.json"

// GetConfig attempts to load the config file from disk and return a valid Config struct
// If config file cannot be opened, empty Config struct and the file open error will be returned
// If config file cannot be decoded, empty Config struct and the decoder error will be returned
func GetConfig() (ct.Config, error) {

	// Try to open the config file
	file, fileOpenError := os.Open(configFileName)
	defer file.Close()

	// If unsuccessful, return an empty config struct and the error obtained when opening file
	if fileOpenError != nil {
		return ct.Config{}, fileOpenError
	}

	// Create a decoder for json format
	decoder := json.NewDecoder(file)

	// And a variable to store the conversion result
	configuration := ct.Config{}

	// Try to decode the file
	decodeError := decoder.Decode(&configuration)

	// If decoding succeeded, return the struct
	if decodeError == nil {
		return configuration, nil
	}

	// Otherwise return an empty struct and the error
	return ct.Config{}, decodeError
}
