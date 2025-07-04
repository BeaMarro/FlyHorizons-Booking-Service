package load_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/models/enums"
	load_test_utils "flyhorizons-bookingservice/tests/load/utils"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Before running the load tests
// Run the microservice at the same time

type BookingLoadTest struct {
	loadTestUtils load_test_utils.LoadTestUtils
}

func getBookingPayload() []byte {
	booking := models.Booking{
		ID:          0,
		UserID:      2,
		FlightCode:  "FR788",
		FlightClass: 0,
		Luggage: []enums.Luggage{
			enums.SmallBag,
			enums.CabinBag,
		},
		Seats: []models.Seat{
			{
				Row:       1,
				Column:    "A",
				Available: true,
			},
			{
				Row:       1,
				Column:    "B",
				Available: true,
			},
		},
		Passengers: []models.Passenger{
			{
				ID:             1,
				FullName:       "John Doe",
				Email:          "john@doe.nl",
				DateOfBirth:    time.Date(1985, 7, 9, 1, 0, 0, 0, time.UTC),
				PassportNumber: "1234",
			},
			{
				ID:             2,
				FullName:       "Jane Doe",
				Email:          "jane@doe.it",
				DateOfBirth:    time.Date(1986, 8, 8, 2, 30, 0, 0, time.UTC),
				PassportNumber: "4321",
			},
		},
	}

	payload, err := json.Marshal(booking)
	if err != nil {
		fmt.Println("Error marshaling booking entity:", err)
		return []byte{}
	}

	return payload
}

func createBooking(t *testing.T, rate vegeta.Rate, duration time.Duration, htmlReport string, title string) vegeta.Metrics {
	loadTest := BookingLoadTest{
		loadTestUtils: load_test_utils.LoadTestUtils{},
	}

	target := vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8083/bookings",
		Body:   getBookingPayload(),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(vegeta.NewStaticTargeter(target), rate, duration, "Load Test CreateBookings") {
		metrics.Add(res)
	}
	metrics.Close()

	// Log the metrics
	loadTest.loadTestUtils.LogMetrics(t, &metrics)
	loadTest.loadTestUtils.EvaluateMetricsSuccess(t, &metrics)
	loadTest.loadTestUtils.GenerateHTMLReport(t, &metrics, fmt.Sprintf("%s.html", htmlReport), fmt.Sprintf("Load Test: %s", title))

	return metrics
}

func createBookingsIncrementally(t *testing.T, rates []int, htmlReport string, title string) {
	loadTest := BookingLoadTest{
		loadTestUtils: load_test_utils.LoadTestUtils{},
	}

	target := vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8083/bookings",
		Body:   getBookingPayload(),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	// Define ramp-up steps
	for _, freq := range rates {
		rate := vegeta.Rate{Freq: freq, Per: time.Second}
		duration := 1 * time.Second

		for res := range attacker.Attack(vegeta.NewStaticTargeter(target), rate, duration, fmt.Sprintf("Ramp-Up %d rps", freq)) {
			metrics.Add(res)
		}
	}

	metrics.Close()

	loadTest.loadTestUtils.LogMetrics(t, &metrics)
	loadTest.loadTestUtils.EvaluateMetricsSuccess(t, &metrics)
	loadTest.loadTestUtils.GenerateHTMLReport(t, &metrics, fmt.Sprintf("%s.html", htmlReport), fmt.Sprintf("Load Test: %s", title))
}

// Load Test (short time): Create bookings at 10 requests per second for 1 second
func TestLoadCreateFewBookingsUsingShortTime(t *testing.T) {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "post_few_bookings_short_time"
	title := "Post bookings at 10 requests per second for 1 second"
	_ = createBooking(t, rate, duration, htmlReport, title)
}

// Load Test (short time): Create bookings at 200 requests per second for 1 second
func TestLoadCreateManyBookingsUsingShortTime(t *testing.T) {
	rate := vegeta.Rate{Freq: 200, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "post_many_bookings_short_time"
	title := "Post bookings at 200 requests per second for 1 second"
	_ = createBooking(t, rate, duration, htmlReport, title)
}

// Load Test (long time): Create bookings at 10 requests per second for 10 seconds
func TestLoadCreateBookingsUsingLongerTime(t *testing.T) {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 10 * time.Second
	htmlReport := "post_many_bookings_long_time"
	title := "Post bookings at 10 requests per second for 1 minute"
	_ = createBooking(t, rate, duration, htmlReport, title)
}

// Spike Test: Create bookings at 500 requests per second for 1 second
func TestLoadSpikeCreateBookings(t *testing.T) {
	rate := vegeta.Rate{Freq: 500, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "post_bookings_spike"
	title := "Spike test bookings at 500 requests per second for 1 second"
	_ = createBooking(t, rate, duration, htmlReport, title)
}

// Spike Test: Create bookings at 1000 requests per second for 1 second
func TestLoadSpikeExtremeCreateBookings(t *testing.T) {
	rate := vegeta.Rate{Freq: 5000, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "post_bookings_spike_extreme"
	title := "Spike test bookings at 5000 requests per second for 1 second"
	_ = createBooking(t, rate, duration, htmlReport, title)
}

// Stress Test: Create bookings by increasing requests from 10 bookings per second to 1000 bookings per second for 1 second
// func TestLoadStressCreateBookings(t *testing.T) { // TODO: Fix
// 	rates := []int{10, 100}
// 	htmlReport := "post_bookings_extreme_stress"
// 	title := "Stress test bookings from 10 to 1000 request per second for 1 second"
// 	createBookingsIncrementally(t, rates, htmlReport, title)
// }

// TODO: Add soak (endurance) testing
