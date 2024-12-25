package migrations

import (
	"context"
	"log"
)

<<<<<<< HEAD
func createTables() error {
	for _, query := range createTablesQueries {
		_, err := pool.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Error doing table initialisations due to %s: %v", query, err)
			return err
		}
	}

	log.Println("All tables which schemas are specified in db/migrations/creation_queries.go exists. If wasn't there before, they have been created")
=======

func createTables() error {
	var err error

	for _, query := range createTablesQueries {
		err = createTable(query)
		if err != nil {
			log.Fatalf("Error creating table with %s: %v", query, err)
			return err
		}
	}

	log.Println("All tables which schemas are specified in db/migrations/creation_queries.go exists. If wasn't there before, they have been created")
	return nil
}


func createTable(query string) error{
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Error running %s", query)
		return err
	}
>>>>>>> 0c747eda8b993aa85bea67e3eacdcb732218ff0c
	return nil
}

// func CreateTables(){
// 	conn, err := Pool.Acquire(context.Background());

// 	conn.Exec(
// 		context.Background(),
// 		generateCreateTableQuery("users", Models.User{}),
// 	);
// 	if err != nil {
// 		log.Fatal("Error while creating 'users' table");
// 	}

// }

// func goTypetoSQL(goType reflect.Type) string {
// 	switch goType.Kind() {
// 	case reflect.String:
// 		return "VARCHAR(255)"
// 	case reflect.Int:
// 		return "INTEGER"
// 	case reflect.Float64:
// 		return "DOUBLE PRECISION"
// 	case reflect.Bool:
// 		return "BOOLEAN"
// 	case reflect.Struct:
// 		if goType == reflect.TypeOf(time.Time{}) {
// 			return "TIMESTAMP"
// 		}
// 	}
// 	return "TEXT"
// }

// func generateCreateTableQuery(tableName string, model interface{}) string{
// 	var query string;
// 	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName)

// 	typ := reflect.TypeOf(model)

// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		columnName := field.Tag.Get("json")

// 		sqlType := goTypetoSQL(field.Type)
// 		query += fmt.Sprintf(" %s %s", columnName, sqlType)

// 		if field.Name == "ID" {
// 			query += " PRIMARY KEY"
// 		}

// 		if i < typ.NumField() - 1 {
// 			query += ","
// 		}
// 		query += "\n"
// 	}

// 	query += ");"
<<<<<<< HEAD
// 	return query;
=======
// 	return query;
>>>>>>> 0c747eda8b993aa85bea67e3eacdcb732218ff0c
