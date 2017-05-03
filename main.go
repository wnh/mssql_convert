package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

type Table struct {
	OriginalName string
	NewName      string
	Columns      []Column
	PrimaryKey   []*Column
	ForiegnKeys  []ForeignKey
}

type Column struct {
	OriginalName string
	NewName      string
	col          *MSSqlColumn
}

type ForeignKey struct {
}

func main() {
	from, to, table_names := getArgs()

	msDB := ConnectAndTest("mssql", from)
	psqlDB := ConnectAndTest("postgres", to)

	tables := []Table{}

	for _, table := range table_names {
		cols := getColumns(table, msDB)
		tt := Table{
			OriginalName: table,
			NewName:      NameToPsql(table),
			Columns:      cols,
		}
		tables = append(tables, tt)

		log.Println("Dropping  ", tt.NewName)
		if _, err := psqlDB.Exec(tt.DropSql()); err != nil {
			log.Fatal(err)
		}

		log.Println("Createing ", tt.NewName)
		if _, err := psqlDB.Exec(tt.CreateSql()); err != nil {
			log.Fatal(err)
		}

		log.Println("Copying   ", tt.NewName)
		if err := CopyTable(msDB, psqlDB, tt); err != nil {
			log.Fatal(err)
		}

	}

}

func CopyTable(from, to *sql.DB, table Table) error {
	tx, err := to.Begin()
	if err != nil {
		return err
	}

	rows, err := from.Query(table.SelectMSSql())
	if err != nil {
		log.Fatal(err)
	}

	rr := make([]interface{}, len(table.Columns))
	ra := make([]interface{}, len(table.Columns))
	for i, _ := range ra {
		ra[i] = &rr[i]
	}

	count := 0
	for rows.Next() {
		count++
		rows.Scan(ra...)
		_, err := tx.Exec(table.InsertPsql(), rr...)
		if err != nil {
			log.Println(err)
		}
		if count%100 == 0 {
			log.Print(count)
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("Copied %d rows", count)
	rows.Close()
	return nil
}

func ConnectAndTest(driverName, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("Error opening %s database: %s", driverName, err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging %s database: %s", driverName, err)
	}
	return db
}

const usage = `mssql_migrate <from> <to> <table> [table ...]`

func getArgs() (from, to string, tables []string) {
	if len(os.Args) < 4 {
		log.Fatal("Usage: ", usage)
	}
	return os.Args[1], os.Args[2], os.Args[3:]
}

func getColumns(table string, db *sql.DB) []Column {
	rows, err := db.Query(fmt.Sprintf("sp_columns %s", table))
	if err != nil {
		log.Fatal(err)
	}

	out := []Column{}
	defer rows.Close()
	for rows.Next() {
		col := MSSqlColumn{}
		col.Scan(rows)
		cc := ToColumn(&col)
		out = append(out, cc)
	}
	return out
}

// Converts names from intercaps to snake case preserving
func NameToPsql(in string) string {
	x := []string{}
	acc := ""
	for i, r := range in {
		if i != 0 && (unicode.IsUpper(r) || unicode.IsDigit(r)) {
			x = append(x, acc)
			acc = ""
		}
		acc += string(unicode.ToLower(r))
	}
	x = append(x, acc)
	out := ""
	lastSmall := false
	for i, part := range x {
		imSmall := len(part) == 1
		if !(imSmall && lastSmall) {
			if i != 0 {
				out += "_"
			}
		}
		lastSmall = imSmall
		out += part
	}
	return out
}
