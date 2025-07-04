package mock_repositories

import (
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockSeatRepository struct {
	mock.Mock
}

var _ interfaces.SeatRepository = (*MockSeatRepository)(nil)

func (m *MockSeatRepository) GetByFlightCode(flightCode string) ([]entities.SeatOptionEntity, error) {
	args := m.Called(flightCode)
	return args.Get(0).([]entities.SeatOptionEntity), nil
}
