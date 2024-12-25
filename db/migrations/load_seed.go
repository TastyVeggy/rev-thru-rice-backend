package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadSeedData(seedDataDir string) error {
	files, err := os.ReadDir(seedDataDir)
	if err != nil {
		log.Fatalf("Failed to read seed data directory: %v", err)
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			tableName := file.Name()[0 : len(file.Name())-5]
			filePath := filepath.Join(seedDataDir, file.Name())

			err := loadAndInsertData(filePath, tableName)

			if err != nil {
				log.Fatalf("Failed to insert data finish from %s: %v", file.Name(), err)
				return err
			} else {
				log.Printf("Successfully inserted data from %s", file.Name())
			}
		}
	}
	return nil
}

func loadAndInsertData(filePath string, tableName string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", filePath, err)
		return err
	}

	var data []map[string]any
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("Failed to get store data into data variable from json file %s: %v", filePath, err)
		return err
	}
	return insertData(tableName, data)

}

func insertData(table string, data []map[string]any) error {
	for _, record := range data {
		var columns []string
		var values []any
		for key, value := range record {
			columns = append(columns, key)
			values = append(values, value)
		}

		columnsStr := strings.Join(columns, ", ")
		placeholders := make([]string, len(columns))
		for i := range placeholders {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}

		placeholdersStr := strings.Join(placeholders, ", ")

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING", table, columnsStr, placeholdersStr)

		_, err := pool.Exec(context.Background(), query, values...)

		if err != nil {
			log.Fatalf("failed to execute insert query: %v", err)
			return err
		}

	}
	return nil
}
