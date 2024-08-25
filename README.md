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

![Birds Eye View](images/architecture-birds-eye-view.png)

The image above is a visual depiction of the `DatabaseEngine`, which is
responsible for glueing together independent layers and providing an interface
for executing SQL statements and `meta` commands (e.g `.dbinfo`).

There are 5 layers: `Pager`, `BTreeEngine`, `Machine`, `Generator` and `Parser`. Each
of these layers is designed to be completely independent of other the layers. A
layer is responsible for defining the interface by which it expects to make calls
to another layer.

This means that layers are decoupled, and don't really on the implementation
details of any other layers.

In the future, I'd like to potentially create different implementations of some
layers and figure out a way to profile them (i.e compare CPU and memory usage).

### Layer 1: Pager

The pager is responsible for efficiently reading and writing pages to/from a file,
and is not concerned with things like:

1. The contents of a page
2. The format of the file
3. The 'context' of the file (e.g a database)

It could be used as an interface for interacting with _any_ page-based file, and
conforms to the following interface:

```go
type Pager interface {
    Close() error
    PageSize() uint64
    ReservedSpace() uint64
    GetPage(pageNum uint64) ([]byte, error)
}
```

It can be instantied with the following constructor:

```go
func NewPager(filepath string, config PagerConfig) (*Pager, error)
```

where the `config` contains the page size and the number of reserved bytes at the
end of each page.

The pager has a cache, which is implemented as a hash map that relates page
numbers to a byte slice (which is the length of a page).

There's likely a huge number of optimisations that could be done here, and I'm
deliberately doing none of them now (other than the very simple cache described
above) because my primary aim is to get an end-to-end simple, extensible and
functional system.

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
