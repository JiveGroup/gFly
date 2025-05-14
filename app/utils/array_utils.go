package utils

import (
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

// ArrContains checks if a slice contains a specific element
func ArrContains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// ArrFilter returns a new slice containing only the elements that satisfy the predicate function
func ArrFilter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// ArrMap applies a function to each element in a slice and returns a new slice with the results
func ArrMap[T any, R any](slice []T, mapFunc func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapFunc(item)
	}
	return result
}

// ArrFind returns the first element in the slice that satisfies the predicate function
// and a boolean indicating whether such an element was found
func ArrFind[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// ArrFindIndex returns the index of the first element in the slice that satisfies the predicate function
// and a boolean indicating whether such an element was found
func ArrFindIndex[T any](slice []T, predicate func(T) bool) (int, bool) {
	for i, item := range slice {
		if predicate(item) {
			return i, true
		}
	}
	return -1, false
}

// ArrUnique returns a new slice with duplicate elements removed
func ArrUnique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// ArrShuffle randomly reorders the elements in the slice
func ArrShuffle[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

// ArrChunk splits a slice into chunks of the specified size
func ArrChunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	chunks := make([][]T, 0, (len(slice)+size-1)/size)

	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// ArrSortedCopy returns a sorted copy of the slice
func ArrSortedCopy[T any](slice []T, less func(i, j T) bool) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return result
}

// ArrReduce applies a reducer function to each element in a slice, resulting in a single output value
func ArrReduce[T any, R any](slice []T, initialValue R, reducer func(acc R, item T) R) R {
	result := initialValue
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// ArrJoin concatenates the elements of a slice into a single string with the specified separator
func ArrJoin[T any](slice []T, separator string, toString func(T) string) string {
	if len(slice) == 0 {
		return ""
	}

	result := toString(slice[0])
	for i := 1; i < len(slice); i++ {
		result += separator + toString(slice[i])
	}
	return result
}

// ArrIntersection returns a slice containing elements that exist in all the provided slices
func ArrIntersection[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	// Use the first slice as a starting point
	elementCounts := make(map[T]int)
	for _, item := range slices[0] {
		elementCounts[item] = 1
	}

	// Count occurrences in other slices
	for i := 1; i < len(slices); i++ {
		seen := make(map[T]bool)
		for _, item := range slices[i] {
			if _, ok := elementCounts[item]; ok && !seen[item] {
				elementCounts[item]++
				seen[item] = true
			}
		}
	}

	// Find elements that appear in all slices
	result := make([]T, 0)
	for item, count := range elementCounts {
		if count == len(slices) {
			result = append(result, item)
		}
	}

	return result
}

// ArrUnion returns a slice containing unique elements from all the provided slices
func ArrUnion[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	seen := make(map[T]struct{})
	result := make([]T, 0)

	for _, slice := range slices {
		for _, item := range slice {
			if _, ok := seen[item]; !ok {
				seen[item] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}

// ArrDifference returns elements in the first slice that are not in any of the other slices
func ArrDifference[T comparable](slice []T, others ...[]T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	// Create a map of all elements in other slices
	exclude := make(map[T]struct{})
	for _, other := range others {
		for _, item := range other {
			exclude[item] = struct{}{}
		}
	}

	// Keep elements from the first slice that are not in the exclude map
	result := make([]T, 0)
	for _, item := range slice {
		if _, ok := exclude[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

// ArrGroupBy groups elements in a slice by a key generated from each element
func ArrGroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// ArrFlatten converts a slice of slices into a single slice containing all elements
func ArrFlatten[T any](slices [][]T) []T {
	totalLen := 0
	for _, slice := range slices {
		totalLen += len(slice)
	}

	result := make([]T, 0, totalLen)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// ArrReverse returns a new slice with the elements in reverse order
func ArrReverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, j := 0, len(slice)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = slice[j]
	}
	return result
}

// ArrAccessible checks if the given value can be accessed as an array/slice
func ArrAccessible(value interface{}) bool {
	if value == nil {
		return false
	}

	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}

// ArrAdd adds a key/value pair to an array/map if the key doesn't exist
func ArrAdd(array map[string]interface{}, key string, value interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range array {
		result[k] = v
	}

	if _, exists := result[key]; !exists {
		result[key] = value
	}

	return result
}

// ArrCollapse collapses a slice of slices into a single slice
func ArrCollapse(arrays [][]interface{}) []interface{} {
	totalLen := 0
	for _, arr := range arrays {
		totalLen += len(arr)
	}

	result := make([]interface{}, 0, totalLen)
	for _, arr := range arrays {
		result = append(result, arr...)
	}

	return result
}

// ArrCrossJoin cross joins the given arrays, returning a cartesian product with all possible permutations
func ArrCrossJoin[T any](arrays ...[]T) [][]T {
	if len(arrays) == 0 {
		return [][]T{}
	}

	if len(arrays) == 1 {
		result := make([][]T, len(arrays[0]))
		for i, item := range arrays[0] {
			result[i] = []T{item}
		}
		return result
	}

	// Get the cartesian product of all but the first array
	subResult := ArrCrossJoin[T](arrays[1:]...)

	// Combine the first array with the sub-result
	result := make([][]T, 0, len(arrays[0])*len(subResult))
	for _, item := range arrays[0] {
		for _, subItem := range subResult {
			newItem := make([]T, 1+len(subItem))
			newItem[0] = item
			copy(newItem[1:], subItem)
			result = append(result, newItem)
		}
	}

	return result
}

// ArrDivide returns two slices, one containing the keys, and the other containing the values of the original map
func ArrDivide(array map[string]interface{}) ([]string, []interface{}) {
	keys := make([]string, 0, len(array))
	values := make([]interface{}, 0, len(array))

	for k, v := range array {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// ArrDot flattens a multi-dimensional map into a single level map with "dot" notation
func ArrDot(array map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	dotRecursive(array, result, "")
	return result
}

// dotRecursive is a helper function for Dot
func dotRecursive(array map[string]interface{}, result map[string]interface{}, prepend string) {
	for key, value := range array {
		if prepend != "" {
			key = prepend + "." + key
		}

		if subArray, ok := value.(map[string]interface{}); ok {
			dotRecursive(subArray, result, key)
		} else {
			result[key] = value
		}
	}
}

// ArrExcept returns the array with the specified keys removed
func ArrExcept(array map[string]interface{}, keys ...string) map[string]interface{} {
	result := make(map[string]interface{})

	// Create a map for faster lookup
	keysMap := make(map[string]struct{})
	for _, key := range keys {
		keysMap[key] = struct{}{}
	}

	for key, value := range array {
		if _, exists := keysMap[key]; !exists {
			result[key] = value
		}
	}

	return result
}

// ArrExists checks if the given key exists in the array
func ArrExists(array map[string]interface{}, key string) bool {
	_, exists := array[key]
	return exists
}

// ArrFirst returns the first element in an array passing a given truth test
func ArrFirst[T any](array []T, callback func(T) bool) (T, bool) {
	for _, value := range array {
		if callback(value) {
			return value, true
		}
	}

	var zero T
	return zero, false
}

// ArrFirstOrDefault returns the first element in the array, or a default if the array is empty
func ArrFirstOrDefault[T any](array []T, defaultValue T) T {
	if len(array) > 0 {
		return array[0]
	}
	return defaultValue
}

// ArrForget removes the given key/value pair from the array
func ArrForget(array map[string]interface{}, keys ...string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range array {
		result[k] = v
	}

	for _, key := range keys {
		delete(result, key)
	}

	return result
}

// ArrGet retrieves a value from an array using "dot" notation
func ArrGet(array map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if array == nil {
		return defaultValue
	}

	if key == "" {
		return array
	}

	keys := strings.Split(key, ".")
	current := array

	for i, segment := range keys {
		if i == len(keys)-1 {
			if val, exists := current[segment]; exists {
				return val
			}
			return defaultValue
		}

		if val, exists := current[segment]; exists {
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return defaultValue
			}
		} else {
			return defaultValue
		}
	}

	return defaultValue
}

// ArrHas determines if any of the keys exist in the array using "dot" notation
func ArrHas(array map[string]interface{}, keys ...string) bool {
	if len(keys) == 0 {
		return false
	}

	for _, key := range keys {
		if !hasDot(array, key) {
			return false
		}
	}

	return true
}

// hasDot is a helper function for Has
func hasDot(array map[string]interface{}, key string) bool {
	if array == nil {
		return false
	}

	if key == "" {
		return false
	}

	keys := strings.Split(key, ".")
	current := array

	for i, segment := range keys {
		if val, exists := current[segment]; exists {
			if i == len(keys)-1 {
				return true
			}

			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return false
			}
		} else {
			return false
		}
	}

	return false
}

// ArrHasAny determines if any of the keys exist in the array using "dot" notation
func ArrHasAny(array map[string]interface{}, keys ...string) bool {
	if len(keys) == 0 {
		return false
	}

	for _, key := range keys {
		if hasDot(array, key) {
			return true
		}
	}

	return false
}

// ArrIsAssoc determines if an array is associative (has string keys)
func ArrIsAssoc(array interface{}) bool {
	if array == nil {
		return false
	}

	value := reflect.ValueOf(array)
	if value.Kind() != reflect.Map {
		return false
	}

	// Check if all keys are strings
	for _, key := range value.MapKeys() {
		if key.Kind() != reflect.String {
			return false
		}
	}

	return true
}

// ArrIsList determines if an array has sequential numeric keys
func ArrIsList(array interface{}) bool {
	if array == nil {
		return false
	}

	value := reflect.ValueOf(array)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return false
	}

	return true
}

// ArrKeyBy keys the array by the given key
func ArrKeyBy[T any, K comparable](array []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range array {
		key := keyFunc(item)
		result[key] = item
	}
	return result
}

// ArrLast returns the last element in an array passing a given truth test
func ArrLast[T any](array []T, callback func(T) bool) (T, bool) {
	for i := len(array) - 1; i >= 0; i-- {
		if callback(array[i]) {
			return array[i], true
		}
	}

	var zero T
	return zero, false
}

// ArrLastOrDefault returns the last element in the array, or a default if the array is empty
func ArrLastOrDefault[T any](array []T, defaultValue T) T {
	if len(array) > 0 {
		return array[len(array)-1]
	}
	return defaultValue
}

// ArrOnly returns the array with only the specified keys
func ArrOnly(array map[string]interface{}, keys ...string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range keys {
		if value, exists := array[key]; exists {
			result[key] = value
		}
	}

	return result
}

// ArrPluck retrieves all of the values for a given key
func ArrPluck[T any, V any](array []T, key func(T) V) []V {
	result := make([]V, len(array))
	for i, item := range array {
		result[i] = key(item)
	}
	return result
}

// ArrPrepend adds an item to the beginning of an array
func ArrPrepend[T any](array []T, values ...T) []T {
	result := make([]T, len(values)+len(array))
	copy(result, values)
	copy(result[len(values):], array)
	return result
}

// ArrPull removes and returns an item from the array by key
func ArrPull[T any](array []T, index int) (T, []T) {
	if index < 0 || index >= len(array) {
		var zero T
		return zero, array
	}

	item := array[index]
	result := append(array[:index], array[index+1:]...)
	return item, result
}

// ArrQuery builds a query string from the array
func ArrQuery(array map[string]interface{}) string {
	values := url.Values{}
	for key, value := range array {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case []string:
			for _, item := range v {
				values.Add(key+"[]", item)
			}
		default:
			// Convert to string using fmt.Sprint
			values.Add(key, fmt.Sprint(v))
		}
	}

	return values.Encode()
}

// ArrRandom returns a random value from an array
func ArrRandom[T any](array []T) (T, bool) {
	if len(array) == 0 {
		var zero T
		return zero, false
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[r.Intn(len(array))], true
}

// ArrRandomOrDefault returns a random value from an array or a default value if the array is empty
func ArrRandomOrDefault[T any](array []T, defaultValue T) T {
	if len(array) == 0 {
		return defaultValue
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return array[r.Intn(len(array))]
}

// ArrSet sets a value within a nested array using "dot" notation
func ArrSet(array map[string]interface{}, key string, value interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range array {
		result[k] = v
	}

	if key == "" {
		return result
	}

	keys := strings.Split(key, ".")
	current := result

	for i, segment := range keys {
		if i == len(keys)-1 {
			current[segment] = value
			break
		}

		if val, exists := current[segment]; exists {
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				// Convert to map if it's not already
				nextMap = make(map[string]interface{})
				current[segment] = nextMap
				current = nextMap
			}
		} else {
			nextMap := make(map[string]interface{})
			current[segment] = nextMap
			current = nextMap
		}
	}

	return result
}

// ArrSortByKey sorts an array by key
func ArrSortByKey(array map[string]interface{}) map[string]interface{} {
	// Get the keys and sort them
	keys := make([]string, 0, len(array))
	for k := range array {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create a new map with the sorted keys
	result := make(map[string]interface{})
	for _, k := range keys {
		result[k] = array[k]
	}

	return result
}

// ArrSortByKeyDesc sorts an array by key in descending order
func ArrSortByKeyDesc(array map[string]interface{}) map[string]interface{} {
	// Get the keys and sort them in descending order
	keys := make([]string, 0, len(array))
	for k := range array {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	// Create a new map with the sorted keys
	result := make(map[string]interface{})
	for _, k := range keys {
		result[k] = array[k]
	}

	return result
}

// ArrSortRecursive recursively sorts an array by keys and values
func ArrSortRecursive(array interface{}) interface{} {
	switch arr := array.(type) {
	case map[string]interface{}:
		// Sort the map by keys
		keys := make([]string, 0, len(arr))
		for k := range arr {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Create a new map with the sorted keys and recursively sorted values
		result := make(map[string]interface{})
		for _, k := range keys {
			result[k] = ArrSortRecursive(arr[k])
		}

		return result

	case []interface{}:
		// Create a new slice with recursively sorted values
		result := make([]interface{}, len(arr))
		for i, v := range arr {
			result[i] = ArrSortRecursive(v)
		}

		return result

	default:
		// Return the value as is
		return array
	}
}

// ArrUndot expands a flattened array with "dot" notation back into a multi-dimensional array
func ArrUndot(array map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range array {
		parts := strings.Split(key, ".")

		// Reference to the current level in the result
		current := result

		// Traverse the parts of the key
		for i, part := range parts {
			// If this is the last part, set the value
			if i == len(parts)-1 {
				current[part] = value
				continue
			}

			// If the next level doesn't exist, create it
			if _, exists := current[part]; !exists {
				current[part] = make(map[string]interface{})
			}

			// Move to the next level
			current = current[part].(map[string]interface{})
		}
	}

	return result
}

// ArrWhereNotNull filters the array using the given callback, removing null values
func ArrWhereNotNull[T any](array []T) []T {
	result := make([]T, 0)

	for _, item := range array {
		// Check if the item is nil
		if !isNil(item) {
			result = append(result, item)
		}
	}

	return result
}

// isNil checks if a value is nil
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	kind := v.Kind()

	// Check for nil pointers, interfaces, maps, slices, and channels
	return (kind == reflect.Ptr || kind == reflect.Interface ||
		kind == reflect.Map || kind == reflect.Slice ||
		kind == reflect.Chan) && v.IsNil()
}

// ArrWrap wraps the given value in an array if it's not already an array
func ArrWrap(value interface{}) []interface{} {
	if value == nil {
		return []interface{}{}
	}

	v := reflect.ValueOf(value)

	// If it's already a slice or array, convert it to []interface{}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result
	}

	// Otherwise, wrap it in a slice
	return []interface{}{value}
}

// MapMergeMaps merges multiple maps into a new map
// If there are duplicate keys, the value from the later map will overwrite the earlier one
func MapMergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// MapKeys returns a slice containing all the keys in the map
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapValues returns a slice containing all the values in the map
func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// MapFindKey returns the first key that maps to the specified value and a boolean indicating whether such a key was found
func MapFindKey[K comparable, V comparable](m map[K]V, value V) (K, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// MapFilterMap returns a new map containing only the key-value pairs that satisfy the predicate function
func MapFilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapMapValues applies a function to each value in the map and returns a new map with the transformed values
func MapMapValues[K comparable, V any, R any](m map[K]V, mapFunc func(V) R) map[K]R {
	result := make(map[K]R, len(m))
	for k, v := range m {
		result[k] = mapFunc(v)
	}
	return result
}

// MapInvertMap creates a new map by swapping the keys and values of the original map
// Note: If multiple keys map to the same value in the original map, only one key-value pair will be in the result
func MapInvertMap[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// MapGetOrDefault returns the value for the given key if it exists, otherwise returns the default value
func MapGetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

// MapGetOrInsert returns the value for the given key if it exists, otherwise inserts and returns the default value
func MapGetOrInsert[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	m[key] = defaultValue
	return defaultValue
}

// MapToSlice converts a map to a slice of key-value pairs
func MapToSlice[K comparable, V any](m map[K]V) []struct {
	Key   K
	Value V
} {
	result := make([]struct {
		Key   K
		Value V
	}, 0, len(m))

	for k, v := range m {
		result = append(result, struct {
			Key   K
			Value V
		}{k, v})
	}

	return result
}

// MapSliceToMap converts a slice of key-value pairs to a map
func MapSliceToMap[K comparable, V any](slice []struct {
	Key   K
	Value V
}) map[K]V {
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		result[item.Key] = item.Value
	}
	return result
}

// MapEqualMaps checks if two maps contain the same key-value pairs
func MapEqualMaps[K, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}

	return true
}

// MapDiffMaps returns the keys that are different between two maps
// The result contains three maps:
// - added: keys in m2 that are not in m1
// - removed: keys in m1 that are not in m2
// - changed: keys that exist in both maps but have different values
func MapDiffMaps[K comparable, V comparable](m1, m2 map[K]V) (added, removed, changed map[K]V) {
	added = make(map[K]V)
	removed = make(map[K]V)
	changed = make(map[K]V)

	// Find keys in m1 that are not in m2 or have different values
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok {
			removed[k] = v1
		} else if v1 != v2 {
			changed[k] = v2
		}
	}

	// Find keys in m2 that are not in m1
	for k, v2 := range m2 {
		if _, ok := m1[k]; !ok {
			added[k] = v2
		}
	}

	return added, removed, changed
}

// SetContains checks if a set contains an element
func SetContains[T comparable](set map[T]struct{}, item T) bool {
	_, ok := set[item]
	return ok
}

// SetToSlice converts a set to a slice
func SetToSlice[T comparable](set map[T]struct{}) []T {
	result := make([]T, 0, len(set))
	for item := range set {
		result = append(result, item)
	}
	return result
}

// SliceToSet converts a slice to a set
func SliceToSet[T comparable](slice []T) map[T]struct{} {
	result := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		result[item] = struct{}{}
	}
	return result
}

// SetUnion returns a new set containing elements from both sets
func SetUnion[T comparable](set1, set2 map[T]struct{}) map[T]struct{} {
	result := make(map[T]struct{}, len(set1)+len(set2))
	for item := range set1 {
		result[item] = struct{}{}
	}
	for item := range set2 {
		result[item] = struct{}{}
	}
	return result
}

// SetIntersection returns a new set containing elements that exist in both sets
func SetIntersection[T comparable](set1, set2 map[T]struct{}) map[T]struct{} {
	result := make(map[T]struct{})

	// Use the smaller set for iteration
	if len(set1) > len(set2) {
		set1, set2 = set2, set1
	}

	for item := range set1 {
		if _, ok := set2[item]; ok {
			result[item] = struct{}{}
		}
	}
	return result
}

// SetDifference returns a new set containing elements in set1 that are not in set2
func SetDifference[T comparable](set1, set2 map[T]struct{}) map[T]struct{} {
	result := make(map[T]struct{})
	for item := range set1 {
		if _, ok := set2[item]; !ok {
			result[item] = struct{}{}
		}
	}
	return result
}
