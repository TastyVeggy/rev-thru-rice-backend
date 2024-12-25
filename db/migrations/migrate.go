package migrations

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func InitDB(db_pool *pgxpool.Pool, wantToCreateTables bool, seedDataDir string) error {
	var err error
	pool = db_pool
	if wantToCreateTables {
		err = createTables()
		if err != nil {
			return err
		}
	}

	if seedDataDir != "" {
		err = loadSeedData(seedDataDir)
		if err != nil {
			return err
		}
	}
	return nil
}
