package drivestate

// DriveState struct
type DriveState struct {
	IsActive bool
	Dispo    string
}

type state map[string]*DriveState

var currentState = make(state)

// GetDriveState get the state of a drive
func GetDriveState(driveID string) *DriveState {
	return currentState[driveID]
}

// NewDriveState create a new store
func NewDriveState(driveID string, state *DriveState) {
	currentState[driveID] = state
}
