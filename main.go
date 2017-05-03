package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

type config struct {
	from   string
	to     string
	tables []string
	drop   bool
}

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

func ToColumn(col *MSSqlColumn) Column {
	return Column{
		OriginalName: col.COLUMN_NAME,
		NewName:      NameToPsql(col.COLUMN_NAME),
		col:          col,
	}
}
func main() {
	from, _, table_names := getArgs()

	db, err := sql.Open("mssql", from)
	if err != nil {
		log.Fatal(err)
	}

	tables := []Table{}

	for _, table := range table_names {
		cols := getColumns(table, db)
		tt := Table{
			OriginalName: table,
			NewName:      NameToPsql(table),
			Columns:      cols,
		}
		tables = append(tables, tt)

		log.Println(tt.DropSql())
		log.Println(tt.CreateSql())
		log.Println(tt.SelectMSSql())
		log.Println(tt.InsertPsql())
	}

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
