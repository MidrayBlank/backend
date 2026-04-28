package abstract

type IRequestAttemptsMap interface {
	Put(requestID int) int
	Get(requestID int) int
	Has(requestID int) bool
}
