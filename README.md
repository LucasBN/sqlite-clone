# SQLite3 Clone - A Serverless Single File Database

A simple SQL database engine written in Go, which reads and stores data on a 
single file in the SQLite3 [file format](https://www.sqlite.org/fileformat.html#storage_of_the_sql_database_schema). 

You can run the database engine with the following command:
```
./sqlite sample.db
```

This project is in it's early stages of development, which means that it does 
not support many features. 

The section on architecture below will go into more depth about the various 
layers that make up the engine, but at present there are only implementations 
for the `Machine`, `BTreeEngine` and `Pager` layers.

This means that you cannot actually pass arbitrary SQL - the `Parser` and
`Generator` layers are mocked and return hard coded instructions. 

# Architecture

### Birds Eye View

TODO

### Layer 1: Pager

TODO

### Layer 2: BTreeEngine

TODO

### Layer 3: Machine

TODO

### Layer 4: Bytecode Generator

TODO

### Layer 5: Parser

TODO

# Useful links

- [SQLite3 File Format](https://www.sqlite.org/fileformat.html#storage_of_the_sql_database_schema)
- [SQLite3 Opcodes](https://www.sqlite.org/opcode.html)