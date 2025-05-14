# Utilities

This directory contains utility functions and helper classes that provide common functionality used throughout the application.

## Table of Contents

- [Array Utilities](#array-utilities)
- [Collection Utilities](#collection-utilities)
- [HTTP Utilities](#http-utilities)
- [Number Utilities](#number-utilities)
- [String Utilities](#string-utilities)
- [Transform Utilities](#transform-utilities)
- [Validation Utilities](#validation-utilities)

## Array Utilities

Array utilities provide functions for working with arrays, slices, maps, and sets in Go.

### ArrContains

```go
func ArrContains[T comparable](slice []T, element T) bool
```

Checks if a slice contains a specific element.

**Parameters:**
- `slice []T`: The slice to search in
- `element T`: The element to search for

**Returns:**
- `bool`: True if the element is found, false otherwise

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
contains := utils.ArrContains(numbers, 3) // true
notContains := utils.ArrContains(numbers, 6) // false
```

### ArrFilter

```go
func ArrFilter[T any](slice []T, predicate func(T) bool) []T
```

Returns a new slice containing only the elements that satisfy the predicate function.

**Parameters:**
- `slice []T`: The input slice
- `predicate func(T) bool`: A function that returns true for elements to include

**Returns:**
- `[]T`: A new slice with filtered elements

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
evenNumbers := utils.ArrFilter(numbers, func(n int) bool {
    return n%2 == 0
}) // [2, 4]
```

### ArrMap

```go
func ArrMap[T any, R any](slice []T, mapFunc func(T) R) []R
```

Applies a function to each element in a slice and returns a new slice with the results.

**Parameters:**
- `slice []T`: The input slice
- `mapFunc func(T) R`: A function to transform each element

**Returns:**
- `[]R`: A new slice with transformed elements

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
doubled := utils.ArrMap(numbers, func(n int) int {
    return n * 2
}) // [2, 4, 6, 8, 10]
```

### ArrFind

```go
func ArrFind[T any](slice []T, predicate func(T) bool) (T, bool)
```

Returns the first element in the slice that satisfies the predicate function and a boolean indicating whether such an element was found.

**Parameters:**
- `slice []T`: The input slice
- `predicate func(T) bool`: A function that returns true for the element to find

**Returns:**
- `T`: The found element or zero value if not found
- `bool`: True if an element was found, false otherwise

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
firstEven, found := utils.ArrFind(numbers, func(n int) bool {
    return n%2 == 0
}) // 2, true

firstNegative, found := utils.ArrFind(numbers, func(n int) bool {
    return n < 0
}) // 0, false
```

### ArrFindIndex

```go
func ArrFindIndex[T any](slice []T, predicate func(T) bool) (int, bool)
```

Returns the index of the first element in the slice that satisfies the predicate function and a boolean indicating whether such an element was found.

**Parameters:**
- `slice []T`: The input slice
- `predicate func(T) bool`: A function that returns true for the element to find

**Returns:**
- `int`: The index of the found element or -1 if not found
- `bool`: True if an element was found, false otherwise

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
evenIndex, found := utils.ArrFindIndex(numbers, func(n int) bool {
    return n%2 == 0
}) // 1, true

negativeIndex, found := utils.ArrFindIndex(numbers, func(n int) bool {
    return n < 0
}) // -1, false
```

### ArrUnique

```go
func ArrUnique[T comparable](slice []T) []T
```

Returns a new slice with duplicate elements removed.

**Parameters:**
- `slice []T`: The input slice

**Returns:**
- `[]T`: A new slice with unique elements

**Example:**
```go
numbers := []int{1, 2, 2, 3, 3, 3, 4, 5}
unique := utils.ArrUnique(numbers) // [1, 2, 3, 4, 5]
```

### ArrShuffle

```go
func ArrShuffle[T any](slice []T) []T
```

Randomly reorders the elements in the slice.

**Parameters:**
- `slice []T`: The input slice

**Returns:**
- `[]T`: A new slice with shuffled elements

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
shuffled := utils.ArrShuffle(numbers) // e.g., [3, 1, 5, 2, 4]
```

### ArrChunk

```go
func ArrChunk[T any](slice []T, size int) [][]T
```

Splits a slice into chunks of the specified size.

**Parameters:**
- `slice []T`: The input slice
- `size int`: The size of each chunk

**Returns:**
- `[][]T`: A slice of slices, each containing at most `size` elements

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
chunks := utils.ArrChunk(numbers, 3) // [[1, 2, 3], [4, 5, 6], [7, 8]]
```

### MapMergeMaps

```go
func MapMergeMaps[K comparable, V any](maps ...map[K]V) map[K]V
```

Merges multiple maps into a new map. If there are duplicate keys, the value from the later map will overwrite the earlier one.

**Parameters:**
- `maps ...map[K]V`: The maps to merge

**Returns:**
- `map[K]V`: A new map containing all key-value pairs from the input maps

**Example:**
```go
map1 := map[string]int{"a": 1, "b": 2}
map2 := map[string]int{"b": 3, "c": 4}
merged := utils.MapMergeMaps(map1, map2) // {"a": 1, "b": 3, "c": 4}
```

## Collection Utilities

Collection utilities provide functions for working with collections (slices and maps) in Go.

### ColAll

```go
func ColAll[T any](collection []T) []T
```

Returns all of the items in the collection.

**Parameters:**
- `collection []T`: The input collection

**Returns:**
- `[]T`: The same collection

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
all := utils.ColAll(numbers) // [1, 2, 3, 4, 5]
```

### ColAvg

```go
func ColAvg[T any](collection []T, valueFunc func(T) float64) float64
```

Returns the average value of a given key.

**Parameters:**
- `collection []T`: The input collection
- `valueFunc func(T) float64`: A function that extracts a numeric value from each element

**Returns:**
- `float64`: The average of the values

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
avg := utils.ColAvg(numbers, func(n int) float64 {
    return float64(n)
}) // 3.0
```

### ColContains

```go
func ColContains[T comparable](collection []T, item T) bool
```

Determines whether the collection contains a given item.

**Parameters:**
- `collection []T`: The input collection
- `item T`: The item to search for

**Returns:**
- `bool`: True if the item is found, false otherwise

**Example:**
```go
numbers := []int{1, 2, 3, 4, 5}
contains := utils.ColContains(numbers, 3) // true
notContains := utils.ColContains(numbers, 6) // false
```

## HTTP Utilities

HTTP utilities provide functions for making HTTP requests and handling HTTP-related operations.

### HttpBuildURL

```go
func HttpBuildURL(baseURL string, queryParams map[string]string) (string, error)
```

Constructs a URL from a base URL and query parameters.

**Parameters:**
- `baseURL string`: The base URL
- `queryParams map[string]string`: Query parameters to add to the URL

**Returns:**
- `string`: The constructed URL
- `error`: An error if the URL is invalid

**Example:**
```go
url, err := utils.HttpBuildURL("https://example.com", map[string]string{
    "q": "search term",
    "page": "1",
}) // "https://example.com?q=search+term&page=1"
```

### HttpGetJSON

```go
func HttpGetJSON(url string, target interface{}, headers map[string]string) error
```

Performs a GET request and unmarshals the JSON response into the provided interface.

**Parameters:**
- `url string`: The URL to request
- `target interface{}`: The target object to unmarshal the response into
- `headers map[string]string`: HTTP headers to include in the request

**Returns:**
- `error`: An error if the request fails or the response cannot be unmarshaled

**Example:**
```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var user User
err := utils.HttpGetJSON("https://api.example.com/users/1", &user, nil)
```

## Number Utilities

Number utilities provide functions for working with numbers in Go.

### NumCeiling

```go
func NumCeiling(number float64, precision ...int) float64
```

Rounds a number up to the nearest integer or precision. If precision is provided, the number is rounded up to that decimal place.

**Parameters:**
- `number float64`: The number to round
- `precision ...int`: Optional precision (decimal places)

**Returns:**
- `float64`: The rounded number

**Example:**
```go
ceiling := utils.NumCeiling(3.14159) // 4.0
ceilingWithPrecision := utils.NumCeiling(3.14159, 2) // 3.15
```

### NumFloor

```go
func NumFloor(number float64, precision ...int) float64
```

Rounds a number down to the nearest integer or precision. If precision is provided, the number is rounded down to that decimal place.

**Parameters:**
- `number float64`: The number to round
- `precision ...int`: Optional precision (decimal places)

**Returns:**
- `float64`: The rounded number

**Example:**
```go
floor := utils.NumFloor(3.14159) // 3.0
floorWithPrecision := utils.NumFloor(3.14159, 2) // 3.14
```

### NumFormat

```go
func NumFormat(number float64, decimals int, decimalSeparator, thousandsSeparator string) string
```

Formats a number with grouped thousands. The default separator is a comma, but a custom separator can be provided.

**Parameters:**
- `number float64`: The number to format
- `decimals int`: The number of decimal places
- `decimalSeparator string`: The character to use as decimal separator
- `thousandsSeparator string`: The character to use as thousands separator

**Returns:**
- `string`: The formatted number as a string

**Example:**
```go
formatted := utils.NumFormat(1234567.89, 2, ".", ",") // "1,234,567.89"
```

## String Utilities

String utilities provide functions for working with strings in Go.

### StrTruncateString

```go
func StrTruncateString(s string, maxLength int) string
```

Truncates a string to the specified length and adds an ellipsis if truncated.

**Parameters:**
- `s string`: The input string
- `maxLength int`: The maximum length

**Returns:**
- `string`: The truncated string

**Example:**
```go
truncated := utils.StrTruncateString("Hello, world!", 5) // "Hello..."
```

### StrSlugifyString

```go
func StrSlugifyString(s string) string
```

Converts a string to a URL-friendly slug.

**Parameters:**
- `s string`: The input string

**Returns:**
- `string`: The slugified string

**Example:**
```go
slug := utils.StrSlugifyString("Hello, World!") // "hello-world"
```

### StrGenerateRandomString

```go
func StrGenerateRandomString(length int) (string, error)
```

Generates a random string of the specified length.

**Parameters:**
- `length int`: The length of the random string

**Returns:**
- `string`: The generated random string
- `error`: An error if random generation fails

**Example:**
```go
randomString, err := utils.StrGenerateRandomString(10) // e.g., "a1b2c3d4e5"
```

## Transform Utilities

Transform utilities provide functions for transforming data structures in Go.

### TransformList

```go
func TransformList[T any, R any](records []T, transformerFn func(T) R) []R
```

Takes a list of records and applies a transformer function to each record, returning a slice of the transformed records.

**Parameters:**
- `records []T`: The input records
- `transformerFn func(T) R`: A function to transform each record

**Returns:**
- `[]R`: A slice of transformed records

**Example:**
```go
type User struct {
    ID   int
    Name string
}

type UserDTO struct {
    ID   int
    Name string
}

users := []User{{1, "Alice"}, {2, "Bob"}}
userDTOs := utils.TransformList(users, func(user User) UserDTO {
    return UserDTO{ID: user.ID, Name: user.Name}
})
```

### TransformMap

```go
func TransformMap[K comparable, V any, R any](m map[K]V, transformerFn func(V) R) map[K]R
```

Transforms a map of one type to a map of another type using a transformer function.

**Parameters:**
- `m map[K]V`: The input map
- `transformerFn func(V) R`: A function to transform each value

**Returns:**
- `map[K]R`: A map with the same keys but transformed values

**Example:**
```go
userAges := map[string]int{"Alice": 30, "Bob": 25}
userAgeStrings := utils.TransformMap(userAges, func(age int) string {
    return fmt.Sprintf("%d years old", age)
}) // {"Alice": "30 years old", "Bob": "25 years old"}
```

## Validation Utilities

Validation utilities provide functions for validating data in Go.

### Validate

```go
func Validate(structData any, msgForTagFunc ...validation.MsgForTagFunc) *response.Error
```

Performs data input checking.

**Parameters:**
- `structData any`: The struct to validate
- `msgForTagFunc ...validation.MsgForTagFunc`: Optional functions to customize error messages

**Returns:**
- `*response.Error`: An error if validation fails, nil otherwise

**Example:**
```go
type UserInput struct {
    Name  string `validate:"required"`
    Email string `validate:"required,email"`
}

input := UserInput{Name: "Alice", Email: "invalid-email"}
err := utils.Validate(input)
if err != nil {
    // Handle validation error
}
```

### IsValidEmail

```go
func IsValidEmail(email string) bool
```

Checks if a string is a valid email address.

**Parameters:**
- `email string`: The email address to validate

**Returns:**
- `bool`: True if the email is valid, false otherwise

**Example:**
```go
validEmail := utils.IsValidEmail("user@example.com") // true
invalidEmail := utils.IsValidEmail("invalid-email") // false
```

### IsValidURL

```go
func IsValidURL(rawURL string) bool
```

Checks if a string is a valid URL.

**Parameters:**
- `rawURL string`: The URL to validate

**Returns:**
- `bool`: True if the URL is valid, false otherwise

**Example:**
```go
validURL := utils.IsValidURL("https://example.com") // true
invalidURL := utils.IsValidURL("not-a-url") // false
```

## Best Practices

- Keep utility functions simple and focused on a single task
- Use descriptive function names that clearly indicate what the function does
- Document parameters, return values, and any side effects
- Write comprehensive tests for utility functions
- Avoid dependencies on application-specific code in utilities
- Consider performance implications for frequently used utilities
