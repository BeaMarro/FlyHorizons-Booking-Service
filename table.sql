-- Booking Table
CREATE TABLE Booking (
    ID INT PRIMARY KEY IDENTITY(1, 1) NOT NULL,
    UserID INT NOT NULL,
    FlightCode NVARCHAR(10) NOT NULL,
    FlightClass INT NOT NULL,
    Luggage NVARCHAR(150) NOT NULL,
    Status NVARCHAR(10) NULL,
    CreatedAt DATETIME NOT NULL
)

-- Passenger Table
CREATE TABLE Passenger (
    ID INT PRIMARY KEY IDENTITY(1, 1) NOT NULL,
    BookingID INT NOT NULL,
    FullName NVARCHAR(100) NOT NULL,
    DateOfBirth DATETIME NOT NULL,
    PassportNumber NVARCHAR(50) NOT NULL,
    Email NVARCHAR(255) NOT NULL,
    FOREIGN KEY (BookingID) REFERENCES Booking(ID)
)

-- Seat Table
-- Seats that can be selected for the Booking
CREATE TABLE Seat (
    ID INT PRIMARY KEY IDENTITY(1, 1) NOT NULL,
    BookingID INT NOT NULL,
    Row INT NOT NULL,
    [Column] CHAR(1) NOT NULL,  
    FOREIGN KEY (BookingID) REFERENCES Booking(ID)
)

-- Seat Option
-- Seat options to choose from for the Booking
CREATE TABLE SeatOption (
    ID INT PRIMARY KEY IDENTITY(1, 1) NOT NULL,
    Row INT NOT NULL,
    [Column] CHAR(1) NOT NULL,
    UNIQUE (Row, [Column]) -- Ensure no duplicate seats
)