package abstract

type IDBTransaction interface {
	Commit() error
	Rollback()
}

type IDBConnection interface {
	Get() any
	BeginTx() IDBTransaction
}
