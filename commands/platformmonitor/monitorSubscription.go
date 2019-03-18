package platformmonitor

// monitorSubscription is a struct containing information related to one subscriber - Discord user that will be notified when
// a person appears on the monitored platform.
type monitorSubscription struct {

	// SubscriberID is the DiscordID of subscriber (user that will be notified when a person is found on the platform)
	SubscriberID string

	// Name to search for on the platform
	SubscribedTo string
}
