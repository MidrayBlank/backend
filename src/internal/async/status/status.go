package status

type CompletionStatus struct {
	RequestId int
	Attempts  int
	Err       error
}

const (
	StatusQueued     = 1
	StatusInProgress = 2
	StatusFailed     = 3
	StatusSuccess    = 4
)
