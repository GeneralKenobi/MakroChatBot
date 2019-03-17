package platformmonitor

import (
	ct "github.com/generalkenobi/makrochatbot/customtypes"
	"github.com/generalkenobi/makrochatbot/logger"

	dg "github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"strings"
)

var monitoredURLs = make(map[string][]monitorTarget)

// Monitor monitors
func Monitor(args *ct.CommandArgs) (*dg.MessageSend, error) {

	resp, err := http.Get("https://platforma.polsl.pl/rau2/")

	if err != nil {
		return nil, nil
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, nil
	}

	htmltext := string(html)

	message := "Result : "

	if strings.Contains(htmltext, "Żogała") {
		message += "found"
	} else {
		message += "not found"
	}

	logger.Log(message)

	return nil, nil
}

func addListener(userID, listenTo, url string) {

	if item, ok := monitoredURLs[url]; ok {
		monitoredURLs[url] = append(monitoredURLs[url], monitorTarget{})
	} else {
		monitoredURLs[url] = []monitorTarget{monitorTarget{
			TargetName: listenTo,
			Listeners: map[string]struct{}{
				userID: {},
			},
		}}
	}

}
