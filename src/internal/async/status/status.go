package status

type CompletionStatus struct {
	RequestId int
	Attempts  int
	Err       error
}
