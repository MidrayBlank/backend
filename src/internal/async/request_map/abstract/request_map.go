package abstract

type IRequestAttemptsMap interface {
	Put(requestID int)
	Delete(requestID int)
	Get(requestID int) int
	Has(requestID int) bool
}
