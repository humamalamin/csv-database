# CSV Processor Library

## Overview
The **CSV Processor Library** is a Go package designed to facilitate the reading and processing of CSV files and inserting data into a PostgreSQL database. It is structured as a reusable library, making it easy for developers to integrate into their projects.

## Features
- Read CSV files with customizable headers.
- Dispatch multiple workers for concurrent data processing.
- Insert batch data into PostgreSQL databases efficiently.
- Modular structure for maintainability and reusability.

## Project Structure
```plaintext
csv_processor/
├── main.go
├── go.mod
├── db/
│   ├── connection.go      # Database connection utilities
│   ├── repo.go          # Database insertion logic
├── csv/
│   ├── reader.go          # CSV reading utilities
├── worker/
│   ├── dispatcher.go      # Worker dispatcher logic
├── utils/
│   ├── helpers.go         # Utility functions
```

## Getting Started

### Prerequisites
- **Go**: Version 1.18 or newer.
- **PostgreSQL**: Ensure PostgreSQL is installed and running.
- **CSV File**: Prepare a CSV file with data to process.

### Installation
Clone this repository and navigate to the project directory:
```bash
git clone https://github.com/your-repo/csv_processor.git
cd csv_processor
```

Install dependencies:
```bash
go mod tidy
```

### Configuration
Update the database configuration in `main.go`:
```go
// Database Configuration
dbConfig := db.Config{
    Host:        "localhost",
    Port:        5432,
    User:        "postgres",
    Password:    "your_password",
    DBName:      "csv_database",
    MaxConns:    20,
    MaxIdleConns: 10,
}
```

Specify the CSV file path in `main.go`:
```go
csvFilePath := "path/to/your-file.csv"
```

### Running the Application
Run the application with the following command:
```bash
go run main.go
```

## Usage

### Library Integration
You can use the modules in your own project by importing the packages:

#### Example: Using the Database Module
```go
import (
    "csv_processor/db"
)

func main() {
    dbConfig := db.Config{
        Host:        "localhost",
        Port:        5432,
        User:        "postgres",
        Password:    "your_password",
        DBName:      "your_database",
        MaxConns:    20,
        MaxIdleConns: 10,
    }

    database, err := db.OpenConnection(dbConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
}
```

#### Example: Reading CSV Files
```go
import (
    "csv_processor/csv"
)

func main() {
    csvReader, file, err := csv.OpenCsvFile("path/to/your-file.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    headers, err := csv.ReadHeaders(csvReader)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("CSV Headers:", headers)
}
```

## Contributing
Contributions are welcome! Feel free to submit issues or pull requests to improve the library.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact
For questions or support, please contact [humamalamin13@gmail.com].
