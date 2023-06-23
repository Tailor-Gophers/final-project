package models

type Reservation struct {
	Model
	PassengerID   uint         `json:"passenger_id,omitempty"`
	FlightClassID uint         `json:"-"`
	OrderID       uint         `json:"-"`
	Price         uint         `json:"price,omitempty"`
	IsCancelled   bool         `json:"is_cancelled,omitempty"`
	FlightClass   *FlightClass `json:"flight_class,omitempty"`
	Passenger     *Passenger   `json:"passenger,omitempty"`
}
