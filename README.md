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

```bash
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

### Architecture of my VBDE


SQLite supports loads of bytecode operations - many of which I'm guessing aren't
needed for the basic operations that I want to perform (I'm sure it's _possible_
to even do very advanced operations with a very small subset of these as well).

I think the best way to build this is to choose a very small subset of
operations which would allow me to execute very basic SQL statements, end to
end. The most simple SQL statement that I can think of is `SELECT 1;` which just
returns the value 1.

The bytecode for this SQL statement looks like this:

addr opcode p1 p2 p3 p4 p5 comment

---

0 Init 0 4 0 0 Start at 4
1 Integer 1 1 0 0 r[1]=1
2 ResultRow 1 1 0 0 output=r[1]
3 Halt 0 0 0 0
4 Goto 0 1 0 0

Each instruction always has 5 operands. I could represent an instruction like:

```
type Instruction struct {
	Opcode String
	P1	   int
	P2	   int
	P3	   int
	P4	   int
	P5	   int
}
```

Ideally with something more like an Enum for the Opcode.

```
var instructions []Instruction

for _, instruction := range instructions {
	machine = machine.Execute(instruction)
}


```

```
type Machine struct {
	Registers map[int]([]byte)
	Out		  ____
}
```

```
addr  opcode         p1    p2    p3    p4             p5  comment
----  -------------  ----  ----  ----  -------------  --  -------------
0     Integer        1     1     0                    0   r[1]=1
1     ResultRow      1     1     0                    0   output=r[1]
2     Halt           0     0     0                    0
```

- engine
  | - instructions
  | - Instruction.go
  | - Halt.go
  | - Integer.go
  | - ResultRow.go
  | - machine.go

addr opcode p1 p2 p3 p4 p5 comment

---

0 Init 0 8 0 0 Start at 8
1 OpenRead 0 2 0 2 0 root=2 iDb=0; users
2 Rewind 0 7 0 0
3 Column 0 0 1 0 r[1]= cursor 0 column 0
4 Column 0 1 2 0 r[2]= cursor 0 column 1
5 ResultRow 1 2 0 0 output=r[1..2]
6 Next 0 3 0 1
7 Halt 0 0 0 0
8 Transaction 0 0 1 0 1 usesStmtJournal=0
9 Goto 0 1 0 0
