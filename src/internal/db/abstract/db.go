package abstract

type IDBConnection interface {
	Get() any
	BeginTx() IDBConnection
	Commit() error
	Rollback()
}
