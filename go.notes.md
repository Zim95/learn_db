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