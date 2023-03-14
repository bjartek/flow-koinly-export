package koinly

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

//"2021-10-20 10:20 UTC",-310.18899997,"ID:7961","0xdf868d4de6d2e0ab","ce8c096607ac931b643e5c196d57c156d4be180ee5a8fe9bcc5ec0cd2e5371dc"
/*
   Date, Sent Amount, Sent Currency, Received Amount, Received Currency
   Fee Amount, Fee Currency, Net Worth Amount, Net Worth Currency, Label, Description, TxHash
*/

type DateTime struct {
	time.Time
}

// Convert the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.UTC().String(), nil
}

// You could also use the standard Stringer interface
func (date *DateTime) String() string {
	return date.String() // Redundant, just for example
}

type Event struct { // Our example struct, you can use "-" to ignore a field
	Date             DateTime `csv:"Date"`
	SentAmount       string   `csv:"Sent Amount,omitempty"`
	SentCurrency     string   `csv:"Sent Currency"`
	ReceivedAmount   string   `csv:"Received Amount"`
	ReceivedCurrency string   `csv:"Received Currency"`
	FeeAmount        string   `csv:"Fee Amount"`
	FeeCurrency      string   `csv:"Fee Currency"`
	Label            string   `csv:"Label"`
	Description      string   `csv:"Description"`
	TxHash           string   `csv:"TxHash"`
}

func Marshal(events []Event, fileName string) error {
	csvContent, err := gocsv.MarshalString(&events) // Get all clients as CSV string
	if err != nil {
		return err
	}
	//fmt.Println(csvContent)

	bytes := []byte(csvContent)
	err = os.WriteFile(fileName, bytes, 0644)
	return err
}
