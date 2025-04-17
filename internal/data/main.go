package data

type MainQ interface {
	New() MainQ
	Transfer() TransferQ
}
