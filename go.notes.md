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
