package delivery

type Service struct {}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Calculate(orderAmount int64) int64 {
	if orderAmount >= 200000 {
		return 0
	}

	return 19900
}
