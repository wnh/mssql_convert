package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	flag "github.com/spf13/pflag"

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

type config struct {
	from   string
	to     string
	tables []string
	drop   bool
	print  bool
}

func main() {
	cfg := getArgs()

	msDB := ConnectAndTest("mssql", cfg.from)

	tables := []Table{}

	for _, table := range cfg.tables {
		cols := getColumns(table, msDB)
		tt := Table{
			OriginalName: table,
			NewName:      NameToPsql(table),
			Columns:      cols,
		}
		tables = append(tables, tt)

		tt.PrimaryKey = getPrimaryKeys(tt, msDB)

		//for _, p := range tt.PrimaryKey {
		//	log.Printf("%#v", p)
		//}
		//log.Println(tt.CreateSql())

		if cfg.print {
			fmt.Println(tt.CreateSql())
		} else {
			psqlDB := ConnectAndTest("postgres", cfg.to)

			if cfg.drop {
				log.Println("Dropping  ", tt.NewName)
				if _, err := psqlDB.Exec(tt.DropSql()); err != nil {
					log.Fatal(err)
				}
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

const usage = `mssql_migrate [--drop] [--print] <from> <to> <table> [table ...]`

func getArgs() config {
	cfg := config{}

	flag.BoolVar(&cfg.drop, "drop", false, "Drop tables before creating them")
	flag.BoolVar(&cfg.print, "print", false, "Dont execute, only print the creation SQL")
	flag.Usage = func() {
		fmt.Println("Usage: ", usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 3 {
		flag.Usage()
		os.Exit(1)
	}
	cfg.from = args[0]
	cfg.to = args[1]
	cfg.tables = args[2:]
	return cfg
}

func getPrimaryKeys(table Table, db *sql.DB) []*Column {
	rows, err := db.Query(fmt.Sprintf("sp_pkeys %s", table.OriginalName))
	if err != nil {
		log.Fatal(err)
	}

	out := []*Column{}
	defer rows.Close()
	for rows.Next() {
		pkey := MSSqlPKey{}
		pkey.Scan(rows)

		for _, c := range table.Columns {
			if c.OriginalName == pkey.COLUMN_NAME {
				out = append(out, &c)
				break
			}
		}
	}
	return out
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
