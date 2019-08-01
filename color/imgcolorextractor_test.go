package color

import (
	"log"
	"os"
	"testing"

	"github.com/deepakjacob/restyle/logger"
)

func TestColor(t *testing.T) {
	if err := logger.Init(-1, ""); err != nil {
		log.Fatal("logger initialization failed for tests")
	}
	path, _ := os.Getwd()
	path += "/test_data/saree2.jpg"
	img, err := os.Open(path)
	if err != nil {
		t.Error("error in opening file", err)
	}
	defer img.Close()
	err = Colors(img)
	if err != nil {
		t.Errorf("error in processing file %v", err)
	}
}
