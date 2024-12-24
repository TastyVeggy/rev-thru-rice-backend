package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	conn *pgxpool.Conn
)

func CreateTables() {
	var err error
	conn, err = Pool.Acquire(context.Background())
	if err != nil {
		log.Println("Error while acquiring connection from the database pool!")
	}

	defer conn.Release()
}

func createTable(query string) {
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Error running %s", query)
	}

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
// 	return query;
// }
