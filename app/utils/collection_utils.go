package utils

import (
	"math/rand"
	"sort"
	"time"
)

// ColAll returns all of the items in the collection
func ColAll[T any](collection []T) []T {
	return collection
}

// ColAvg returns the average value of a given key
func ColAvg[T any](collection []T, valueFunc func(T) float64) float64 {
	if len(collection) == 0 {
		return 0
	}

	sum := 0.0
	for _, item := range collection {
		sum += valueFunc(item)
	}

	return sum / float64(len(collection))
}

// ColChunk breaks the collection into multiple, smaller collections of a given size
func ColChunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	chunks := make([][]T, 0, (len(collection)+size-1)/size)

	for i := 0; i < len(collection); i += size {
		end := i + size
		if end > len(collection) {
			end = len(collection)
		}
		chunks = append(chunks, collection[i:end])
	}

	return chunks
}

// ColCollapse collapses a collection of arrays into a single, flat collection
func ColCollapse[T any](collection [][]T) []T {
	totalLen := 0
	for _, slice := range collection {
		totalLen += len(slice)
	}

	result := make([]T, 0, totalLen)
	for _, slice := range collection {
		result = append(result, slice...)
	}
	return result
}

// ColContains determines whether the collection contains a given item
func ColContains[T comparable](collection []T, item T) bool {
	for _, value := range collection {
		if value == item {
			return true
		}
	}
	return false
}

// ColContainsFunc determines whether the collection contains an item with the given predicate
func ColContainsFunc[T any](collection []T, predicate func(T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return true
		}
	}
	return false
}

// ColCount returns the total number of items in the collection
func ColCount[T any](collection []T) int {
	return len(collection)
}

// ColCountBy counts the occurrences of values in the collection
func ColCountBy[T any, K comparable](collection []T, keyFunc func(T) K) map[K]int {
	result := make(map[K]int)
	for _, item := range collection {
		key := keyFunc(item)
		result[key]++
	}
	return result
}

// ColCrossJoin cross joins the collection with the given arrays or collections
func ColCrossJoin[T any](collection []T, arrays ...[]T) [][]T {
	if len(collection) == 0 {
		return [][]T{}
	}

	// Start with the original collection as single-item arrays
	result := make([][]T, len(collection))
	for i, item := range collection {
		result[i] = []T{item}
	}

	// Cross join with each additional array
	for _, array := range arrays {
		newResult := make([][]T, 0, len(result)*len(array))
		for _, item := range result {
			for _, value := range array {
				newItem := make([]T, len(item)+1)
				copy(newItem, item)
				newItem[len(item)] = value
				newResult = append(newResult, newItem)
			}
		}
		result = newResult
	}

	return result
}

// ColDiff compares the collection against another collection or array
func ColDiff[T comparable](collection []T, items []T) []T {
	// Create a map for faster lookup
	itemMap := make(map[T]struct{})
	for _, item := range items {
		itemMap[item] = struct{}{}
	}

	// Keep elements from the collection that are not in items
	result := make([]T, 0)
	for _, item := range collection {
		if _, exists := itemMap[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}

// ColDiffAssoc compares the collection against another collection or array based on its keys and values
func ColDiffAssoc[K comparable, V comparable](collection map[K]V, items map[K]V) map[K]V {
	result := make(map[K]V)

	for key, value := range collection {
		if itemValue, exists := items[key]; !exists || itemValue != value {
			result[key] = value
		}
	}

	return result
}

// ColDiffKeys compares the collection against another collection or array based on its keys
func ColDiffKeys[K comparable, V any](collection map[K]V, items map[K]V) map[K]V {
	result := make(map[K]V)

	for key, value := range collection {
		if _, exists := items[key]; !exists {
			result[key] = value
		}
	}

	return result
}

// ColEach iterates over the collection and passes each item to the given callback
func ColEach[T any](collection []T, callback func(T, int) bool) {
	for i, item := range collection {
		if !callback(item, i) {
			break
		}
	}
}

// ColEvery determines if all elements of the collection pass a given truth test
func ColEvery[T any](collection []T, predicate func(T) bool) bool {
	for _, item := range collection {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// ColExcept returns all items in the collection except for those with the specified keys
func ColExcept[K comparable, V any](collection map[K]V, keys []K) map[K]V {
	result := make(map[K]V)

	// Create a map for faster lookup
	keysMap := make(map[K]struct{})
	for _, key := range keys {
		keysMap[key] = struct{}{}
	}

	for key, value := range collection {
		if _, exists := keysMap[key]; !exists {
			result[key] = value
		}
	}

	return result
}

// ColFilter filters the collection using the given callback
func ColFilter[T any](collection []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range collection {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// ColFirst returns the first element in the collection that passes a given truth test
func ColFirst[T any](collection []T, predicate func(T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// ColFirstOrDefault returns the first element in the collection or a default value if the collection is empty
func ColFirstOrDefault[T any](collection []T, defaultValue T) T {
	if len(collection) > 0 {
		return collection[0]
	}
	return defaultValue
}

// ColFlatMap iterates through the collection and passes each value to the given callback
func ColFlatMap[T any, R any](collection []T, callback func(T) []R) []R {
	result := make([]R, 0)
	for _, item := range collection {
		result = append(result, callback(item)...)
	}
	return result
}

// ColFlatten flattens a multi-dimensional collection into a single dimension
func ColFlatten[T any](collection [][]T) []T {
	totalLen := 0
	for _, slice := range collection {
		totalLen += len(slice)
	}

	result := make([]T, 0, totalLen)
	for _, slice := range collection {
		result = append(result, slice...)
	}
	return result
}

// ColFlip swaps the collection's keys with their corresponding values
func ColFlip[K comparable, V comparable](collection map[K]V) map[V]K {
	result := make(map[V]K)
	for key, value := range collection {
		result[value] = key
	}
	return result
}

// ColForget removes an item from the collection by its key
func ColForget[K comparable, V any](collection map[K]V, keys ...K) map[K]V {
	result := make(map[K]V)
	for k, v := range collection {
		result[k] = v
	}

	for _, key := range keys {
		delete(result, key)
	}

	return result
}

// ColGet retrieves an item from the collection by its key
func ColGet[K comparable, V any](collection map[K]V, key K, defaultValue V) V {
	if value, exists := collection[key]; exists {
		return value
	}
	return defaultValue
}

// ColGroupBy groups the collection's items by a given key
func ColGroupBy[T any, K comparable](collection []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range collection {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// ColHas determines if a given key exists in the collection
func ColHas[K comparable, V any](collection map[K]V, key K) bool {
	_, exists := collection[key]
	return exists
}

// ColImplode joins the items in a collection
func ColImplode[T any](collection []T, separator string, toString func(T) string) string {
	if len(collection) == 0 {
		return ""
	}

	result := toString(collection[0])
	for i := 1; i < len(collection); i++ {
		result += separator + toString(collection[i])
	}
	return result
}

// ColIntersect removes any values from the original collection that are not present in the given array or collection
func ColIntersect[T comparable](collection []T, items []T) []T {
	// Create a map for faster lookup
	itemMap := make(map[T]struct{})
	for _, item := range items {
		itemMap[item] = struct{}{}
	}

	// Keep elements from the collection that are in items
	result := make([]T, 0)
	for _, item := range collection {
		if _, exists := itemMap[item]; exists {
			result = append(result, item)
		}
	}

	return result
}

// ColIntersectByKeys removes any keys from the original collection that are not present in the given array or collection
func ColIntersectByKeys[K comparable, V any](collection map[K]V, keys []K) map[K]V {
	result := make(map[K]V)

	// Create a map for faster lookup
	keysMap := make(map[K]struct{})
	for _, key := range keys {
		keysMap[key] = struct{}{}
	}

	for key, value := range collection {
		if _, exists := keysMap[key]; exists {
			result[key] = value
		}
	}

	return result
}

// ColIsEmpty determines if the collection is empty
func ColIsEmpty[T any](collection []T) bool {
	return len(collection) == 0
}

// ColIsNotEmpty determines if the collection is not empty
func ColIsNotEmpty[T any](collection []T) bool {
	return len(collection) > 0
}

// ColKeyBy keys the collection by the given key
func ColKeyBy[T any, K comparable](collection []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range collection {
		key := keyFunc(item)
		result[key] = item
	}
	return result
}

// ColKeys returns all of the collection's keys
func ColKeys[K comparable, V any](collection map[K]V) []K {
	keys := make([]K, 0, len(collection))
	for key := range collection {
		keys = append(keys, key)
	}
	return keys
}

// ColLast returns the last element in the collection that passes a given truth test
func ColLast[T any](collection []T, predicate func(T) bool) (T, bool) {
	for i := len(collection) - 1; i >= 0; i-- {
		if predicate(collection[i]) {
			return collection[i], true
		}
	}
	var zero T
	return zero, false
}

// ColLastOrDefault returns the last element in the collection or a default value if the collection is empty
func ColLastOrDefault[T any](collection []T, defaultValue T) T {
	if len(collection) > 0 {
		return collection[len(collection)-1]
	}
	return defaultValue
}

// ColMap iterates through the collection and passes each value to the given callback
func ColMap[T any, R any](collection []T, callback func(T) R) []R {
	result := make([]R, len(collection))
	for i, item := range collection {
		result[i] = callback(item)
	}
	return result
}

// ColMax returns the maximum value of a given key
func ColMax[T any, V float64 | int | int64 | float32 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8](collection []T, valueFunc func(T) V) V {
	if len(collection) == 0 {
		var zero V
		return zero
	}

	max := valueFunc(collection[0])
	for i := 1; i < len(collection); i++ {
		value := valueFunc(collection[i])
		if value > max {
			max = value
		}
	}

	return max
}

// ColMerge merges the given array or collection with the original collection
func ColMerge[K comparable, V any](collection map[K]V, items map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range collection {
		result[k] = v
	}
	for k, v := range items {
		result[k] = v
	}
	return result
}

// ColMin returns the minimum value of a given key
func ColMin[T any, V float64 | int | int64 | float32 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8](collection []T, valueFunc func(T) V) V {
	if len(collection) == 0 {
		var zero V
		return zero
	}

	min := valueFunc(collection[0])
	for i := 1; i < len(collection); i++ {
		value := valueFunc(collection[i])
		if value < min {
			min = value
		}
	}

	return min
}

// ColOnly returns the items in the collection with the specified keys
func ColOnly[K comparable, V any](collection map[K]V, keys []K) map[K]V {
	result := make(map[K]V)

	// Create a map for faster lookup
	keysMap := make(map[K]struct{})
	for _, key := range keys {
		keysMap[key] = struct{}{}
	}

	for key, value := range collection {
		if _, exists := keysMap[key]; exists {
			result[key] = value
		}
	}

	return result
}

// ColPad fills the array to the specified size with a value
func ColPad[T any](collection []T, size int, value T) []T {
	if size <= len(collection) {
		return collection
	}

	result := make([]T, size)
	copy(result, collection)
	for i := len(collection); i < size; i++ {
		result[i] = value
	}
	return result
}

// ColPartition separates elements that pass a given truth test from those that don't
func ColPartition[T any](collection []T, predicate func(T) bool) ([]T, []T) {
	pass := make([]T, 0)
	fail := make([]T, 0)

	for _, item := range collection {
		if predicate(item) {
			pass = append(pass, item)
		} else {
			fail = append(fail, item)
		}
	}

	return pass, fail
}

// ColPluck retrieves all of the values for a given key
func ColPluck[T any, V any](collection []T, key func(T) V) []V {
	result := make([]V, len(collection))
	for i, item := range collection {
		result[i] = key(item)
	}
	return result
}

// ColPrepend adds an item to the beginning of the collection
func ColPrepend[T any](collection []T, values ...T) []T {
	result := make([]T, len(values)+len(collection))
	copy(result, values)
	copy(result[len(values):], collection)
	return result
}

// ColPull removes and returns an item from the collection by key
func ColPull[T any](collection []T, index int) (T, []T) {
	if index < 0 || index >= len(collection) {
		var zero T
		return zero, collection
	}

	item := collection[index]
	result := append(collection[:index], collection[index+1:]...)
	return item, result
}

// ColPush adds an item to the end of the collection
func ColPush[T any](collection []T, values ...T) []T {
	return append(collection, values...)
}

// ColPut sets the given key and value in the collection
func ColPut[K comparable, V any](collection map[K]V, key K, value V) map[K]V {
	result := make(map[K]V)
	for k, v := range collection {
		result[k] = v
	}
	result[key] = value
	return result
}

// ColRandom retrieves a random item from the collection
func ColRandom[T any](collection []T) (T, bool) {
	if len(collection) == 0 {
		var zero T
		return zero, false
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return collection[r.Intn(len(collection))], true
}

// ColRandomOrDefault retrieves a random item from the collection or a default value if the collection is empty
func ColRandomOrDefault[T any](collection []T, defaultValue T) T {
	if len(collection) == 0 {
		return defaultValue
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return collection[r.Intn(len(collection))]
}

// ColReduce reduces the collection to a single value
func ColReduce[T any, R any](collection []T, initialValue R, callback func(R, T) R) R {
	result := initialValue
	for _, item := range collection {
		result = callback(result, item)
	}
	return result
}

// ColReject filters the collection using the given callback
func ColReject[T any](collection []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range collection {
		if !predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// ColReverse reverses the order of the collection's items
func ColReverse[T any](collection []T) []T {
	result := make([]T, len(collection))
	for i, j := 0, len(collection)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = collection[j]
	}
	return result
}

// ColSearch searches the collection for a given value and returns the corresponding key if successful
func ColSearch[T comparable](collection []T, value T) (int, bool) {
	for i, item := range collection {
		if item == value {
			return i, true
		}
	}
	return -1, false
}

// ColSearchFunc searches the collection using the given callback
func ColSearchFunc[T any](collection []T, predicate func(T) bool) (int, bool) {
	for i, item := range collection {
		if predicate(item) {
			return i, true
		}
	}
	return -1, false
}

// ColShift removes and returns the first item from the collection
func ColShift[T any](collection []T) (T, []T) {
	if len(collection) == 0 {
		var zero T
		return zero, collection
	}

	item := collection[0]
	result := collection[1:]
	return item, result
}

// ColShuffle randomly shuffles the items in the collection
func ColShuffle[T any](collection []T) []T {
	result := make([]T, len(collection))
	copy(result, collection)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

// ColSlice returns a slice of the collection starting at the given index
func ColSlice[T any](collection []T, start int) []T {
	if start >= len(collection) {
		return []T{}
	}

	if start < 0 {
		start = len(collection) + start
		if start < 0 {
			start = 0
		}
	}

	return collection[start:]
}

// ColSliceWithLength returns a slice of the collection starting at the given index with the specified length
func ColSliceWithLength[T any](collection []T, start int, length int) []T {
	if start >= len(collection) || length <= 0 {
		return []T{}
	}

	if start < 0 {
		start = len(collection) + start
		if start < 0 {
			start = 0
		}
	}

	end := start + length
	if end > len(collection) {
		end = len(collection)
	}

	return collection[start:end]
}

// ColSort sorts the collection
func ColSort[T any](collection []T, less func(i, j T) bool) []T {
	result := make([]T, len(collection))
	copy(result, collection)

	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return result
}

// ColSortBy sorts the collection by the given key
func ColSortBy[T any, K comparable](collection []T, keyFunc func(T) K, less func(i, j K) bool) []T {
	result := make([]T, len(collection))
	copy(result, collection)

	sort.Slice(result, func(i, j int) bool {
		return less(keyFunc(result[i]), keyFunc(result[j]))
	})

	return result
}

// ColSortByDesc sorts the collection by the given key in descending order
func ColSortByDesc[T any, K comparable](collection []T, keyFunc func(T) K, less func(i, j K) bool) []T {
	result := make([]T, len(collection))
	copy(result, collection)

	sort.Slice(result, func(i, j int) bool {
		return less(keyFunc(result[j]), keyFunc(result[i]))
	})

	return result
}

// TODO ColSplice removes and returns a slice of items starting at the specified index
func ColSplice[T any](collection []T, start int, length int) ([]T, []T) {
	if start >= len(collection) || length <= 0 {
		return []T{}, collection
	}

	if start < 0 {
		start = len(collection) + start
		if start < 0 {
			start = 0
		}
	}

	end := start + length
	if end > len(collection) {
		end = len(collection)
	}

	removed := collection[start:end]
	result := append(collection[:start], collection[end:]...)
	return removed, result
}

// ColSplit breaks a collection into the given number of groups
func ColSplit[T any](collection []T, numberOfGroups int) [][]T {
	if numberOfGroups <= 0 {
		return [][]T{}
	}

	if len(collection) == 0 {
		return [][]T{}
	}

	result := make([][]T, numberOfGroups)
	for i := 0; i < numberOfGroups; i++ {
		result[i] = make([]T, 0)
	}

	for i, item := range collection {
		result[i%numberOfGroups] = append(result[i%numberOfGroups], item)
	}

	return result
}

// ColSum returns the sum of all items in the collection
func ColSum[T any, V float64 | int | int64 | float32 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8](collection []T, valueFunc func(T) V) V {
	var sum V
	for _, item := range collection {
		sum += valueFunc(item)
	}
	return sum
}

// ColTake returns a new collection with the specified number of items
func ColTake[T any](collection []T, limit int) []T {
	if limit <= 0 {
		return []T{}
	}

	if limit >= len(collection) {
		return collection
	}

	return collection[:limit]
}

// ColTap passes the collection to the given callback then returns the collection
func ColTap[T any](collection []T, callback func([]T)) []T {
	callback(collection)
	return collection
}

// ColUnique returns all of the unique items in the collection
func ColUnique[T comparable](collection []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)

	for _, item := range collection {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// ColUniqueBy returns all of the unique items in the collection using the given key
func ColUniqueBy[T any, K comparable](collection []T, keyFunc func(T) K) []T {
	seen := make(map[K]struct{})
	result := make([]T, 0)

	for _, item := range collection {
		key := keyFunc(item)
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// ColUnless executes the given callback when the first argument is falsy
func ColUnless[T any](condition bool, collection []T, callback func([]T) []T) []T {
	if !condition {
		return callback(collection)
	}
	return collection
}

// ColUnlessEmpty executes the given callback when the collection is not empty
func ColUnlessEmpty[T any](collection []T, callback func([]T) []T) []T {
	if len(collection) > 0 {
		return callback(collection)
	}
	return collection
}

// ColUnlessNotEmpty executes the given callback when the collection is empty
func ColUnlessNotEmpty[T any](collection []T, callback func([]T) []T) []T {
	if len(collection) == 0 {
		return callback(collection)
	}
	return collection
}

// ColValues returns all of the values in the collection
func ColValues[K comparable, V any](collection map[K]V) []V {
	values := make([]V, 0, len(collection))
	for _, value := range collection {
		values = append(values, value)
	}
	return values
}

// ColWhen executes the given callback when the first argument is truthy
func ColWhen[T any](condition bool, collection []T, callback func([]T) []T) []T {
	if condition {
		return callback(collection)
	}
	return collection
}

// ColWhenEmpty executes the given callback when the collection is empty
func ColWhenEmpty[T any](collection []T, callback func([]T) []T) []T {
	if len(collection) == 0 {
		return callback(collection)
	}
	return collection
}

// ColWhenNotEmpty executes the given callback when the collection is not empty
func ColWhenNotEmpty[T any](collection []T, callback func([]T) []T) []T {
	if len(collection) > 0 {
		return callback(collection)
	}
	return collection
}

// ColWhere filters the collection using the given callback
func ColWhere[T any](collection []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range collection {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// ColWhereIn filters the collection by a given key / value contained within the given array
func ColWhereIn[T any, K comparable](collection []T, keyFunc func(T) K, values []K) []T {
	// Create a map for faster lookup
	valuesMap := make(map[K]struct{})
	for _, value := range values {
		valuesMap[value] = struct{}{}
	}

	result := make([]T, 0)
	for _, item := range collection {
		key := keyFunc(item)
		if _, exists := valuesMap[key]; exists {
			result = append(result, item)
		}
	}

	return result
}

// ColWhereNotIn filters the collection by a given key / value not contained within the given array
func ColWhereNotIn[T any, K comparable](collection []T, keyFunc func(T) K, values []K) []T {
	// Create a map for faster lookup
	valuesMap := make(map[K]struct{})
	for _, value := range values {
		valuesMap[value] = struct{}{}
	}

	result := make([]T, 0)
	for _, item := range collection {
		key := keyFunc(item)
		if _, exists := valuesMap[key]; !exists {
			result = append(result, item)
		}
	}

	return result
}

// ColZip merges together the values of the given arrays with the values of the original collection
func ColZip[T any](collection []T, arrays ...[]T) [][]T {
	if len(collection) == 0 {
		return [][]T{}
	}

	// Determine the length of the result (minimum length of all arrays)
	minLength := len(collection)
	for _, array := range arrays {
		if len(array) < minLength {
			minLength = len(array)
		}
	}

	// Create the result array
	result := make([][]T, minLength)
	for i := 0; i < minLength; i++ {
		// Each inner array has a length of 1 (for the collection) + len(arrays)
		result[i] = make([]T, 1+len(arrays))
		result[i][0] = collection[i]

		// Add items from the other arrays
		for j, array := range arrays {
			result[i][j+1] = array[i]
		}
	}

	return result
}
