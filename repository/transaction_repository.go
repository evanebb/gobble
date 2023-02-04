package repository

// For later use, I don't need this right now :)

type TransactionRepository interface {
	BeginTransaction() error
	CommitTransaction() error
	AbortTransaction() error
}
