package file_reader

import (
	"divya07514-hftest-golang/db"
	"divya07514-hftest-golang/logger"
	"divya07514-hftest-golang/model"
	"encoding/json"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

const (
	batchSize  = 10000 // Size of each batch
	numWorkers = 50    // Number of concurrent workers (goroutines)
)

type FileReader interface {
	ProcessFile(filename string) error
}

type JsonFileReader struct {
	Wg sync.WaitGroup
	Db db.DataSource
}

// ProcessFile Read json from a given file and insert the records in the db
func (r *JsonFileReader) ProcessFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		logger.ErrorLogger.Error().Msgf("error opening file: %v", err.Error())
	}
	defer file.Close()
	data := make(chan []model.RecipeData, numWorkers)
	r.Wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go r.processRecords(data)
	}
	r.readData(data, file)
	r.Wg.Wait()
	logger.InfoLogger.Info().Msgf("input file processed")
	return nil
}

// read data and insert into db in batches
func (r *JsonFileReader) readData(data chan []model.RecipeData, file *os.File) {
	go func() {
		defer close(data)
		decoder := json.NewDecoder(file)
		var batch []model.RecipeData
		if _, err := decoder.Token(); err != nil {
			logger.ErrorLogger.Error().Msgf("Error reading array start: %v", err.Error())
			return
		}
		for {
			var record model.RecipeData
			if err := decoder.Decode(&record); err != nil {
				if err.Error() != "EOF" {
					logger.ErrorLogger.Error().Msgf("error decoding json: %v", err.Error())
				}
				break
			}
			batch = append(batch, record)
			if len(batch) >= batchSize {
				data <- batch
				batch = nil
			}
		}
		if _, err := decoder.Token(); err != nil {
			logger.ErrorLogger.Error().Msgf("Error reading array end: %v", err.Error())
		}
	}()
}

func (r *JsonFileReader) processRecords(dataChannel <-chan []model.RecipeData) {
	defer r.Wg.Done()
	log.Info().Msgf("inserting data into db")
	r.Db.InsertBatch(dataChannel)
}
