package mock_repositories

import (
	"flyhorizons-bookingservice/models/enums"
	entities "flyhorizons-bookingservice/repositories/entity"
	"flyhorizons-bookingservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockBookingRepository struct {
	mock.Mock
}

var _ interfaces.BookingRepository = (*MockBookingRepository)(nil)

func (m *MockBookingRepository) GetAll() []entities.BookingEntity {
	args := m.Called()
	return args.Get(0).([]entities.BookingEntity)
}

func (m *MockBookingRepository) GetByID(id int) entities.BookingEntity {
	args := m.Called(id)
	return args.Get(0).(entities.BookingEntity)
}

func (m *MockBookingRepository) GetByUserID(userID int) []entities.BookingEntity {
	args := m.Called(userID)
	return args.Get(0).([]entities.BookingEntity)
}

func (m *MockBookingRepository) Create(booking entities.BookingEntity) *entities.BookingEntity {
	args := m.Called(booking)
	return args.Get(0).(*entities.BookingEntity)
}

func (m *MockBookingRepository) DeleteByBookingID(ID int) bool {
	args := m.Called(ID)
	return args.Bool(0)
}

func (m *MockBookingRepository) UpdateStatus(bookingID int, status enums.Status) {
	m.Called(bookingID, status)
}

func (m *MockBookingRepository) Update(booking entities.BookingEntity) entities.BookingEntity {
	args := m.Called(booking)
	return args.Get(0).(entities.BookingEntity)
}
