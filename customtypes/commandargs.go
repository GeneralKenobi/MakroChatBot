package customtypes

// CommandArgs is a struct for arguments that are passed to command handlers
type CommandArgs struct {

	// Name of the invoked command
	CommandName string

	// Name of the user that invoked the command
	Username string

	// ID of the user that invoked the command
	UserID string

	// Arguments passed by the user
	UserArgs []string
}
