package reactions

import (
	"errors"
	dg "github.com/bwmarrin/discordgo"
	ct "github.com/generalkenobi/makrochatbot/customtypes"
	"os"
)

// fileNames contains mapping of command name to corresponding name of file containing reaction image
var fileNames = map[string]string{
	"group1": "omegalul.png",
	"group2": "pogchamp.png",
}

// reactionImagesPath is the path to the directory containing reaction images used by commands
const reactionImagesPath = "resources/imagereactions/"

// ImageReaction returns a message containing some reaction image (depending on the invoked command)
func ImageReaction(args *ct.CommandArgs) (*dg.MessageSend, error) {

	// Variable for full path of the file - initialize it with directory name
	fullFilePath := reactionImagesPath

	// Try to get the filename from map
	if item, ok := fileNames[args.CommandName]; ok {
		// And add it to the file path
		fullFilePath += item
	} else {
		// If we couldn't obtain it, return an error - this shouldn't happen (every ImageReaction command should have an image assigned)
		return nil, errors.New("ImageReaction command was recognized and fired but no reaction image was found. Command: " + args.CommandName)
	}

	// Try to create file reader
	fileReader, err := os.Open(fullFilePath)

	// Try to open the file
	if err != nil {
		// In case of failure, return the error (add some information about command as well)
		return nil, errors.New("ImageReaction command couldn't open reaction image file. Command: " + args.CommandName +
			", Full file path: " + fullFilePath +
			", File open error: " + err.Error())
	}

	// Create file for the message
	file := &dg.File{
		Name:   fullFilePath,
		Reader: fileReader,
	}

	// Create the message and include the file in it
	message := &dg.MessageSend{
		Files: []*dg.File{
			file,
		},
	}

	return message, nil
}
