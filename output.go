package main

import (
	"fmt"
	"log"
	"strings"
)

// Generate a DROP TABLE statment
func (t *Table) DropSql() string {
	return fmt.Sprintf("DROP TABLE %s", t.NewName)
}

// Generate a CREATE statement for building the table
func (t *Table) CreateSql() string {
	cols := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		cols[i] = c.CreateSql()
	}
	return fmt.Sprintf("CREATE TABLE %s (\n   %s\n)", t.NewName, strings.Join(cols, ",\n   "))
}

// Generate a SELECT statement for the original MS Sql Server Table
func (t *Table) SelectMSSql() string {
	names := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		names[i] = c.OriginalName
	}
	nameList := strings.Join(names, ", ")
	return fmt.Sprintf("SELECT %s FROM %s", nameList, t.OriginalName)
}

// Generate the INSERT statement for Postgres
func (t *Table) InsertPsql() string {
	names := make([]string, len(t.Columns))
	place := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		names[i] = c.NewName
		place[i] = fmt.Sprintf(":%d", i+1)
	}
	nameList := strings.Join(names, ", ")
	placeList := strings.Join(place, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.NewName, nameList, placeList)
}

// Build the name/type pair for use in a create statement
func (c *Column) CreateSql() string {
	return fmt.Sprintf("%s %s", c.NewName, c.PostgresType())
}

// Convert MS SQL column to a Postgres type string
// Help: http://www.sqlines.com/sql-server-to-postgresql
func (c *Column) PostgresType() string {
	out := ""
	switch c.col.DATA_TYPE {
	case 4: //int
		return "INT"
	case -10: // ntext
		return "TEXT"
	case -7: // BIT
		return "BOOL"
	case 11: // smalldatetime
		return "TIMESTAMP(0)"
	case -1: // smalldatetime
		return "TEXT"
	case -9: // smalldatetime
		return fmt.Sprintf("VARCHAR(%v)", c.col.PRECISION)
	case 12: // varchar
		return fmt.Sprintf("VARCHAR(%v)", c.col.PRECISION)
	case 1: // char
		return fmt.Sprintf("CHAR(%v)", c.col.PRECISION)
	case 6: // float
		return "FLOAT"
	default:
		log.Fatalf("Dont know how to translate %d (%s)", c.col.DATA_TYPE, c.col.TYPE_NAME)
	}
	return out
}
