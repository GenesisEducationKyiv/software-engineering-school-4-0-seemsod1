package customer

type Service struct {
	CustomerRepo    Database
	SagaCoordinator *SagaCoordinator
}

func NewService(customer Database, saga *SagaCoordinator) *Service {
	return &Service{
		CustomerRepo:    customer,
		SagaCoordinator: saga,
	}
}
