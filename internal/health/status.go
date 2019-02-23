package health

// Status defines the type for a health check status
type Status int

const (
	// Ok means everything is fine
	Ok Status = 0

	// Warning means it's working, but a condition exists that could lead to a failure
	Warning Status = 1

	// Failure means something went wrong
	Failure Status = 2
)

func (s Status) String() string {
	switch s {
	case Ok:
		return "Ok"
	case Warning:
		return "Warning"
	case Failure:
		return "Failure"
	default:
		return "Unknown"
	}
}
