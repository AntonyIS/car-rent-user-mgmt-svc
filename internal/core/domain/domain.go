package domain

type Payment struct {
	ID string
}

type DriverLicense struct {
	ID                  string
	DriverLicenseNumber string
	ExpiryDate          string
	LicenseClass        string
	LicensePhone        string
}

type Vehicle struct {
	ID                 string
	Make               string
	Model              string
	Year               int
	Body               string
	Color              string
	Transmission       string
	FuelType           string
	SeatCapacity       string
	LisencePlateNumber string
	VIN                string
	RentalPrice        float64
	Condition          string
	Image              string
	InsuranceCoverage  string
}

type RentalHistory struct {
	ID                 string
	RentalStatus       string
	InsuranceCover     string
	PaymentInformation []Payment
	Vehicle            Vehicle
}

type User struct {
	ID            string
	Firstname     string
	Lastname      string
	Email         string
	Phone         int
	DOB           string
	Address       string
	Payment       Payment
	Vehicle       Vehicle
	DriverLicense DriverLicense
	RentalHistory RentalHistory
}
