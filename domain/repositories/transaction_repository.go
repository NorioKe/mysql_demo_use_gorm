package repositories

type TransactionManager interface {
	Transaction(func() error) error // 统一事务接口
}
