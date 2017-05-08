NAME
     mssql_migrate -- copy MS Sql Server Database to a Postgres Database

SYNOPSIS
     mssql_migrate [--drop] [--print] <from> <to> <table> [table ...]

DESCRIPTION

     Arguments <from> and <to> are both URI's referencing the MSSQL and posgres
     databases respectively, examples as follows:

          sqlserver://user:password@hostname?database=DatabaseName&other_options=true
          postgres://user:password@hostname/database?options=true

     Options:

     --drop    Drop tables before creating them

     --print   Dont execute, only print the creation SQL

TODO
     * Add NOT NULL to fields
     * Add Foreign keys
