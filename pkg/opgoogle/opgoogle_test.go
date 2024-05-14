package opgoogle

import (
	"log"
	"testing"

	_ "embed"
)

var (
	//go:embed credentials.json
	credJSON []byte
	//go:embed token.json
	tokenJSON []byte
	//go:embed secret.json
	secretJSON []byte
)

func TestGetDoc(t *testing.T) {
	docID := "1HrwVZLRTGiK9ZsKmG8OqsBQdpAr1rCkE9zxOIIOvAgA"

	docSv, err := GetSheetService(credJSON, tokenJSON)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := docSv.Spreadsheets.Values.Get(docID, "sheet1!F2:F").Do()
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Values) == 0 {
		log.Println("No data found.")
	} else {
		log.Println("Name:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			if len(row) == 0 {
				continue
			}
			log.Printf("%s\n", row[0])
		}
	}
}
