package enums

type FlightClass int

const (
	Economy  FlightClass = 0
	Business FlightClass = 1
)

func FlightClassFromInt(value int) FlightClass {
	switch value {
	case 0:
		return Economy
	case 1:
		return Business
	default:
		return Economy
	}
}
