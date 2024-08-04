# Reading the database file

I've made a few design decisions related to how pages, cells and records are
stored, as well as how tables are read.

A cell is treated as something that exists entirely within the context of a
b-tree page, and the type of the cell is entirely dependent on the type of
b-tree page. A b-tree page looks like:

```
type BTreePage struct {
	Header         BTreeHeader
	IntIdxCells    []IntIdxCell
	IntTableCells  []IntTableCell
	LeafIdxCells   []LeafIdxCell
	LeafTableCells []LeafTableCell
}
```

The header contains information like the page type, number of cells and cell
content offset (among others, as defined by the SQLite database file format).

I briefly considered separating out cells and then attempting to traverse all
overflow pages to bring the entire payload together, but I decided this was a
really bad idea because:

i. I'm doing this inside a function called `readBTreePage` which takes a page
number and it feels weird for this to make subsequent calls to read other pages
and to return data that is actually stored on other pages.

ii. I'm not sure how this would end up looking for interior pages (which I
haven't yet implemented). My suspicion is that it would be quite horrible.

I decided to go for an approach which seems much more sensible: the
`readBTreePage` function would only read and return the data on the page that it
was asked about, and the caller would be responsible for interpreting that data.

Right now, the only caller is `readTable`, which takes a `rootPage`. This
function first calls `readBTreePage` on the root page itself, but of course the
data of the table may be split across multiple pages - which we only find out
after we have read the root page. After reading the root page, `readTable` will
read all other necessary pages to create a list of 'raw records' which are
essentially byte slices that contain _all_ of the data for a record and have no
dependency or reference to pages. This allows us to pass the entire raw record
to a function that is responsible for decoding the record into a useable format
(see `record.go`).

So far, I think this is a good structure and will allow me to (fingers crossed)
easily add support for things like overflown cells and interior pages

# Virtual Database Engine (VDBE)

I've now got my hands dirty with reading data from a simple SQLite database file
and feel somewhat confident that I could continue down this path to flesh out
the functionality a bit (i.e handling interior pages and supporting indexes).

However, I'm going to switch my attention to a very different aspect: bytecode.
Unlike database engines like PostgreSQL and MySQL which execute SQL by walking a
tree of objects (similar but not the same as an AST), SQLite executes SQL by
executing bytecode on a virtual machine.

The [SQLite docs](https://www.sqlite.org/whybytecode.html) go into a lot of
detail as to why the developers made this choice, but to briefly summarise:

1. Easier to understand: linear and "atomic" instructions
2. Easier to debug: clearer separation between frontend and backend
3. Can be run incrementally (important since it runs locally, not on a server)
4. Bytecode is smaller than the AST representation (important for caching)
5. It _might_ be faster but a fair comparison is difficult

### What does SQLite bytecode look like?

It's really easy to see what bytecode is produced by SQLite:

```
$ sqlite3
sqlite> EXPLAIN SELECT 1;
| addr | opcode    | p1 | p2 | p3 | p4 | p5 |
|------|-----------|----|----|----|----|----|
| 0    | Init      | 0  | 4  | 0  | 0  | 0  | <- Jump to address 4
| 1    | Integer   | 1  | 1  | 0  | 0  | 0  | <- Put the value 1 into register 1
| 2    | ResultRow | 1  | 1  | 0  | 0  | 0  | <- Output the value of register 1
| 3    | Halt      | 0  | 0  | 0  | 0  | 0  | <- Halt execution
| 4    | Goto      | 0  | 1  | 0  | 0  | 0  | <- Jump to address 1
```

This output gives me a fairly clear idea on how to make a virtual machine
capable of running such instructions. Since I'm more interested in the actual
internals of a database rather than the higher level interfaces, I'm going to
continue working "back to front" - which means that for now I'm not going to
spend any time worrying about writing an SQL parser that generates bytecode, I'm
just going to write a simple virtual machine that can execute a subset of the
instructions described in the SQLite docs.

### Architecture of my VDBE

I'm unsure about my initial implementation (I don't know if it's "efficient" or
not - and I don't know if that really even matters?), but the directory
structure looks like:

```
machine
|-- instructions
|    +-- instruction.go
|    +-- integer.go
|    +-- halt.go
|    +-- result_row.go
|-- registers
|    +-- register.go
|-- state
|    +-- state.go
+-- machine.go
```

I've also created a few data types (structs) to represent the machine:

1. Machine: machine state, the program and the output buffer
2. State: current address, registers, halted
3. RegisterFile: map from int to int
4. Instruction: interface with an execute function

The purpose of creating these abstractions was to get to a point where I could
extend the capabilities of the VM by only having to define the execute function
on new instructions.

A 'Machine' has a `Run()` function on it which repeatedly calls execute on the
instruction at the current address (instruction pointer) until the halt state is
reached. Each call to `Execute` is passed a pointer to the current machine
state, and updates the state according to the instructions specification.

The `State` is not directly embedded in the `Machine` struct because the state
is relevant to instructions (and is therefore imported by the instructions
package), but the other fields on the machine (output and program) are not
relevant to execution of a single instruction (and therefore don't need to be
accessible / imported). This enforces a separation of concerns (and avoids
import cycles).

The following instructions are sufficient to execute bytecode that represents
the query `SELECT 1, 2;`:

```go
instructions := []instructions.Instruction{
	instructions.Integer{Register: 1, Value: 1},
	instructions.Integer{Register: 2, Value: 2},
	instructions.ResultRow{FromRegister: 1, ToRegister: 2},
	instructions.Halt{},
}
m := machine.Init(instructions)
output := m.Run()
```

This is cool and lays down the foundations for a VDBE - we now need to implement
instructions that are able to interact with the database file!
