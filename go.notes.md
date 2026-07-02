## Mimicking Classes:
---------------------

- Go does not have classes.
- Instead, we use structs and methods.
- Structs hold data.
- Methods provide behavior.
- Constructor-like functions are typically named `NewSomething()`.

Example:

```go
type DataBase struct {
	filePath string
}

func NewDataBase(path string) *DataBase {
	return &DataBase{
		filePath: path,
	}
}
```

Usage:

```go
db := NewDataBase("db.log")
```

## Pass by Value and Pass by Reference:
---------------------------------------

### Python

- Variables hold references to objects by default.

Example:

```python
class DB:
    def __init__(self):
        self.name = "old"

def change(db):
    db.name = "new"

db = DB()

change(db)

print(db.name)
```

Output:

```text
new
```

The function receives a reference to the same object.

---

### Go

- Variables hold copies of values by default.

Example:

```go
type DB struct {
	name string
}

func change(db DB) {
	db.name = "new"
}

func main() {
	db := DB{
		name: "old",
	}

	change(db)

	fmt.Println(db.name)
}
```

Output:

```text
old
```

The struct was copied.

---

### Using Pointers

```go
type DB struct {
	name string
}

func change(db *DB) {
	db.name = "new"
}

func main() {
	db := DB{
		name: "old",
	}

	change(&db)

	fmt.Println(db.name)
}
```

Output:

```text
new
```

---

### Rule of Thumb

- Go copies values by default.
- Use pointers when multiple parts of the program should operate on the same struct instance.
- Use pointers for large structs and shared state.
- Use values for small immutable-like data.


## Error Handling Basics:
-------------------------

- Go does not use exceptions for normal error handling.
- Functions typically return:

```go
result, err
```

Example:

```go
file, err := os.Open("db.log")
```

- Always check errors immediately.

```go
if err != nil {
	return err
}
```

---

### Creating Errors

```go
err := fmt.Errorf("something went wrong")
```

---

### Wrapping Errors

```go
if err != nil {
	return fmt.Errorf(
		"failed to open database file: %w",
		err,
	)
}
```

Output:

```text
failed to open database file: open db.log: no such file or directory
```

---

### %v vs %w

`%v`

```go
fmt.Errorf("error: %v", err)
```

- Converts error to text.
- Loses original error information.

`%w`

```go
fmt.Errorf("error: %w", err)
```

- Converts error to text.
- Preserves original error.
- Allows:

```go
errors.Is(...)
errors.As(...)
```

---

### Multiple Error Checks

```go
file, err := os.Open(...)

if err != nil {
	return fmt.Errorf(
		"failed to open file: %w",
		err,
	)
}

_, err = file.WriteString(...)

if err != nil {
	return fmt.Errorf(
		"failed to write record: %w",
		err,
	)
}
```

Each `if err != nil` checks the operation immediately before it.


## File Handling Basics:
------------------------

### Open File

```go
file, err := os.Open("db.log")
```

---

### Open File For Writing

```go
file, err := os.OpenFile(
	"db.log",
	os.O_APPEND|os.O_CREATE|os.O_WRONLY,
	0644,
)
```

---

### File Flags

```go
os.O_APPEND
```

Write at end of file.

```go
os.O_CREATE
```

Create file if it does not exist.

```go
os.O_WRONLY
```

Write only.

```go
os.O_RDONLY
```

Read only.

```go
os.O_RDWR
```

Read and write.

---

### Close Files

```go
file, err := os.Open(...)

if err != nil {
	return err
}

defer file.Close()
```

---

### What defer Means

```go
defer file.Close()
```

means:

```text
When this function exits,
close the file automatically.
```

Example:

```go
func demo() {
	fmt.Println("1")

	defer fmt.Println("3")

	fmt.Println("2")
}
```

Output:

```text
1
2
3
```

---

### Writing To File

```go
_, err = file.WriteString(
	"SET name namah\n",
)
```

---

### Writing Formatted Data

```go
record := fmt.Sprintf(
	"SET %s %s\n",
	key,
	value,
)

_, err = file.WriteString(record)
```

---

### Read Entire File

```go
data, err := os.ReadFile("db.log")
```

Convert bytes to string:

```go
fmt.Println(string(data))
```

---

### Read Line By Line

```go
scanner := bufio.NewScanner(file)

for scanner.Scan() {
	line := scanner.Text()

	fmt.Println(line)
}
```

---

### Why Scanner?

Good for large files.

Instead of loading:

```text
5 GB
```

into memory, it reads:

```text
1 line
↓
1 line
↓
1 line
```

## Short Variable Declarations (`:=`) vs Assignment (`=`)

---

### What does `:=` mean?

Many beginners think:

```go
x := 10
```

means:

> Infer the type of x.

While type inference does happen, the real meaning is:

> Declare a new variable and assign a value.

Example:

```go
name := "Namah"
age := 30
```

Equivalent to:

```go
var name string = "Namah"
var age int = 30
```

---

### What does `=` mean?

`=` assigns a value to an existing variable.

Example:

```go
var name string

name = "Namah"
```

or

```go
name := "Namah"

name = "John"
```

---

### Rule of Thumb

Use:

```go
:=
```

when creating a variable for the first time.

Use:

```go
=
```

when updating an existing variable.

---

### Common Example

```go
file, err := os.Open("db.log")

if err != nil {
	return err
}

_, err = file.WriteString("hello")
```

The first line creates:

```go
file
err
```

The second operation reuses:

```go
err
```

so we use:

```go
=
```

instead of:

```go
:=
```

---

### Why This Error Happens

Example:

```go
file, err := os.Open("db.log")

_, err := file.WriteString("hello")
```

Produces:

```text
no new variables on left side of :=
```

Reason:

```go
err
```

already exists.

And:

```go
_
```

is not a real variable.

Therefore there are no new variables being declared.

---

### Redeclaration Rule

Go allows:

```go
file, err := os.Open(...)

count, err := file.WriteString(...)
```

because:

```go
count
```

is a new variable.

Rule:

> At least one variable on the left side of `:=` must be new.

Valid:

```go
a := 1

a, b := 2, 3
```

Invalid:

```go
a := 1

a := 2
```

---

### Mental Model

```text
:=  -> Declare + Assign

=   -> Assign to Existing Variable
```

Type inference is a side effect of creating the variable.

The primary purpose of `:=` is variable declaration.


## Functions vs Methods in Go

---

### Normal Functions

A normal function belongs to the package.

Example:

```go
func SetKey(
	db *DataBase,
	key string,
	value string,
) error {
	...
}
```

Usage:

```go
SetKey(db, "name", "namah")
```

The database instance must be passed explicitly.

---

### Methods

A method belongs to a type.

Example:

```go
func (db *DataBase) Set(
	key string,
	value string,
) error {
	...
}
```

Usage:

```go
db.Set("name", "namah")
```

The receiver:

```go
(db *DataBase)
```

is similar to:

```python
self
```

in Python.

---

### Understanding the Receiver

Example:

```go
func (db *DataBase) Get(
	key string,
) (string, error) {
	...
}
```

Breakdown:

```go
func
```

Defines a function.

```go
(db *DataBase)
```

Receiver.

Indicates that this method belongs to the `DataBase` type.

```go
Get
```

Method name.

```go
key string
```

Input parameter.

```go
(string, error)
```

Return values.

This reads as:

> Define a method named Get that belongs to DataBase, takes a string key, and returns a string and an error.

---

### Methods vs Functions

These are conceptually similar:

```go
func SetKey(
	db *DataBase,
	key string,
	value string,
) error {
	...
}
```

```go
func (db *DataBase) Set(
	key string,
	value string,
) error {
	...
}
```

The method version simply associates the behavior with the type.

---

### Why Use Methods?

Suppose we have:

```go
type DataBase struct {
	filePath string
}
```

Which is easier to read?

Function style:

```go
SetKey(db, "name", "namah")
GetKey(db, "name")
DeleteKey(db, "name")
```

Method style:

```go
db.Set("name", "namah")
db.Get("name")
db.Delete("name")
```

Methods clearly indicate that the operations belong to the database.

---

### Why Constructors Are Normal Functions

Methods require an existing instance.

Example:

```go
func (db *DataBase) Set(...) error
```

assumes:

```go
db
```

already exists.

A constructor's job is to create the object.

Therefore constructors are usually normal functions:

```go
func NewDataBase(
	path string,
) *DataBase {

	return &DataBase{
		filePath: path,
	}
}
```

Usage:

```go
db := NewDataBase("db.log")
```

---

### Pointer Receivers

Most database methods should use:

```go
func (db *DataBase) Set(...)
```

instead of:

```go
func (db DataBase) Set(...)
```

Reason:

* Avoid copying the struct.
* Allow modification of the original object.
* Keep all callers operating on the same database instance.

---

### Automatic Pointer Conversion

Suppose:

```go
db := DataBase{}
```

and:

```go
func (db *DataBase) Set(...)
```

exists.

This works:

```go
db.Set(...)
```

Go automatically converts:

```go
db.Set(...)
```

to:

```go
(&db).Set(...)
```

when needed.

---

### Rule of Thumb

Use methods when behavior belongs to a specific instance:

```go
db.Set(...)
db.Get(...)
db.Delete(...)
```

Use normal functions when behavior does not belong to an instance:

```go
NewDataBase(...)
ParseConfig(...)
LoadEnv(...)
```

Most business logic around a struct should be implemented as methods.

# Maps and `make()`

## Why `make()` Exists

Some Go types need memory allocation before use.

Example:

```go
var index map[string]int64
```

This only declares the map.

It is:

```text
nil
```

Trying to write:

```go
index["test"] = 100
```

causes:

```text
panic: assignment to entry in nil map
```

---

## Allocating the Map

```go
index := make(map[string]int64)
```

Now memory is allocated and the map can be used.

```go
index["test"] = 100
```

works.

---

## Why Not Just Use `var`?

```go
var index map[string]int64
```

creates:

```text
Declaration only
```

No backing storage exists.

```go
index = make(map[string]int64)
```

creates the backing storage.

---

## Types That Need `make()`

### Maps

```go
make(map[string]int64)
```

### Slices

```go
make([]string, 0)
```

### Channels

```go
make(chan int)
```

---

## Rule of Thumb

```text
Structs:
    var s Struct

Maps:
    make(map[K]V)

Slices:
    make([]T, size)

Channels:
    make(chan T)
```

---

# Command Line Arguments

## Accessing Arguments

Go provides:

```go
os.Args
```

Type:

```go
[]string
```

Example:

```bash
./learn_db set name namah
```

Produces:

```text
[
    "./learn_db",
    "set",
    "name",
    "namah",
]
```

---

## Common Pattern

```go
if len(os.Args) < 2 {
    return
}

command := os.Args[1]
```

---

## Slicing Arguments

To skip executable and command:

```go
args := os.Args[2:]
```

Example:

```text
[
    "./learn_db",
    "set",
    "name",
    "namah",
]
```

becomes:

```text
[
    "name",
    "namah",
]
```

---

# Slices

## Creating a Slice

```go
names := []string{
    "john",
    "jane",
    "bob",
}
```

---

## Slice Syntax

```go
arr[start:end]
```

Examples:

```go
arr := []string{
    "a",
    "b",
    "c",
    "d",
}
```

```go
arr[1:]
```

Result:

```text
[b c d]
```

```go
arr[:2]
```

Result:

```text
[a b]
```

```go
arr[1:3]
```

Result:

```text
[b c]
```

---

## Mental Model

```text
[start, end)
```

Start inclusive.

End exclusive.

---

# Methods as Values

## Methods Can Be Stored

Given:

```go
func (db *DataBase) GetKey(
    key string,
) (string, error)
```

Then:

```go
handler := db.GetKey
```

has type:

```go
func(
    string,
) (string, error)
```

Notice:

```text
db is already attached.
```

---

## Method Values

```go
db.GetKey
```

Produces:

```go
func(string) (string, error)
```

Receiver already bound.

---

## Method Expressions

```go
(*DataBase).GetKey
```

Produces:

```go
func(
    *DataBase,
    string,
) (string, error)
```

Receiver not bound.

---

## Mental Model

```text
db.GetKey
    =
receiver attached

(*DataBase).GetKey
    =
receiver must be passed manually
```

---

# Interfaces

## What Is An Interface?

An interface describes behavior.

Example:

```go
type Engine interface {
    SetKey(
        key string,
        value string,
    ) error

    GetKey(
        key string,
    ) (string, error)
}
```

---

## Implicit Implementation

Unlike Java:

```java
class DB implements Engine
```

Go does not require:

```go
implements Engine
```

anywhere.

---

Given:

```go
type DataBase struct {
    filePath string
}
```

and:

```go
func (db *DataBase) SetKey(...) error
func (db *DataBase) GetKey(...) (string, error)
```

Go automatically determines:

```text
*DataBase satisfies Engine
```

because all required methods exist.

---

## Why Use Interfaces?

Without interfaces:

```go
func RunCommand(
    db *logengine.DataBase,
)
```

Only works with:

```text
logengine.DataBase
```

---

With interfaces:

```go
func RunCommand(
    db engine.Engine,
)
```

Works with:

```text
Log Engine
BTree Engine
LSM Engine
```

as long as they implement:

```go
SetKey(...)
GetKey(...)
```

---

## Interfaces Describe Capabilities

Struct:

```text
What data exists?
```

Interface:

```text
What actions are available?
```

---

## Important: Do Not Use Pointers To Interfaces

Bad:

```go
*engine.Engine
```

Good:

```go
engine.Engine
```

Interfaces already contain a reference to the concrete value.

---

## Rule of Thumb

```text
Struct:
    Usually use *Struct

Interface:
    Usually use Interface
```

---

# Scanners and Iteration

## Reading Line By Line

```go
scanner := bufio.NewScanner(file)
```

---

## Iterating

```go
for scanner.Scan() {
    line := scanner.Text()

    fmt.Println(line)
}
```

---

## What Does Scan() Return?

```go
scanner.Scan()
```

returns:

```go
bool
```

Meaning:

```text
true  -> another line exists

false -> EOF or error
```

---

## Where Is The Current Value?

The current line is stored internally inside the scanner.

Retrieve it using:

```go
scanner.Text()
```

Example:

```go
for scanner.Scan() {
    line := scanner.Text()

    fmt.Println(line)
}
```

---

## Mental Model

Python:

```python
for line in file:
```

Go:

```go
for scanner.Scan() {
    line := scanner.Text()
}
```

---

# Package Visibility (Exported vs Unexported)

## Exported Names

Names starting with uppercase letters are exported.

Example:

```go
func SetKey(...)
```

Accessible from other packages.

---

## Unexported Names

Names starting with lowercase letters are private.

Example:

```go
func setKey(...)
```

Only accessible inside the package.

---

Examples:

```go
type DataBase struct {}
```

Exported.

```go
type dataBase struct {}
```

Private.

---

```go
func CreateDatabase(...)
```

Exported.

---

```go
func createDatabase(...)
```

Private.

---

## Rule of Thumb

```text
Uppercase:
    Public

Lowercase:
    Private
```

## Why We Usually Do Not Use Pointers To Interfaces

When learning Go, it is natural to think:

```go
func RunCommand(
	db *engine.Engine,
)
```

because we have learned that structs are usually passed as pointers.

However, interfaces are different.

---

### Structs Need Pointers

Given:

```go
type DataBase struct {
	filePath string
}
```

Passing:

```go
func DoSomething(
	db DataBase,
)
```

creates a copy of the struct.

To avoid copying and to modify the original object, we use:

```go
func DoSomething(
	db *DataBase,
)
```

---

### Interfaces Already Hold A Reference

Given:

```go
type Engine interface {
	SetKey(...)
	GetKey(...)
}
```

and:

```go
db := CreateDatabase(...)
```

we can write:

```go
func RunCommand(
	db engine.Engine,
)
```

Even though we are not using a pointer.

Why?

Because an interface internally contains:

```text
1. The concrete type
2. A reference to the concrete value
```

Conceptually:

```text
Engine
┌─────────────────────────┐
│ Type: *DataBase         │
│ Value: 0x12345678       │
└─────────────────────────┘
```

The interface already knows:

- What type it contains.
- Where the actual object lives in memory.

---

### Think Of An Interface As A Box

Suppose:

```go
var db engine.Engine
```

After assignment:

```go
db = CreateDatabase(...)
```

Conceptually:

```text
db
│
▼

┌─────────────────────────┐
│ Type: *DataBase         │
│ Value: 0x12345678       │
└─────────────────────────┘
```

The interface is already a box containing a pointer to the real object.

---

### What Happens With A Pointer To An Interface?

Suppose we write:

```go
*engine.Engine
```

Now we have:

```text
Pointer
│
▼

Interface Box
│
▼

Pointer To DataBase
│
▼

Actual DataBase
```

Visualized:

```text
*Engine
    │
    ▼
 Engine
    │
    ▼
 *DataBase
    │
    ▼
 DataBase
```

This extra level of indirection is almost never useful.

---

### Why This Causes Errors

Suppose:

```go
func RunCommand(
	db *engine.Engine,
)
```

and later:

```go
db.SetKey(...)
```

Go reports:

```text
db.SetKey undefined
(type *engine.Engine is pointer to interface)
```

because:

```text
The methods belong to the interface,
not to a pointer to the interface.
```

---

### Correct Usage

Bad:

```go
func RunCommand(
	db *engine.Engine,
)
```

Good:

```go
func RunCommand(
	db engine.Engine,
)
```

---

### Rule Of Thumb

```text
Struct:
    Usually use *Struct

Interface:
    Usually use Interface
```

If you ever see:

```go
*SomeInterface
```

it is often a sign that the design should be reconsidered.

---

### Mental Model

Struct:

```text
Need pointer?
Usually yes.
```

Interface:

```text
Already contains a reference.
No additional pointer needed.
```

# Database Engine Evolution Notes

## Why BuildIndex Exists

When the database starts:

```go
db := NewDatabase(...)
db.BuildIndex()
```

We scan the entire log once and reconstruct:

```go
map[string]int64
```

Example:

```text
db.log

test,value1
user,hello
test,value2
```

Index becomes:

```go
{
    "test": 22,
    "user": 12,
}
```

The latest occurrence wins.

After startup:

```go
SetKey()
```

maintains the index in memory.

No more rebuilding.

---

## What is an Offset?

An offset is:

> Number of bytes from the beginning of the file.

Example:

```text
0         10        20
|---------|---------|

test,value1
user,hello
```

If:

```text
user,hello
```

starts at byte 12:

```go
index["user"] = 12
```

---

## How Reads Work

Instead of:

```go
Scan entire file
```

we do:

```go
offset := index["user"]

file.Seek(offset, io.SeekStart)

Read one record
```

Read complexity becomes:

```text
O(1)
```

instead of:

```text
O(n)
```

---

# os.Open vs os.OpenFile

## os.Open

```go
os.Open("db.log")
```

Equivalent to:

```go
os.OpenFile(
    "db.log",
    os.O_RDONLY,
    0,
)
```

Read only.

---

## os.OpenFile

Allows:

```go
O_RDONLY
O_WRONLY
O_RDWR
O_CREATE
O_APPEND
O_TRUNC
```

Example:

```go
os.OpenFile(
    "db.log",
    os.O_RDWR|os.O_CREATE,
    0644,
)
```

---

# What Does Seek Do?

Files have a cursor.

```text
cursor
  ↓
test,value1
user,hello
```

Move cursor:

```go
file.Seek(0, io.SeekStart)
```

Move to beginning.

---

```go
file.Seek(0, io.SeekEnd)
```

Move to EOF.

Used for appends.

---

```go
file.Seek(0, io.SeekCurrent)
```

Move 0 bytes relative to current position.

Effectively:

> Tell me where I am.

---

# Why Seek(0, io.SeekCurrent) Looked Attractive

We thought:

```go
offset, _ := file.Seek(0, io.SeekCurrent)
```

would give:

```text
Current record offset
```

before reading each line.

Unfortunately:

```go
bufio.Reader
```

breaks this assumption.

---

# The Buffered Reader Bug

Suppose:

```go
reader := bufio.NewReader(file)
```

Reader may immediately fetch:

```text
8192 bytes
```

from the OS.

The OS cursor jumps:

```text
0 -> 8192
```

even though you've only processed:

```text
test,value1
```

from memory.

Now:

```go
file.Seek(0, io.SeekCurrent)
```

returns:

```text
8192
```

instead of:

```text
12
```

or wherever you really are.

---

# Two Different Cursors Exist

## OS Cursor

Represents:

```text
Position in file descriptor
```

Example:

```text
8192
```

---

## Buffer Cursor

Represents:

```text
Position inside buffered memory
```

Example:

```text
12
```

These are not always the same.

This was the key insight.

---

# The One-Line Explanation

> The OS cursor tracks bytes read from the file descriptor, while bufio tracks bytes consumed from its in-memory buffer.

Example:

```text
Bytes read from OS      = 8192
Bytes consumed by you   = 12
```

Therefore:

```text
OS Cursor      = 8192
Logical Cursor = 12
```

They diverge.

---

# Why Does bufio Read Ahead?

Because syscalls are expensive.

Without buffering:

```go
Read 1 byte
Read 1 byte
Read 1 byte
...
```

Potentially millions of syscalls.

---

With buffering:

```text
Read 8KB once
```

Then consume data directly from RAM.

Far fewer syscalls.

Much faster.

---

# Grocery Store Analogy

Without buffering:

```text
Need tomato
→ Drive to store

Need onion
→ Drive to store

Need milk
→ Drive to store
```

Many trips.

---

With buffering:

```text
Need tomato
```

Buffer says:

```text
While I'm there,
I'll grab tomatoes,
onions,
milk,
bread,
eggs.
```

Now future requests come from the fridge.

No extra trips.

---

# What Happens When File Is Larger Than 8KB?

Suppose:

```text
100 MB file
```

Buffer loads:

```text
First 8KB
```

Consume data.

When buffer becomes empty:

```text
Load next 8KB
```

Repeat until EOF.

This is called streaming.

We never load the entire file into memory.

---

# Why We Track Offsets Ourselves

Instead of:

```go
file.Seek(0, io.SeekCurrent)
```

we do:

```go
offset += bytesRead
```

Example:

```go
offset := int64(0)

for scanner.Scan() {
    index[key] = offset

    offset += int64(len(scanner.Bytes()) + 1)
}
```

Now:

```text
offset
```

represents our logical position in the file.

Independent of buffering.

---

# Why This Works

An offset is ultimately:

> Number of bytes from the beginning of the file.

Whether we obtain it from:

```go
file.Seek(...)
```

or:

```go
offset += bytesRead
```

doesn't matter.

The number is the same.

---

# Runtime Writes

For writes we use the filesystem offset:

```go
offset, _ := file.Seek(0, io.SeekEnd)

file.WriteString(record)

index[key] = offset
```

This gives the true starting position of the new record.

---

# Storage Engine Lesson

A recurring pattern in databases:

> Don't repeatedly ask the OS where you are. Track your own logical state.

Examples:

- Database page positions
- Kafka offsets
- Log sequence numbers
- Memory allocators
- Storage engine indexes

The database often already knows where it is.

---

# Current Architecture

## Disk

```text
Append-only log
```

## Memory

```go
map[string]int64
```

## Write Path

```text
Append record
Update index
```

Equivalent to:

```go
offset, _ := file.Seek(0, io.SeekEnd)

file.WriteString(record)

index[key] = offset
```

Complexity:

```text
O(1)
```

---

## Read Path

```text
Lookup offset
Seek
Read one record
```

Equivalent to:

```go
offset := index[key]

file.Seek(offset, io.SeekStart)

Read record
```

Complexity:

```text
O(1)
```

---

# Final Insight

The moment this project stopped being "write a file" and became:

- file descriptors
- offsets
- buffering
- syscalls
- logical vs physical position
- append-only logs
- indexing

...it became a storage engine.

And the concepts learned here reappear everywhere:

- Kafka partitions
- PostgreSQL WAL
- Cassandra SSTables
- RocksDB
- Redis AOF
- Filesystems
- Distributed logs

Different systems.

Same ideas.


# Bytes vs Characters vs Runes

## The Question

Does:

```go
len(scanner.Bytes())
```

mean:

> "How many characters did I read?"

Not necessarily.

It means:

> "How many bytes did I read?"

These are not always the same thing.

---

# ASCII Characters

For simple ASCII text:

```text
hello
```

each character occupies:

```text
1 byte
```

Example:

```go
fmt.Println(len("hello"))
```

Output:

```text
5
```

Because:

```text
h = 1 byte
e = 1 byte
l = 1 byte
l = 1 byte
o = 1 byte

Total = 5 bytes
```

In ASCII:

```text
1 character = 1 byte
```

which makes it easy to forget they're different concepts.

---

# Unicode Changes Everything

Not all characters are ASCII.

Example:

```text
é
```

```go
fmt.Println(len("é"))
```

Output:

```text
2
```

because UTF-8 stores:

```text
é = 2 bytes
```

---

Example:

```text
न
```

```go
fmt.Println(len("न"))
```

Output:

```text
3
```

because:

```text
न = 3 bytes
```

---

Example:

```text
🚀
```

```go
fmt.Println(len("🚀"))
```

Output:

```text
4
```

because:

```text
🚀 = 4 bytes
```

---

# UTF-8 Encoding

Go strings use UTF-8.

Different characters occupy different numbers of bytes.

| Character | Bytes |
|------------|---------|
| a | 1 |
| z | 1 |
| 1 | 1 |
| , | 1 |
| é | 2 |
| न | 3 |
| 🚀 | 4 |

---

# len() Measures Bytes

Example:

```go
fmt.Println(len("hello"))
```

Output:

```text
5
```

because:

```text
5 bytes
```

---

Example:

```go
fmt.Println(len("🚀"))
```

Output:

```text
4
```

because:

```text
4 bytes
```

---

# What Is a Rune?

Go represents a Unicode character using:

```go
rune
```

A rune represents:

```text
One Unicode code point
```

Think:

```text
One character
```

---

Example:

```go
s := "🚀"

fmt.Println(len(s))
```

Output:

```text
4
```

because:

```text
4 bytes
```

---

But:

```go
fmt.Println(len([]rune(s)))
```

Output:

```text
1
```

because:

```text
1 character
```

---

Another Example

```go
s := "नमस्ते"

fmt.Println(len(s))
```

might output:

```text
18
```

because UTF-8 uses multiple bytes per character.

But:

```go
fmt.Println(len([]rune(s)))
```

might output:

```text
6
```

because there are six characters.

---

# Why Does scanner.Bytes() Work For Offsets?

Our index builder uses:

```go
offset += int64(len(scanner.Bytes()) + 1)
```

At first glance this seems odd.

Why not count characters?

Because:

> File offsets are measured in bytes.

Not characters.

---

Suppose we have:

```text
name,🚀
```

The filesystem stores:

```text
n
a
m
e
,
🚀
\n
```

as bytes.

The rocket occupies:

```text
4 bytes
```

on disk.

Therefore:

```go
len(scanner.Bytes())
```

correctly measures:

```text
How many bytes this record occupies in the file
```

which is exactly what an offset represents.

---

# Storage Engine Insight

Humans think in:

```text
Characters
Words
Strings
```

Computers think in:

```text
Bytes
Memory Addresses
Offsets
```

The filesystem does not know what a character is.

The filesystem only sees:

```text
01001010
01100101
01101100
...
```

(raw bytes)

---

# Why We Use len(scanner.Bytes())

We are calculating:

```text
Record size on disk
```

not:

```text
Number of visible characters
```

Therefore:

```go
offset += int64(len(scanner.Bytes()) + 1)
```

is correct.

Even for:

```text
hello
नमस्ते
🚀🚀🚀
```

because offsets are measured in bytes.

---

# Final Takeaway

For storage engines:

```go
len(string)
```

means:

```text
Number of bytes
```

not:

```text
Number of characters
```

And that is exactly what we want because:

> A file offset is the number of bytes from the beginning of the file.

# Memory Profiling

- TotalAlloc:
	- Think of it as a speedometer on a bike. It keeps going.
	- Let's say we do `a:=make([]byte, 100)`, `TotalAlloc=100`.
	- Then we do `b:=make([]byte, 200)`, `TotalAlloc=200 + prev 100 = 300`.
	- Then lets say GC cleans up data.  `TotalAlloc is still 300`.
	- It measures all the bytes that have ever been allocated.
- TotalAlloc Usage:
	- Suppose we have a code change, we used to allocate `8kb` before and now we allocate `80kb`.
	- This means our older version created much less garbage.
	- This means:
		- Fewer GC Cycles, less CPU cleaning memory, better throughput

- Mallocs:
	- The number of objects allocated in the heap.
	- Let's say we do `a:=make([]byte, 10)`, `Mallocs += 1`.
	- Then again we do `b:=make([]byte, 2000)`, `Mallocs += 1, still 1`.
	- It does not care about bytes, it cares about the number of objects in heap.
- Malloc Usage:
	- Suppose our program creates 1000 tiny strings vs 1 large buffer.
	- The version which produces 1000 tiny strings has 1000 Mallocs.
	- This means, We have more pointers, more graph traversal for GC, more bookkeeping.
	- Basically, an indicator to how much work am I creating for the GC.

- HeapAlloc:
	- How much live memory exists right now.
	- Let's say we do `make([]byte, 100 MB)`, `HeapAlloc = 100MB`.
	- GC runs later and then `HeapAlloc = 0`.
- HeapAlloc Usage:
	- We can use this to measure memory usage.
	- If our previous version used `8MB` but now uses `5MB`, thats an improvement in memory usage by `3MB`. `3MB` is what is shown in the graph.
	- Another example, lets say We used, `20MB`, but now use `15MB`.

# CPU Profiling


