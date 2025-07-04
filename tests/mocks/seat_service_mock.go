package mock_repositories

import (
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockSeatService struct {
	mock.Mock
}

var _ interfaces.SeatService = (*MockSeatService)(nil)

func (m *MockSeatService) GetByFlightCode(flightCode string) ([]models.Seat, error) {
	args := m.Called(flightCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Seat), args.Error(1)
}
