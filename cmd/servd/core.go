package servd

// Service compiles inquiry, payment and checkstatus as a service instance
type Service struct{}

// New returns an instantiated Service
func New() *Service {
	return &Service{}
}
