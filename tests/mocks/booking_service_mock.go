package mock_repositories

import (
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockBookingService struct {
	mock.Mock
}

var _ interfaces.BookingService = (*MockBookingService)(nil)

func (m *MockBookingService) BookingExists(bookingID int) bool {
	args := m.Called(bookingID)
	return args.Bool(0)
}

func (m *MockBookingService) GetByUserID(userID int) []models.Booking {
	args := m.Called(userID)
	return args.Get(0).([]models.Booking)
}

func (m *MockBookingService) Create(booking models.Booking) (*models.Booking, error) {
	args := m.Called(booking)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockBookingService) DeleteByBookingID(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookingService) Update(booking models.Booking) (*models.Booking, error) {
	args := m.Called(booking)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}
