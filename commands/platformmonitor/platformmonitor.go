package platformmonitor

import (
	"github.com/generalkenobi/makrochatbot/communication"
	ct "github.com/generalkenobi/makrochatbot/customtypes"
	"github.com/generalkenobi/makrochatbot/logger"

	"errors"
	dg "github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MinimumSleepSeconds is the minimum number of seconds that have to pass between two subsequent platform checks.
// It is not possible to configure platform monitor to check with greater freqency.
const MinimumSleepSeconds = 10

// monitoredURLs holds all subscribers assigned to the monitored URL
// Key is the monitored URL.
// Value is a map used as a set (value is an empty struct) which holds monitorSubscriptions - structs used to hold subscribers
// and their targets.
var monitoredURLs = make(map[string]map[monitorSubscription]struct{})

// monitoredURLsMutex is a mutex used to take ownership of monitoredURLs
var monitoredURLsMutex sync.Mutex

// whitelistedURLs contains whitelisted urls and their aliases - those are only urls that may be subscribed to.
var whitelistedURLs = map[string]string{
	"rau2": "https://platforma.polsl.pl/rau2",
	"rau3": "https://platforma.polsl.pl/rau3",
}

// CreateMonitorSubscription adds a new subscription to the monitor service
// Arguments:
// 1. Name to subscribe to
// 2. Platform alias (e.g. "rau1", "rau2")
func CreateMonitorSubscription(args *ct.CommandArgs) (*dg.MessageSend, error) {

	// Message that will be sent to the user as feedback at the end of the call
	message := &dg.MessageSend{}
	defer communication.SendToUser(args.UserID, message)

	if len(args.UserArgs) < 2 {
		// Can't subscribe - not enough arguments
		return nil, errors.New("CreateMonitorSubscription: Not enough arguments: need 2, received " + strconv.Itoa(len(args.UserArgs)))
	}

	listenToName, urlAlias := args.UserArgs[0], args.UserArgs[1]

	// Try to get the url out of whitelist
	url, ok := whitelistedURLs[urlAlias]

	if !ok {
		// TODO: Notify user that the specified platform is incorrect
		message.Content = "The specified platform is incorrect"
		return nil, errors.New("CreateMonitorSubscription: Incorrect platform specified: " + urlAlias)
	}

	// If everything was correct, add a new subscription. Use its return value as feedback for user
	message.Content = addSubscriber(args.UserID, listenToName, url)

	return nil, nil
}

// RemoveAllSubscriptions removes all subscriptions from the user that invoked the command
func RemoveAllSubscriptions(args *ct.CommandArgs) (*dg.MessageSend, error) {

	// Remove the subscriptions
	removedSubscriptions := removeSubscriber(args.UserID)

	var messageContent string

	if len(removedSubscriptions) > 0 {

		// Add a header to the slice of removed subscriptions
		removedSubscriptions = append([]string{"Removed:"}, removedSubscriptions...)

		// And join it using newline as separator.
		// It will look like that:
		// Removed:
		// https://platform.polsl.pl/rau1 : Smiths
		// https://platform.polsl.pl/rau2 : Thompson
		messageContent = strings.Join(removedSubscriptions, "\n")
	} else {
		messageContent = "You weren't subscribed to anyone"
	}

	// Create a message with the string message
	message := &dg.MessageSend{
		Content: messageContent,
	}

	// And send it to the user
	communication.SendToUser(args.UserID, message)

	return nil, nil
}

// Start starts a new monitoring routine that will check all registered subscriptions, sleep for the provided number of seconds and
// then repeating the process indefinitely. The minimum number of seconds is MinimumSleepSeconds.
// Return value is a channel which is used to stop the created monitoring routine - it will be stopped when value 1 is passed to
// the channel or the channel is closed.
func Start(seconds int) chan int {

	// Create a channel that will be used to send stop signal
	stopChannel := make(chan int)

	// Start the monitoring routine using the provided number of seconds and the created stop channel
	go monitorRoutine(seconds, stopChannel)

	return stopChannel
}

// monitorRoutine fetches all urls and checks whether people that are monitored are present on the platform.
// If they are then subscribers are notified.
// Then it sleeps for the requested number of seconds (for safety reasons seconds cannot be smaller than MinimumSleepSeconds, if it is
// then it the minimum value of 10 will be used instead).
// Routine will stop if the channel is closed or value 1 is sent to the channel
func monitorRoutine(seconds int, stopChannel chan int) {

	// Make sure that there are at least 10 seconds between subsequent platform checks
	if seconds < MinimumSleepSeconds {
		seconds = MinimumSleepSeconds
	}

	// Keep going indefinitely
	for {

		// First check if there is any data sent to us
		select {
		case value, ok := <-stopChannel:
			{
				// Return if the channel is closed or the error code was passed - condition of closing
				if !ok || value == 1 {
					return
				}
			}

		// Otherwise run the platform check
		default:
			{
				runPlatfromCheck()
			}
		}

		// Finally go to sleep for the provided number of seconds
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

// runPlatformCheck checks all registered platforms and subscriptions to those platforms and notifies all subscribers about eventual
// matches.
func runPlatfromCheck() {

	logger.Log("Running platform check")

	// Take ownership of the mutex in order to work with monitoredURLs - we don't want it to be modified in the process
	monitoredURLsMutex.Lock()
	defer monitoredURLsMutex.Unlock()

	// Contains users that should be notified - key is userID and value is the message to send to him
	toNotify := make(map[string]string)

	for url, subscriptions := range monitoredURLs {

		// Use helper to download the html
		htmlText, err := downloadURLAsHTML(url)

		// If there was an error, log it and continue to another iteration
		if err != nil {
			logger.LogError(err)
			continue
		}

		// For each subscription found in the html
		for _, subscription := range findNamesInHTML(htmlText, subscriptions) {

			// If the subsriber is not yet in the toNotify map
			if _, ok := toNotify[subscription.SubscriberID]; !ok {

				// Add a header for him
				toNotify[subscription.SubscriberID] = "Platform notification:"
			}

			// The if above made sure that the header is present in the map, add to it newline for the triggered subscription
			toNotify[subscription.SubscriberID] += "\n" + subscription.SubscribedTo + " appeared on " + url[strings.LastIndex(url, "/"):]
		}
	}

	// Finally send notifications to users
	sendToUsers(toNotify)
}

// sendToUsers sends messages to a group of users.
// Each key in the provided map is a userID. Each value corresponding to it is the content of the message that should be sent to that user.
func sendToUsers(messages map[string]string) {

	// For each user to notify
	for userID, message := range messages {

		// Create a message send basing on the message text
		messageSend := &dg.MessageSend{
			Content: message,
		}

		// And send it
		communication.SendToUser(userID, messageSend)
	}
}

// findNamesInHTML searches through the htmlText for every SubscribedTo name in every monitorSubscription in subscriptions.
// Returns all subscriptions whose SubscribedTo was found in htmlText.
// htmlText is the html (converted to string) to search through.
// subscriptions are all subscriptions to check
func findNamesInHTML(htmlText string, subscriptions map[monitorSubscription]struct{}) []monitorSubscription {

	// Turn the text into lower-case, otherwise trivial mismatches would happen (such as:"Smith" and "smith" wouldn't match)
	htmlText = strings.ToLower(htmlText)

	// Slice for matched subscriptions
	var matchedSubscriptions []monitorSubscription

	// Cache containing already found names
	foundNamesCache := make(map[string]struct{})

	// Cache with names that were already checked and were not found
	notFoundNamesCache := make(map[string]struct{})

	// For each subscription for this url
	for subscription := range subscriptions {

		// Check if the name is in the not found cache
		if _, ok := notFoundNamesCache[subscription.SubscribedTo]; ok {
			// If so, continue because it's not present
			continue
		}

		// Check if the subscribed to is already in cache, if not check if he's in the hmtl text
		if _, ok := foundNamesCache[subscription.SubscribedTo]; ok || strings.Contains(htmlText, subscription.SubscribedTo) {
			// Add the found name to cache - it doesn't matter if it's overwritten, it just has to be there
			foundNamesCache[subscription.SubscribedTo] = struct{}{}

			// Add the matched subscription to result slice
			matchedSubscriptions = append(matchedSubscriptions, subscription)

		} else {
			// Else add the name to not found cache - if we got here it means it wasn't in that cache yet and it wasn't found either
			notFoundNamesCache[subscription.SubscribedTo] = struct{}{}
		}
	}

	return matchedSubscriptions
}

// downloadURLAsHTML tries to download and extract html from the given url.
// Returns the document (as string) and nil error if everything went well or empty string and an error indicating what went wrong
func downloadURLAsHTML(url string) (string, error) {

	// Try to get the url
	resp, err := http.Get(url)

	// Return error if it wasn't possible to download it
	if err != nil {
		return "", errors.New("Couldn't download " + url + ". Error: " + err.Error())
	}

	// Later close the body
	defer resp.Body.Close()

	// Try to read the whole body
	html, err := ioutil.ReadAll(resp.Body)

	// If there was an error
	if err != nil {
		// Return the error
		return "", errors.New("Couldn't read http response body, from " + url + ". Error: " + err.Error())
	}

	// Convert the result to string and return it
	return string(html), nil
}

// addSubscriber adds a new subscriber to the specified url
// Returns that can be sent to user as feedback (Information about success, failure, etc.).
func addSubscriber(userID, subscribeTo, url string) string {

	// Take ownership of the mutex in order to work with monitoredURLs
	monitoredURLsMutex.Lock()
	defer monitoredURLsMutex.Unlock()

	// Create a subscription containing the required information
	subscription := monitorSubscription{
		SubscriberID: userID,
		SubscribedTo: strings.ToLower(subscribeTo),
	}

	// First, check if the url exists
	if _, ok := monitoredURLs[url]; !ok {
		// If not, add it to the map
		monitoredURLs[url] = make(map[monitorSubscription]struct{})
	}

	// Check if such subscription is already present (url is guaranteed to be in the map)
	if _, ok := monitoredURLs[url][subscription]; ok {
		// If so, notify the user that he's already subscribed to this particular person on this particular platform
		return "You're already subscribed to this name on this url"
	}

	// If the subscription is new, add it to the map
	monitoredURLs[url][subscription] = struct{}{}
	return "Subscribed successfully"
}

// removeSubscriber removes all subscriptions assigned to the given userID.
// Returns all removed subscriptions as a slice. Each entry has a form "url : subscribed_to".
func removeSubscriber(userID string) []string {

	// Take ownership of the mutex in order to work with monitoredURLs
	monitoredURLsMutex.Lock()
	defer monitoredURLsMutex.Unlock()

	// Slice for all removed subscriptions - it will be used to inform the user about all removed subscriptions
	removedSubscriptions := []string{}

	// For each registered url
	for url, urlSubscriptions := range monitoredURLs {

		// For each subscription in the url
		for subscription := range monitoredURLs[url] {

			// If its userID matches the searched one
			if subscription.SubscriberID == userID {

				// Add the pair of url and subscribed to name to the designated slice
				removedSubscriptions = append(removedSubscriptions, url[strings.LastIndex(url, "/"):]+" : "+subscription.SubscribedTo)

				// And delete subscription key from the collection
				delete(urlSubscriptions, subscription)
			}

			// If there are no more subscriptions for this url, remove it from the collection
			if len(urlSubscriptions) == 0 {
				delete(monitoredURLs, url)
			}
		}
	}

	return removedSubscriptions
}

// removeSubscriberFrom removes the specified subscription (userID & subscribedTo pair) from the given url.
// Notifies user (given by userID) about the result.
// If the subscription wasn't found then nothing will happen.
func removeSubscriberFrom(userID, subscribedTo, url string) {

	// Take ownership of the mutex in order to work with monitoredURLs
	monitoredURLsMutex.Lock()
	defer monitoredURLsMutex.Unlock()

	// Check if the url is present in the collection
	if _, ok := monitoredURLs[url]; !ok {
		// Nothing to do - user was not subscribed to anyone on that URL
		// TODO: Notify user that he was not subscribed
		return
	}

	// Create a subscription containing the required information
	subscription := monitorSubscription{
		SubscriberID: userID,
		SubscribedTo: subscribedTo,
	}

	// Check if the subscription is present
	if _, ok := monitoredURLs[url][subscription]; ok {
		// Subscription found: remove it
		delete(monitoredURLs[url], subscription)

		// If it was the last subscription to that url then remove the url as well
		if len(monitoredURLs[url]) == 0 {
			delete(monitoredURLs, url)
		}

		// TODO: Notif the user that he unsubscribed successfully
	} else {
		// TODO: Notif the user that the requested subscription wasn't found
	}
}
