# Jsonify

[Datum](#datum) | [Value](#value) | [List](#list) | [Map](#map) |
[Context](#context) | [Seek](#seek)

Jsonify is a helper for writing JSON files in a way that the typical
[`json`](https://pkg.go.dev/encoding/json) package can not easily do.
Specifically writing specific datum based on a context and seeking
within the data.

See [JSON](https://www.json.org/json-en.html) for information for
how datum will be formatted.

## Datum

Datum is one part in the tree of data. It can be a singular value,
a list of datum, a key/value map with string keys and datums.

### Value

The value is a leaf node in the JSON tree.
A value can be a number (e.g. int, int64, uint8, float32),
a boolean, a string, or a null value.

Examples: `null`, `true`, `"Hello World"`, `345`, `1.02`

### List

The list is a linear collection of zero or more datum in a given order.
A list may have a mix of any datum in it.

Examples: `[]`, `[ 12, 23, 8 ]`, `[ 1, 2, [3, 4], {}, null, "cat" ]`

### Map

The map is a collection of key/value pairs where the key is a string and
the value is any datum. A map may have a mix of any datum in it but must
have unique keys. A map is typically written with keys ordered in ascending
ascii order.

Examples: `{}`, `{ "a": 12, "b": 5, "x": 8 }`, `{ "1": [], "&": {}, "cat": 42 }`

## Context

Jsonify allows a context to be passed from parents into children while
the data is being generated. To use the context add the following to a
Go object. The context allows children to be outputted in a specific
way based on the parent.

```Go
func (e Example) ToJson(ctx *Context) Datum {
    // ...
}
```

The makes it possible to write a complex graph of data which can't be
outputted to JSON by normal. The context fixes this by making it possible
to output a node graph, X, in detail but any nodes connected to X, can
be output as some identifier, name, reference, or index. Doing this
can convert the graph into a tree that represents the graph whilst being
able to be represented in JSON.

The context can also be used to control what information is being outputted
to the JSON file. This allows some additional debugging data to be added
or removed easily as well as making it possible to have multiple output modes.

## Seek

Seek uses a set of strings to specify a path of data.
Each step in the path may be one of the following:

- **counter**: `#` will return the number of datum in a datum.
  - For a list this will return the number of elements in the list.
  - For a map this will return the number of key/value pairs.
  - This is the only step allowed for a value and will always return 1.

- **index**: Any integer or string containing an integer will return
    one element at the given index in a list or panic if out-of-bounds.

- **key match**: Any text that doesn't match a different step
    will return the value with a matching key in a map.
    If no key is found this will panic. If the key is a number, starts with
    a `!` or `~`, or contains a `=` or `..`, then the test will have to be
    quoted again, e.g. `"\"!cat\""`.

- **not key match**: Step starting with `!` will return a map with all of
    the key/values pairs that do not match the given key, e.g. `!cat`.
    If the text after the `!` is quoted, the text will be unquoted.

- **regex key match**: Step starting with `~` will return a map with all of
    the key/value pairs that have keys matching the regular expression
    following the `~`.
    If the text after the `~` is quoted, the text will be unquoted.

- **no regex key match**: Step starting with `!~` will return a map with all
    the key/value pairs that have keys that do not match the regular expression.
    If the text after the `!~` is quoted, the text will be unquoted.

- **sub-value match**: Step starts with a word followed by a `=` then followed
    by text. The first word is a key and the second text is the value match.
    The value may be quoted, be negated (`=!`), be a regular expression (`=~`),
    or be a negated regular expression (`=!~`).
  - For a list, this checks every element in the list and create a list of
    all matching elements. If the element is a map, the given key is used
    to get the value and the value must be a value datum that matches
    the given value.
    - For example a list of people objects can be reduced to a smaller list
      of people with the name of "Bill" or "Jill" with `name=~^[BJ]ill$`.
  - For a map, this checks every key/value pair, X, in the map and creates
    a map of all key/value pairs where the value of X is a map that matches.
    Each value of X that is a map, gets the value for the given key and
    checks if that value matches the given value.
    - For example a map of people keyed with an identifier can be reduced to
      a map of identifiers with matching people not named "Bill" with
      `name=!"Bill"`.

- **range**: Is an optional start index followed by `..` and an optional end
    index. Both the start and end are inclusive.
    If both are given, e.g. `12..56`, the items from index 12 up to and
    including index 56 are used. If the start isn't given, e.g. `..56`, the
    start index is 0. If the end isn't given, e.g. `12..`, the end index
    is the size of the map or list. If neither are given, e.g. `..`, the whole
    list is used.
  - For a list, this gets the sub-list of the given range as a list.
    - For example a list of people objects can be reduced to a list of names
      with `[ "..", "name" ]`.
  - For a map, this gets the i'th key/value pair for the given range as a map.
    The key/value pairs are ordered by the key.
