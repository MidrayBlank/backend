package abstract

type IAIAPIRepository interface {
	Upsert(hash string) error

	GetRequestsCount(hash string) (int, error)

	ResetRequestsCount(hash string) error
}
