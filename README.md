# Architecture & Design Decisions

The entry point is `main.go`, which is responsible for opening the database
file, reading the database header and handling the execution of the command by
handing the responsibility to another function.

## Reading the SQLite Schema

Inside `schema.go` is the function `readSQLiteSchema` which takes the database
file and the database header, and returns a `SQLiteSchema` struct which contains
a list of `SQLiteSchemaTuple`s which corresponds to the rows on the 
sqlite_schema table.

My aim for this function is that it reads the first page of the database, and
then 'bubbles' 


The first page is either a leaf page or an interior page.

What I'm really interested in is all of the cells that make up the entire table,
and when I come to decode these cells, I don't really care that they're on
separate pages.

func readSQLiteSchema() {
    // Read the first page

    // Get all of the cells relevant to the table

    // Convert each cell into a record, and then into a SQLiteSchemaTuple
}