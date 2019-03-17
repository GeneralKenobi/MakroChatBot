package reactions

import (
	ct "../../customtypes"
	dg "github.com/bwmarrin/discordgo"
	"os"
)

// ImageReaction returns a message containing some reaction image (depending on the invoked command)
func ImageReaction(args *ct.CommandArgs) (*dg.MessageSend, error) {

	//reader, _ := os.Open("resources/imagereactions/omegalul.jpeg")
	//
	//image := dg.File{
	//	Name:        "OMEGALUL",
	//	ContentType: "image",
	//	Reader:      reader,
	//}

	filename := "resources/imagereactions/omegalul.jpeg"

	f, _ := os.Open(filename)

	embed := dg.MessageEmbed{
		Image: &dg.MessageEmbedImage{
			URL: "attachment://" + filename,
		},
	}

	message := dg.MessageSend{
		Content: "xd",
		Embed:   &embed,
		Files: []*dg.File{
			&dg.File{
				Name:   filename,
				Reader: f,
			},
		},
		//Files:   []*dg.File{&image},
	}

	return &message, nil
}
