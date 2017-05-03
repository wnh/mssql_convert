package main

import (
	"database/sql"
)

type MSSqlColumn struct {
	TABLE_QUALIFIER   string //Object qualifier name. This field can be NULL.
	TABLE_OWNER       string //Object owner name. This field always returns a value.
	TABLE_NAME        string //Object name. This field always returns a value.
	COLUMN_NAME       string //Column name, for each column of the TABLE_NAME returned. This field always returns a value.
	DATA_TYPE         int    //Integer code for ODBC data type. If this is a data type that cannot be mapped to an ODBC type, it is NULL. The native data type name is returned in the TYPE_NAME column.
	TYPE_NAME         string //String representing a data type. The underlying DBMS presents this data type name.
	PRECISION         int    //Number of significant digits. The return value for the PRECISION column is in base 10.
	LENGTH            int    //Transfer size of the data.1
	SCALE             int    //Number of digits to the right of the decimal point.
	RADIX             int    //Base for numeric data types.
	NULLABLE          int    //Specifies nullability.  1 = NULL is possible.  0 = NOT NULL.
	REMARKS           string //This field always returns NULL.
	COLUMN_DEF        string //Default value of the column.
	SQL_DATA_TYPE     int    //Value of the SQL data type as it appears in the TYPE field of the descriptor. This column is the same as the DATA_TYPE column, except for the datetime and SQL-92 interval data types. This column always returns a value.
	SQL_DATETIME_SUB  int    //Subtype code for datetime and SQL-92 interval data types. For other data types, this column returns NULL.
	CHAR_OCTET_LENGTH int    //Maximum length in bytes of a character or integer data type column. For all other data types, this column returns NULL.
	ORDINAL_POSITION  int    //Ordinal position of the column in the object. The first column in the object is 1. This column always returns a value.

	//Nullability of the column in the object. ISO rules are followed to determine nullability. An ISO SQL-compliant DBMS cannot return an empty string.
	//YES = Column can include NULLS.
	//NO = Column cannot include NULLS.
	//This column returns a zero-length string if nullability is unknown.
	//The value returned for this column is different from the value returned for the NULLABLE column.
	IS_NULLABLE string

	SS_DATA_TYPE int //SQL Server data type used by extended stored procedures. For more information, see Data Types (Transact-SQL).
}

func (col *MSSqlColumn) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&col.TABLE_QUALIFIER,
		&col.TABLE_OWNER,
		&col.TABLE_NAME,
		&col.COLUMN_NAME,
		&col.DATA_TYPE,
		&col.TYPE_NAME,
		&col.PRECISION,
		&col.LENGTH,
		&col.SCALE,
		&col.RADIX,
		&col.NULLABLE,
		&col.REMARKS,
		&col.COLUMN_DEF,
		&col.SQL_DATA_TYPE,
		&col.SQL_DATETIME_SUB,
		&col.CHAR_OCTET_LENGTH,
		&col.ORDINAL_POSITION,
		&col.IS_NULLABLE,
		&col.SS_DATA_TYPE,
	)
}

func ToColumn(col *MSSqlColumn) Column {
	return Column{
		OriginalName: col.COLUMN_NAME,
		NewName:      NameToPsql(col.COLUMN_NAME),
		col:          col,
	}
}

type MSSqlFKey struct {
	PKTABLE_QUALIFIER string // Name of the table (with the primary key) qualifier. This field can be NULL.
	PKTABLE_OWNER     string // Name of the table (with the primary key) owner. This field always returns a value.
	PKTABLE_NAME      string // Name of the table (with the primary key). This field always returns a value.
	PKCOLUMN_NAME     string // Name of the primary key columns, for each column of the TABLE_NAME returned. This field always returns a value.
	FKTABLE_QUALIFIER string // Name of the table (with a foreign key) qualifier. This field can be NULL.
	FKTABLE_OWNER     string // Name of the table (with a foreign key) owner. This field always returns a value.
	FKTABLE_NAME      string // Name of the table (with a foreign key). This field always returns a value.
	FKCOLUMN_NAME     string // Name of the foreign key column, for each column of the TABLE_NAME returned. This field always returns a value.
	KEY_SEQ           int    //Sequence number of the column in a multicolumn primary key. This field always returns a value.

	//Action applied to the foreign key when the SQL operation is an update. SQL Server returns 0 or 1 for these columns:
	//0=CASCADE changes to foreign key.
	//1=NO ACTION changes if foreign key is present.
	UPDATE_RULE int
	//Action applied to the foreign key when the SQL operation is a deletion. SQL Server returns 0 or 1 for these columns:
	//0=CASCADE changes to foreign key.
	//1=NO ACTION changes if foreign key is present.
	DELETE_RULE int

	FK_NAME string //Foreign key identifier. It is NULL if not applicable to the data source. SQL Server returns the FOREIGN KEY constraint name.
	PK_NAME string //Primary key identifier. It is NULL if not applicable to the data source. SQL Server returns the PRIMARY KEY constraint name.
}

func (key *MSSqlFKey) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&key.PKTABLE_QUALIFIER,
		&key.PKTABLE_OWNER,
		&key.PKTABLE_NAME,
		&key.PKCOLUMN_NAME,
		&key.FKTABLE_QUALIFIER,
		&key.FKTABLE_OWNER,
		&key.FKTABLE_NAME,
		&key.FKCOLUMN_NAME,
		&key.KEY_SEQ,
		&key.UPDATE_RULE,
		&key.DELETE_RULE,
		&key.FK_NAME,
		&key.PK_NAME,
	)
}

type MSSqlPKey struct {
	TABLE_QUALIFIER string //Name of the table qualifier. This field can be NULL.
	TABLE_OWNER     string //Name of the table owner. This field always returns a value.
	TABLE_NAME      string //Name of the table. In SQL Server, this column represents the table name as listed in the sysobjects table. This field always returns a value.
	COLUMN_NAME     string //Name of the column, for each column of the TABLE_NAME returned. In SQL Server, this column represents the column name as listed in the sys.columns table. This field always returns a value.
	KEY_SEQ         int    //Sequence number of the column in a multicolumn primary key.
	PK_NAME         string //Primary key identifier. Returns NULL if not applicable to the data source.

}

func (key *MSSqlPKey) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&key.TABLE_QUALIFIER,
		&key.TABLE_OWNER,
		&key.TABLE_NAME,
		&key.COLUMN_NAME,
		&key.KEY_SEQ,
		&key.PK_NAME,
	)
}
