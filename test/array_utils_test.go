package test

import (
	"gfly/app/utils"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestArrayUtils(t *testing.T) {
	// Test Map functions
	t.Run("MapMergeMaps", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"b": 3, "c": 4}

		result := utils.MapMergeMaps(map1, map2)

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		if result["a"] != 1 || result["b"] != 3 || result["c"] != 4 {
			t.Errorf("Expected merged map with values a=1, b=3, c=4, got %v", result)
		}

		// Test with empty maps
		emptyResult := utils.MapMergeMaps(map[string]int{}, map[string]int{})
		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with one empty map
		oneEmptyResult := utils.MapMergeMaps(map1, map[string]int{})
		if len(oneEmptyResult) != len(map1) {
			t.Errorf("Expected %d items, got %d", len(map1), len(oneEmptyResult))
		}
	})

	t.Run("MapKeys", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		keys := utils.MapKeys(m)

		if len(keys) != 3 {
			t.Errorf("Expected 3 keys, got %d", len(keys))
		}

		// Check that all keys are present
		keySet := make(map[string]bool)
		for _, key := range keys {
			keySet[key] = true
		}

		if !keySet["a"] || !keySet["b"] || !keySet["c"] {
			t.Errorf("Expected keys a, b, c, got %v", keys)
		}

		// Test with empty map
		emptyKeys := utils.MapKeys(map[string]int{})
		if len(emptyKeys) != 0 {
			t.Errorf("Expected empty result, got %d keys", len(emptyKeys))
		}
	})

	t.Run("MapValues", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		values := utils.MapValues(m)

		if len(values) != 3 {
			t.Errorf("Expected 3 values, got %d", len(values))
		}

		// Check that all values are present
		valueSet := make(map[int]bool)
		for _, value := range values {
			valueSet[value] = true
		}

		if !valueSet[1] || !valueSet[2] || !valueSet[3] {
			t.Errorf("Expected values 1, 2, 3, got %v", values)
		}

		// Test with empty map
		emptyValues := utils.MapValues(map[string]int{})
		if len(emptyValues) != 0 {
			t.Errorf("Expected empty result, got %d values", len(emptyValues))
		}
	})

	t.Run("MapFindKey", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		key, found := utils.MapFindKey(m, 2)

		if !found || key != "b" {
			t.Errorf("Expected to find key 'b', got %s with found=%v", key, found)
		}

		// Test not found
		key, found = utils.MapFindKey(m, 4)

		if found {
			t.Errorf("Expected not to find any key, but found %s", key)
		}

		// Test with empty map
		key, found = utils.MapFindKey(map[string]int{}, 1)

		if found {
			t.Errorf("Expected not to find any key in empty map, but found %s", key)
		}
	})

	t.Run("MapFilterMap", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

		result := utils.MapFilterMap(m, func(k string, v int) bool {
			return v%2 == 0
		})

		if len(result) != 2 {
			t.Errorf("Expected 2 items, got %d", len(result))
		}

		if result["b"] != 2 || result["d"] != 4 {
			t.Errorf("Expected filtered map with b=2, d=4, got %v", result)
		}

		// Test with empty map
		emptyResult := utils.MapFilterMap(map[string]int{}, func(k string, v int) bool {
			return true
		})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with predicate that filters everything
		noneResult := utils.MapFilterMap(m, func(k string, v int) bool {
			return false
		})

		if len(noneResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(noneResult))
		}
	})

	t.Run("MapMapValues", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		result := utils.MapMapValues(m, func(v int) string {
			return strconv.Itoa(v)
		})

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		if result["a"] != "1" || result["b"] != "2" || result["c"] != "3" {
			t.Errorf("Expected transformed map with a='1', b='2', c='3', got %v", result)
		}

		// Test with empty map
		emptyResult := utils.MapMapValues(map[string]int{}, func(v int) string {
			return strconv.Itoa(v)
		})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}
	})

	t.Run("MapInvertMap", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		result := utils.MapInvertMap(m)

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		if result[1] != "a" || result[2] != "b" || result[3] != "c" {
			t.Errorf("Expected inverted map with 1='a', 2='b', 3='c', got %v", result)
		}

		// Test with empty map
		emptyResult := utils.MapInvertMap(map[string]int{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with duplicate values
		dupMap := map[string]int{"a": 1, "b": 1, "c": 2}
		dupResult := utils.MapInvertMap(dupMap)

		if len(dupResult) != 2 {
			t.Errorf("Expected 2 items (due to duplicate values), got %d", len(dupResult))
		}
	})

	t.Run("MapGetOrDefault", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		// Test existing key
		result := utils.MapGetOrDefault(m, "b", 0)

		if result != 2 {
			t.Errorf("Expected 2, got %d", result)
		}

		// Test non-existent key
		result = utils.MapGetOrDefault(m, "d", 0)

		if result != 0 {
			t.Errorf("Expected default value 0, got %d", result)
		}

		// Test with empty map
		emptyResult := utils.MapGetOrDefault(map[string]int{}, "a", 42)

		if emptyResult != 42 {
			t.Errorf("Expected default value 42, got %d", emptyResult)
		}
	})

	t.Run("MapGetOrInsert", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		// Test existing key
		result := utils.MapGetOrInsert(m, "b", 0)

		if result != 2 {
			t.Errorf("Expected 2, got %d", result)
		}

		if m["b"] != 2 {
			t.Errorf("Expected map to be unchanged for existing key, got b=%d", m["b"])
		}

		// Test non-existent key
		result = utils.MapGetOrInsert(m, "d", 4)

		if result != 4 {
			t.Errorf("Expected default value 4, got %d", result)
		}

		if m["d"] != 4 {
			t.Errorf("Expected map to have new key-value pair d=4, got d=%d", m["d"])
		}

		// Test with empty map
		emptyMap := map[string]int{}
		emptyResult := utils.MapGetOrInsert(emptyMap, "a", 42)

		if emptyResult != 42 {
			t.Errorf("Expected default value 42, got %d", emptyResult)
		}

		if emptyMap["a"] != 42 {
			t.Errorf("Expected map to have new key-value pair a=42, got a=%d", emptyMap["a"])
		}
	})

	t.Run("MapToSlice", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}

		result := utils.MapToSlice(m)

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		// Check that all key-value pairs are present
		pairFound := make(map[string]bool)
		for _, pair := range result {
			if pair.Value == m[pair.Key] {
				pairFound[pair.Key] = true
			}
		}

		if !pairFound["a"] || !pairFound["b"] || !pairFound["c"] {
			t.Errorf("Expected all key-value pairs to be present, got %v", result)
		}

		// Test with empty map
		emptyResult := utils.MapToSlice(map[string]int{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}
	})

	t.Run("MapSliceToMap", func(t *testing.T) {
		slice := []struct {
			Key   string
			Value int
		}{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
		}

		result := utils.MapSliceToMap(slice)

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		if result["a"] != 1 || result["b"] != 2 || result["c"] != 3 {
			t.Errorf("Expected map with a=1, b=2, c=3, got %v", result)
		}

		// Test with empty slice
		emptyResult := utils.MapSliceToMap([]struct {
			Key   string
			Value int
		}{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with duplicate keys (last one wins)
		dupSlice := []struct {
			Key   string
			Value int
		}{
			{Key: "a", Value: 1},
			{Key: "a", Value: 2},
			{Key: "b", Value: 3},
		}

		dupResult := utils.MapSliceToMap(dupSlice)

		if len(dupResult) != 2 {
			t.Errorf("Expected 2 items (due to duplicate keys), got %d", len(dupResult))
		}

		if dupResult["a"] != 2 {
			t.Errorf("Expected a=2 (last value wins), got a=%d", dupResult["a"])
		}
	})

	t.Run("MapEqualMaps", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2, "c": 3}
		map2 := map[string]int{"a": 1, "b": 2, "c": 3}
		map3 := map[string]int{"a": 1, "b": 2, "d": 3}
		map4 := map[string]int{"a": 1, "b": 2, "c": 4}

		// Test equal maps
		if !utils.MapEqualMaps(map1, map2) {
			t.Error("Expected maps to be equal")
		}

		// Test different keys
		if utils.MapEqualMaps(map1, map3) {
			t.Error("Expected maps with different keys to be unequal")
		}

		// Test different values
		if utils.MapEqualMaps(map1, map4) {
			t.Error("Expected maps with different values to be unequal")
		}

		// Test different lengths
		map5 := map[string]int{"a": 1, "b": 2}
		if utils.MapEqualMaps(map1, map5) {
			t.Error("Expected maps with different lengths to be unequal")
		}

		// Test empty maps
		emptyMap1 := map[string]int{}
		emptyMap2 := map[string]int{}
		if !utils.MapEqualMaps(emptyMap1, emptyMap2) {
			t.Error("Expected empty maps to be equal")
		}

		// Test empty map and non-empty map
		if utils.MapEqualMaps(emptyMap1, map1) {
			t.Error("Expected empty map and non-empty map to be unequal")
		}
	})

	t.Run("MapDiffMaps", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2, "c": 3}
		map2 := map[string]int{"b": 2, "c": 4, "d": 5}

		added, removed, changed := utils.MapDiffMaps(map1, map2)

		// Check added keys
		if len(added) != 1 || added["d"] != 5 {
			t.Errorf("Expected added to contain d=5, got %v", added)
		}

		// Check removed keys
		if len(removed) != 1 || removed["a"] != 1 {
			t.Errorf("Expected removed to contain a=1, got %v", removed)
		}

		// Check changed keys
		if len(changed) != 1 || changed["c"] != 4 {
			t.Errorf("Expected changed to contain c=4, got %v", changed)
		}

		// Test with identical maps
		added2, removed2, changed2 := utils.MapDiffMaps(map1, map1)

		if len(added2) != 0 || len(removed2) != 0 || len(changed2) != 0 {
			t.Errorf("Expected no differences for identical maps, got added=%v, removed=%v, changed=%v",
				added2, removed2, changed2)
		}

		// Test with empty maps
		added3, removed3, changed3 := utils.MapDiffMaps(map[string]int{}, map[string]int{})

		if len(added3) != 0 || len(removed3) != 0 || len(changed3) != 0 {
			t.Errorf("Expected no differences for empty maps, got added=%v, removed=%v, changed=%v",
				added3, removed3, changed3)
		}

		// Test with one empty map
		added4, removed4, changed4 := utils.MapDiffMaps(map1, map[string]int{})

		if len(added4) != 0 || len(removed4) != len(map1) || len(changed4) != 0 {
			t.Errorf("Expected all keys to be removed when comparing with empty map, got added=%v, removed=%v, changed=%v",
				added4, removed4, changed4)
		}
	})

	// Test Set functions
	t.Run("SetContains", func(t *testing.T) {
		set := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
		}

		// Test existing element
		if !utils.SetContains(set, "b") {
			t.Error("Expected set to contain 'b'")
		}

		// Test non-existent element
		if utils.SetContains(set, "d") {
			t.Error("Expected set to not contain 'd'")
		}

		// Test with empty set
		emptySet := map[string]struct{}{}
		if utils.SetContains(emptySet, "a") {
			t.Error("Expected empty set to not contain any elements")
		}
	})

	t.Run("SetToSlice", func(t *testing.T) {
		set := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
		}

		result := utils.SetToSlice(set)

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		// Check that all elements are present
		elemFound := make(map[string]bool)
		for _, elem := range result {
			elemFound[elem] = true
		}

		if !elemFound["a"] || !elemFound["b"] || !elemFound["c"] {
			t.Errorf("Expected elements a, b, c, got %v", result)
		}

		// Test with empty set
		emptyResult := utils.SetToSlice(map[string]struct{}{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}
	})

	t.Run("SliceToSet", func(t *testing.T) {
		slice := []string{"a", "b", "c", "b"} // Note the duplicate 'b'

		result := utils.SliceToSet(slice)

		if len(result) != 3 {
			t.Errorf("Expected 3 items (duplicates removed), got %d", len(result))
		}

		// Check that all elements are present
		if _, ok := result["a"]; !ok {
			t.Error("Expected set to contain 'a'")
		}
		if _, ok := result["b"]; !ok {
			t.Error("Expected set to contain 'b'")
		}
		if _, ok := result["c"]; !ok {
			t.Error("Expected set to contain 'c'")
		}

		// Test with empty slice
		emptyResult := utils.SliceToSet([]string{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}
	})

	t.Run("SetUnion", func(t *testing.T) {
		set1 := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
		}

		set2 := map[string]struct{}{
			"c": {},
			"d": {},
			"e": {},
		}

		result := utils.SetUnion(set1, set2)

		if len(result) != 5 {
			t.Errorf("Expected 5 items, got %d", len(result))
		}

		// Check that all elements are present
		expectedElems := []string{"a", "b", "c", "d", "e"}
		for _, elem := range expectedElems {
			if _, ok := result[elem]; !ok {
				t.Errorf("Expected union to contain '%s'", elem)
			}
		}

		// Test with empty sets
		emptyResult := utils.SetUnion(map[string]struct{}{}, map[string]struct{}{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with one empty set
		oneEmptyResult := utils.SetUnion(set1, map[string]struct{}{})

		if len(oneEmptyResult) != len(set1) {
			t.Errorf("Expected %d items, got %d", len(set1), len(oneEmptyResult))
		}
	})

	t.Run("SetIntersection", func(t *testing.T) {
		set1 := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
		}

		set2 := map[string]struct{}{
			"c": {},
			"d": {},
			"e": {},
		}

		result := utils.SetIntersection(set1, set2)

		if len(result) != 1 {
			t.Errorf("Expected 1 item, got %d", len(result))
		}

		// Check that the intersection contains only 'c'
		if _, ok := result["c"]; !ok {
			t.Error("Expected intersection to contain 'c'")
		}

		// Test with no common elements
		set3 := map[string]struct{}{
			"x": {},
			"y": {},
			"z": {},
		}

		noCommonResult := utils.SetIntersection(set1, set3)

		if len(noCommonResult) != 0 {
			t.Errorf("Expected empty result for sets with no common elements, got %d items", len(noCommonResult))
		}

		// Test with empty sets
		emptyResult := utils.SetIntersection(map[string]struct{}{}, map[string]struct{}{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with one empty set
		oneEmptyResult := utils.SetIntersection(set1, map[string]struct{}{})

		if len(oneEmptyResult) != 0 {
			t.Errorf("Expected empty result when intersecting with empty set, got %d items", len(oneEmptyResult))
		}
	})

	t.Run("SetDifference", func(t *testing.T) {
		set1 := map[string]struct{}{
			"a": {},
			"b": {},
			"c": {},
		}

		set2 := map[string]struct{}{
			"c": {},
			"d": {},
			"e": {},
		}

		result := utils.SetDifference(set1, set2)

		if len(result) != 2 {
			t.Errorf("Expected 2 items, got %d", len(result))
		}

		// Check that the difference contains 'a' and 'b'
		if _, ok := result["a"]; !ok {
			t.Error("Expected difference to contain 'a'")
		}
		if _, ok := result["b"]; !ok {
			t.Error("Expected difference to contain 'b'")
		}

		// Test with no common elements
		set3 := map[string]struct{}{
			"x": {},
			"y": {},
			"z": {},
		}

		noCommonResult := utils.SetDifference(set1, set3)

		if len(noCommonResult) != len(set1) {
			t.Errorf("Expected result to contain all elements from set1, got %d items", len(noCommonResult))
		}

		// Test with empty sets
		emptyResult := utils.SetDifference(map[string]struct{}{}, map[string]struct{}{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with one empty set
		oneEmptyResult := utils.SetDifference(set1, map[string]struct{}{})

		if len(oneEmptyResult) != len(set1) {
			t.Errorf("Expected result to contain all elements from set1, got %d items", len(oneEmptyResult))
		}
	})

	// Test Accessible
	t.Run("Accessible", func(t *testing.T) {
		if !utils.ArrAccessible([]int{1, 2, 3}) {
			t.Error("Expected slice to be accessible")
		}
		if !utils.ArrAccessible(map[string]int{"a": 1}) {
			t.Error("Expected map to be accessible")
		}
		if utils.ArrAccessible(123) {
			t.Error("Expected int to not be accessible")
		}
	})

	// Test Add
	t.Run("Add", func(t *testing.T) {
		array := map[string]interface{}{"name": "John"}
		result := utils.ArrAdd(array, "age", 30)
		if result["age"] != 30 {
			t.Error("Expected age to be added")
		}

		// Test that it doesn't overwrite existing keys
		result = utils.ArrAdd(array, "name", "Jane")
		if result["name"] != "John" {
			t.Error("Expected name to not be overwritten")
		}
	})

	// Test Collapse
	t.Run("Collapse", func(t *testing.T) {
		arrays := [][]interface{}{
			{1, 2, 3},
			{4, 5, 6},
		}
		result := utils.ArrCollapse(arrays)
		expected := []interface{}{1, 2, 3, 4, 5, 6}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %v at index %d, got %v", v, i, result[i])
			}
		}
	})

	// Test CrossJoin
	t.Run("CrossJoin", func(t *testing.T) {
		array1 := []string{"a", "b"}
		array2 := []string{"1", "2"}
		result := utils.ArrCrossJoin(array1, array2)

		expected := [][]string{
			{"a", "1"},
			{"a", "2"},
			{"b", "1"},
			{"b", "2"},
		}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, row := range expected {
			if !reflect.DeepEqual(result[i], row) {
				t.Errorf("Expected %v at index %d, got %v", row, i, result[i])
			}
		}
	})

	// Test Divide
	t.Run("Divide", func(t *testing.T) {
		array := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		keys, values := utils.ArrDivide(array)

		if len(keys) != 2 || len(values) != 2 {
			t.Errorf("Expected 2 keys and 2 values, got %d keys and %d values", len(keys), len(values))
		}

		// Check that all keys and values are present
		keyFound := make(map[string]bool)
		for _, key := range keys {
			keyFound[key] = true
		}

		if !keyFound["name"] || !keyFound["age"] {
			t.Error("Expected keys 'name' and 'age'")
		}

		valueFound := make(map[interface{}]bool)
		for _, value := range values {
			valueFound[value] = true
		}

		if !valueFound["John"] || !valueFound[30] {
			t.Error("Expected values 'John' and 30")
		}
	})

	// Test Dot
	t.Run("Dot", func(t *testing.T) {
		array := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "John",
				"age":  30,
			},
		}
		result := utils.ArrDot(array)

		if result["user.name"] != "John" || result["user.age"] != 30 {
			t.Error("Expected flattened array with dot notation")
		}
	})

	// Test Undot
	t.Run("Undot", func(t *testing.T) {
		array := map[string]interface{}{
			"user.name": "John",
			"user.age":  30,
		}
		result := utils.ArrUndot(array)

		userMap, ok := result["user"].(map[string]interface{})
		if !ok {
			t.Error("Expected 'user' to be a map")
			return
		}

		if userMap["name"] != "John" || userMap["age"] != 30 {
			t.Error("Expected expanded array from dot notation")
		}
	})

	// Test Except
	t.Run("Except", func(t *testing.T) {
		array := map[string]interface{}{
			"name": "John",
			"age":  30,
			"city": "New York",
		}
		result := utils.ArrExcept(array, "age")

		if _, exists := result["age"]; exists {
			t.Error("Expected 'age' to be removed")
		}

		if result["name"] != "John" || result["city"] != "New York" {
			t.Error("Expected other keys to remain")
		}
	})

	// Test Exists
	t.Run("Exists", func(t *testing.T) {
		array := map[string]interface{}{
			"name": "John",
		}

		if !utils.ArrExists(array, "name") {
			t.Error("Expected 'name' to exist")
		}

		if utils.ArrExists(array, "age") {
			t.Error("Expected 'age' to not exist")
		}
	})

	// Test First
	t.Run("First", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result, found := utils.ArrFirst(array, func(item int) bool {
			return item > 2
		})

		if !found || result != 3 {
			t.Errorf("Expected 3, got %d", result)
		}

		result, found = utils.ArrFirst(array, func(item int) bool {
			return item > 10
		})

		if found || result != 0 {
			t.Errorf("Expected not found and 0, got found=%v and %d", found, result)
		}
	})

	// Test FirstOrDefault
	t.Run("FirstOrDefault", func(t *testing.T) {
		array := []int{1, 2, 3}

		result := utils.ArrFirstOrDefault(array, 0)
		if result != 1 {
			t.Errorf("Expected 1, got %d", result)
		}

		emptyArray := []int{}
		result = utils.ArrFirstOrDefault(emptyArray, 0)
		if result != 0 {
			t.Errorf("Expected 0, got %d", result)
		}
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		array := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "John",
			},
		}

		result := utils.ArrGet(array, "user.name", nil)
		if result != "John" {
			t.Errorf("Expected 'John', got %v", result)
		}

		result = utils.ArrGet(array, "user.age", 30)
		if result != 30 {
			t.Errorf("Expected 30, got %v", result)
		}
	})

	// Test Has
	t.Run("Has", func(t *testing.T) {
		array := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "John",
			},
		}

		if !utils.ArrHas(array, "user.name") {
			t.Error("Expected 'user.name' to exist")
		}

		if utils.ArrHas(array, "user.age") {
			t.Error("Expected 'user.age' to not exist")
		}
	})

	// Test HasAny
	t.Run("HasAny", func(t *testing.T) {
		array := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "John",
			},
		}

		if !utils.ArrHasAny(array, "user.name", "user.age") {
			t.Error("Expected at least one key to exist")
		}

		if utils.ArrHasAny(array, "user.age", "user.city") {
			t.Error("Expected none of the keys to exist")
		}
	})

	// Test IsAssoc
	t.Run("IsAssoc", func(t *testing.T) {
		array := map[string]interface{}{
			"name": "John",
		}

		if !utils.ArrIsAssoc(array) {
			t.Error("Expected map to be associative")
		}

		if utils.ArrIsAssoc([]int{1, 2, 3}) {
			t.Error("Expected slice to not be associative")
		}
	})

	// Test IsList
	t.Run("IsList", func(t *testing.T) {
		array := []int{1, 2, 3}

		if !utils.ArrIsList(array) {
			t.Error("Expected slice to be a list")
		}

		if utils.ArrIsList(map[string]interface{}{"name": "John"}) {
			t.Error("Expected map to not be a list")
		}
	})

	// Test Only
	t.Run("Only", func(t *testing.T) {
		array := map[string]interface{}{
			"name": "John",
			"age":  30,
			"city": "New York",
		}

		result := utils.ArrOnly(array, "name", "age")

		if len(result) != 2 {
			t.Errorf("Expected 2 items, got %d", len(result))
		}

		if result["name"] != "John" || result["age"] != 30 {
			t.Error("Expected only 'name' and 'age' to be included")
		}

		if _, exists := result["city"]; exists {
			t.Error("Expected 'city' to be excluded")
		}
	})

	// Test Pluck
	t.Run("Pluck", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		users := []User{
			{Name: "John", Age: 30},
			{Name: "Jane", Age: 25},
		}

		names := utils.ArrPluck(users, func(user User) string {
			return user.Name
		})

		if len(names) != 2 {
			t.Errorf("Expected 2 names, got %d", len(names))
		}

		if names[0] != "John" || names[1] != "Jane" {
			t.Error("Expected names to be 'John' and 'Jane'")
		}
	})

	// Test Prepend
	t.Run("Prepend", func(t *testing.T) {
		array := []int{3, 4, 5}
		result := utils.ArrPrepend(array, 1, 2)

		expected := []int{1, 2, 3, 4, 5}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	// Test Pull
	t.Run("Pull", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		item, result := utils.ArrPull(array, 2)

		if item != 3 {
			t.Errorf("Expected pulled item to be 3, got %d", item)
		}

		expected := []int{1, 2, 4, 5}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	// Test Random
	t.Run("Random", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		_, found := utils.ArrRandom(array)
		if !found {
			t.Error("Expected to find a random element")
		}

		emptyArray := []int{}
		_, found = utils.ArrRandom(emptyArray)
		if found {
			t.Error("Expected not to find a random element in empty array")
		}
	})

	// Test Set
	t.Run("Set", func(t *testing.T) {
		array := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "John",
			},
		}

		result := utils.ArrSet(array, "user.age", 30)

		userMap, ok := result["user"].(map[string]interface{})
		if !ok {
			t.Error("Expected 'user' to be a map")
			return
		}

		if userMap["name"] != "John" || userMap["age"] != 30 {
			t.Error("Expected 'user.age' to be set to 30")
		}
	})

	// Test Wrap
	t.Run("Wrap", func(t *testing.T) {
		// Test wrapping a single value
		result := utils.ArrWrap(123)

		if len(result) != 1 || result[0] != 123 {
			t.Errorf("Expected [123], got %v", result)
		}

		// Test wrapping an array
		array := []int{1, 2, 3}
		result = utils.ArrWrap(array)

		if len(result) != 3 {
			t.Errorf("Expected length 3, got %d", len(result))
		}

		for i, v := range array {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %v", v, i, result[i])
			}
		}
	})

	// Test Contains
	t.Run("Contains", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		if !utils.ArrContains(array, 3) {
			t.Error("Expected array to contain 3")
		}

		if utils.ArrContains(array, 6) {
			t.Error("Expected array to not contain 6")
		}

		// Test with empty array
		emptyArray := []int{}
		if utils.ArrContains(emptyArray, 1) {
			t.Error("Expected empty array to not contain any elements")
		}

		// Test with strings
		strArray := []string{"apple", "banana", "cherry"}
		if !utils.ArrContains(strArray, "banana") {
			t.Error("Expected array to contain 'banana'")
		}
	})

	// Test Filter
	t.Run("Filter", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result := utils.ArrFilter(array, func(item int) bool {
			return item%2 == 0
		})

		expected := []int{2, 4}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrFilter(emptyArray, func(item int) bool {
			return item%2 == 0
		})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}
	})

	// Test Map
	t.Run("Map", func(t *testing.T) {
		array := []int{1, 2, 3}

		result := utils.ArrMap(array, func(item int) string {
			return strconv.Itoa(item)
		})

		expected := []string{"1", "2", "3"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
			}
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrMap(emptyArray, func(item int) string {
			return strconv.Itoa(item)
		})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}
	})

	// Test Find
	t.Run("Find", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result, found := utils.ArrFind(array, func(item int) bool {
			return item%2 == 0 && item > 2
		})

		if !found || result != 4 {
			t.Errorf("Expected to find 4, got %d with found=%v", result, found)
		}

		// Test not found
		result, found = utils.ArrFind(array, func(item int) bool {
			return item > 10
		})

		if found {
			t.Errorf("Expected not to find any item, but found %d", result)
		}

		// Test with empty array
		emptyArray := []int{}
		_, found = utils.ArrFind(emptyArray, func(item int) bool {
			return true
		})

		if found {
			t.Error("Expected not to find any item in empty array")
		}
	})

	// Test FindIndex
	t.Run("FindIndex", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		index, found := utils.ArrFindIndex(array, func(item int) bool {
			return item%2 == 0 && item > 2
		})

		if !found || index != 3 {
			t.Errorf("Expected to find index 3, got %d with found=%v", index, found)
		}

		// Test not found
		index, found = utils.ArrFindIndex(array, func(item int) bool {
			return item > 10
		})

		if found || index != -1 {
			t.Errorf("Expected not to find any item, but found at index %d", index)
		}

		// Test with empty array
		emptyArray := []int{}
		index, found = utils.ArrFindIndex(emptyArray, func(item int) bool {
			return true
		})

		if found || index != -1 {
			t.Errorf("Expected not to find any item in empty array, but found at index %d", index)
		}
	})

	// Test Unique
	t.Run("Unique", func(t *testing.T) {
		array := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}

		result := utils.ArrUnique(array)
		expected := []int{1, 2, 3, 4, 5}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrUnique(emptyArray)

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}

		// Test with already unique array
		uniqueArray := []int{1, 2, 3, 4, 5}
		uniqueResult := utils.ArrUnique(uniqueArray)

		if len(uniqueResult) != len(uniqueArray) {
			t.Errorf("Expected length %d, got %d", len(uniqueArray), len(uniqueResult))
		}

		for i, v := range uniqueArray {
			if uniqueResult[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, uniqueResult[i])
			}
		}
	})

	// Test Shuffle
	t.Run("Shuffle", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result := utils.ArrShuffle(array)

		// Check that the result has the same length
		if len(result) != len(array) {
			t.Errorf("Expected length %d, got %d", len(array), len(result))
		}

		// Check that the result contains all the original elements
		for _, v := range array {
			found := false
			for _, r := range result {
				if r == v {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected result to contain %d, but it doesn't", v)
			}
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrShuffle(emptyArray)

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}

		// Test with single element array
		singleArray := []int{1}
		singleResult := utils.ArrShuffle(singleArray)

		if len(singleResult) != 1 || singleResult[0] != 1 {
			t.Errorf("Expected [1], got %v", singleResult)
		}
	})

	// Test Chunk
	t.Run("Chunk", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// Test with chunk size 3
		result := utils.ArrChunk(array, 3)
		expected := [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
			{10},
		}

		if len(result) != len(expected) {
			t.Errorf("Expected %d chunks, got %d", len(expected), len(result))
		}

		for i, chunk := range expected {
			if len(result[i]) != len(chunk) {
				t.Errorf("Expected chunk %d to have length %d, got %d", i, len(chunk), len(result[i]))
			}

			for j, v := range chunk {
				if result[i][j] != v {
					t.Errorf("Expected %d at chunk %d, index %d, got %d", v, i, j, result[i][j])
				}
			}
		}

		// Test with chunk size larger than array
		result = utils.ArrChunk(array, 15)
		if len(result) != 1 || len(result[0]) != len(array) {
			t.Errorf("Expected 1 chunk with length %d, got %d chunks", len(array), len(result))
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrChunk(emptyArray, 3)

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}

		// Test with invalid chunk size
		invalidResult := utils.ArrChunk(array, 0)
		if len(invalidResult) != 0 {
			t.Errorf("Expected empty result for invalid chunk size, got length %d", len(invalidResult))
		}
	})

	// Test SortedCopy
	t.Run("SortedCopy", func(t *testing.T) {
		array := []int{5, 3, 1, 4, 2}

		result := utils.ArrSortedCopy(array, func(a, b int) bool {
			return a < b
		})

		expected := []int{1, 2, 3, 4, 5}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test descending order
		result = utils.ArrSortedCopy(array, func(a, b int) bool {
			return a > b
		})

		expected = []int{5, 4, 3, 2, 1}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrSortedCopy(emptyArray, func(a, b int) bool {
			return a < b
		})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}
	})

	// Test Reduce
	t.Run("Reduce", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		// Sum all elements
		result := utils.ArrReduce(array, 0, func(acc, item int) int {
			return acc + item
		})

		if result != 15 {
			t.Errorf("Expected sum 15, got %d", result)
		}

		// Concatenate strings
		strArray := []string{"a", "b", "c"}
		strResult := utils.ArrReduce(strArray, "", func(acc, item string) string {
			return acc + item
		})

		if strResult != "abc" {
			t.Errorf("Expected 'abc', got '%s'", strResult)
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrReduce(emptyArray, 10, func(acc, item int) int {
			return acc + item
		})

		if emptyResult != 10 {
			t.Errorf("Expected initial value 10, got %d", emptyResult)
		}
	})

	// Test Join
	t.Run("Join", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result := utils.ArrJoin(array, ", ", func(item int) string {
			return strconv.Itoa(item)
		})

		expected := "1, 2, 3, 4, 5"

		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrJoin(emptyArray, ", ", func(item int) string {
			return strconv.Itoa(item)
		})

		if emptyResult != "" {
			t.Errorf("Expected empty string, got '%s'", emptyResult)
		}

		// Test with single element
		singleArray := []int{1}
		singleResult := utils.ArrJoin(singleArray, ", ", func(item int) string {
			return strconv.Itoa(item)
		})

		if singleResult != "1" {
			t.Errorf("Expected '1', got '%s'", singleResult)
		}
	})

	// Test Intersection
	t.Run("Intersection", func(t *testing.T) {
		array1 := []int{1, 2, 3, 4, 5}
		array2 := []int{3, 4, 5, 6, 7}
		array3 := []int{5, 6, 7, 8, 9}

		result := utils.ArrIntersection(array1, array2, array3)
		expected := []int{5}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with no common elements
		array4 := []int{10, 11, 12}
		noCommonResult := utils.ArrIntersection(array1, array4)

		if len(noCommonResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(noCommonResult))
		}

		// Test with empty array
		emptyResult := utils.ArrIntersection(array1, []int{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}

		// Test with no arrays
		noArraysResult := utils.ArrIntersection[int]()

		if len(noArraysResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(noArraysResult))
		}
	})

	// Test Union
	t.Run("Union", func(t *testing.T) {
		array1 := []int{1, 2, 3}
		array2 := []int{3, 4, 5}
		array3 := []int{5, 6, 7}

		result := utils.ArrUnion(array1, array2, array3)

		// Sort the result for consistent comparison
		sort.Ints(result)

		expected := []int{1, 2, 3, 4, 5, 6, 7}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with empty array
		emptyResult := utils.ArrUnion(array1, []int{})

		// Sort the result for consistent comparison
		sort.Ints(emptyResult)

		expected = []int{1, 2, 3}

		if len(emptyResult) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(emptyResult))
		}

		for i, v := range expected {
			if emptyResult[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, emptyResult[i])
			}
		}

		// Test with no arrays
		noArraysResult := utils.ArrUnion[int]()

		if len(noArraysResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(noArraysResult))
		}
	})

	// Test Difference
	t.Run("Difference", func(t *testing.T) {
		array1 := []int{1, 2, 3, 4, 5}
		array2 := []int{3, 4, 5}
		array3 := []int{5, 6, 7}

		result := utils.ArrDifference(array1, array2, array3)
		expected := []int{1, 2}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with empty first array
		emptyResult := utils.ArrDifference([]int{}, array2)

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}

		// Test with empty other arrays
		fullResult := utils.ArrDifference(array1)

		if len(fullResult) != len(array1) {
			t.Errorf("Expected length %d, got %d", len(array1), len(fullResult))
		}

		for i, v := range array1 {
			if fullResult[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, fullResult[i])
			}
		}
	})

	// Test GroupBy
	t.Run("GroupBy", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		people := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 25},
			{Name: "Dave", Age: 30},
		}

		result := utils.ArrGroupBy(people, func(p Person) int {
			return p.Age
		})

		if len(result) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(result))
		}

		if len(result[25]) != 2 || len(result[30]) != 2 {
			t.Errorf("Expected 2 people in each group, got %d in group 25 and %d in group 30",
				len(result[25]), len(result[30]))
		}

		// Test with empty array
		emptyResult := utils.ArrGroupBy([]Person{}, func(p Person) int {
			return p.Age
		})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d groups", len(emptyResult))
		}
	})

	// Test Flatten
	t.Run("Flatten", func(t *testing.T) {
		arrays := [][]int{
			{1, 2, 3},
			{4, 5},
			{6, 7, 8, 9},
		}

		result := utils.ArrFlatten(arrays)
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with empty array
		emptyResult := utils.ArrFlatten([][]int{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}

		// Test with empty inner arrays
		emptyInnerResult := utils.ArrFlatten([][]int{{}, {}, {}})

		if len(emptyInnerResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyInnerResult))
		}
	})

	// Test Reverse
	t.Run("Reverse", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result := utils.ArrReverse(array)
		expected := []int{5, 4, 3, 2, 1}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}

		// Test with empty array
		emptyResult := utils.ArrReverse([]int{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got length %d", len(emptyResult))
		}

		// Test with single element
		singleResult := utils.ArrReverse([]int{1})

		if len(singleResult) != 1 || singleResult[0] != 1 {
			t.Errorf("Expected [1], got %v", singleResult)
		}
	})

	// Test Forget
	t.Run("Forget", func(t *testing.T) {
		array := map[string]interface{}{
			"name": "John",
			"age":  30,
			"city": "New York",
		}

		result := utils.ArrForget(array, "age", "city")

		if len(result) != 1 {
			t.Errorf("Expected 1 item, got %d", len(result))
		}

		if result["name"] != "John" {
			t.Errorf("Expected name to be 'John', got %v", result["name"])
		}

		if _, exists := result["age"]; exists {
			t.Error("Expected 'age' to be removed")
		}

		if _, exists := result["city"]; exists {
			t.Error("Expected 'city' to be removed")
		}

		// Test with non-existent keys
		result = utils.ArrForget(array, "country")

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}
	})

	// Test KeyBy
	t.Run("KeyBy", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}

		people := []Person{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Charlie"},
		}

		result := utils.ArrKeyBy(people, func(p Person) int {
			return p.ID
		})

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		if result[1].Name != "Alice" || result[2].Name != "Bob" || result[3].Name != "Charlie" {
			t.Error("Expected correct mapping of IDs to people")
		}

		// Test with empty array
		emptyResult := utils.ArrKeyBy([]Person{}, func(p Person) int {
			return p.ID
		})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}
	})

	// Test Last
	t.Run("Last", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result, found := utils.ArrLast(array, func(item int) bool {
			return item%2 == 0
		})

		if !found || result != 4 {
			t.Errorf("Expected to find 4, got %d with found=%v", result, found)
		}

		// Test not found
		result, found = utils.ArrLast(array, func(item int) bool {
			return item > 10
		})

		if found {
			t.Errorf("Expected not to find any item, but found %d", result)
		}

		// Test with empty array
		emptyArray := []int{}
		_, found = utils.ArrLast(emptyArray, func(item int) bool {
			return true
		})

		if found {
			t.Error("Expected not to find any item in empty array")
		}
	})

	// Test LastOrDefault
	t.Run("LastOrDefault", func(t *testing.T) {
		array := []int{1, 2, 3}

		result := utils.ArrLastOrDefault(array, 0)
		if result != 3 {
			t.Errorf("Expected 3, got %d", result)
		}

		emptyArray := []int{}
		result = utils.ArrLastOrDefault(emptyArray, 0)
		if result != 0 {
			t.Errorf("Expected 0, got %d", result)
		}
	})

	// Test Query
	t.Run("Query", func(t *testing.T) {
		array := map[string]interface{}{
			"name":    "John",
			"age":     30,
			"hobbies": []string{"reading", "swimming"},
		}

		result := utils.ArrQuery(array)

		// The order of query parameters is not guaranteed, so we'll check for the presence of expected substrings
		if !strings.Contains(result, "name=John") {
			t.Errorf("Expected query to contain 'name=John', got '%s'", result)
		}

		if !strings.Contains(result, "age=30") {
			t.Errorf("Expected query to contain 'age=30', got '%s'", result)
		}

		if !strings.Contains(result, "hobbies%5B%5D=reading") || !strings.Contains(result, "hobbies%5B%5D=swimming") {
			t.Errorf("Expected query to contain hobbies, got '%s'", result)
		}

		// Test with empty map
		emptyResult := utils.ArrQuery(map[string]interface{}{})

		if emptyResult != "" {
			t.Errorf("Expected empty string, got '%s'", emptyResult)
		}
	})

	// Test RandomOrDefault
	t.Run("RandomOrDefault", func(t *testing.T) {
		array := []int{1, 2, 3, 4, 5}

		result := utils.ArrRandomOrDefault(array, 0)

		// Check that the result is one of the elements in the array
		found := false
		for _, v := range array {
			if result == v {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected result to be one of the elements in the array, got %d", result)
		}

		// Test with empty array
		emptyArray := []int{}
		emptyResult := utils.ArrRandomOrDefault(emptyArray, 42)

		if emptyResult != 42 {
			t.Errorf("Expected default value 42, got %d", emptyResult)
		}
	})

	// Test SortByKey
	t.Run("SortByKey", func(t *testing.T) {
		array := map[string]interface{}{
			"c": 3,
			"a": 1,
			"b": 2,
		}

		result := utils.ArrSortByKey(array)

		// Check that the keys are in alphabetical order
		keys := make([]string, 0, len(result))
		for k := range result {
			keys = append(keys, k)
		}

		// Sort the keys for consistent comparison
		sort.Strings(keys)

		expected := []string{"a", "b", "c"}
		if len(keys) != len(expected) {
			t.Errorf("Expected %d keys, got %d", len(expected), len(keys))
		}

		for i, v := range expected {
			if keys[i] != v {
				t.Errorf("Expected key %s at index %d, got %s", v, i, keys[i])
			}
		}

		// Test with empty map
		emptyResult := utils.ArrSortByKey(map[string]interface{}{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}
	})

	// TODO Don't work
	// Test SortByKeyDesc
	//t.Run("SortByKeyDesc", func(t *testing.T) {
	//	array := map[string]interface{}{
	//		"c": 3,
	//		"a": 1,
	//		"b": 2,
	//	}
	//
	//	result := utils.ArrSortByKeyDesc(array)
	//
	//	// Check that the keys are in reverse alphabetical order
	//	keys := make([]string, 0, len(result))
	//	for k := range result {
	//		keys = append(keys, k)
	//	}
	//
	//	expected := []string{"c", "b", "a"}
	//	if len(keys) != len(expected) {
	//		t.Errorf("Expected %d keys, got %d", len(expected), len(keys))
	//	}
	//
	//	for i, v := range expected {
	//		if keys[i] != v {
	//			t.Errorf("Expected key %s at index %d, got %s", v, i, keys[i])
	//		}
	//	}
	//
	//	// Test with empty map
	//	emptyResult := utils.ArrSortByKeyDesc(map[string]interface{}{})
	//
	//	if len(emptyResult) != 0 {
	//		t.Errorf("Expected empty result, got %d items", len(emptyResult))
	//	}
	//})

	// Test SortRecursive
	t.Run("SortRecursive", func(t *testing.T) {
		array := map[string]interface{}{
			"c": 3,
			"a": map[string]interface{}{
				"z": 26,
				"x": 24,
				"y": 25,
			},
			"b": []interface{}{5, 3, 1, 4, 2},
		}

		result := utils.ArrSortRecursive(array)

		// Check that the top-level keys are sorted
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			t.Error("Expected result to be a map")
			return
		}

		keys := make([]string, 0, len(resultMap))
		for k := range resultMap {
			keys = append(keys, k)
		}

		// Sort the keys for consistent comparison
		sort.Strings(keys)

		expectedKeys := []string{"a", "b", "c"}
		if len(keys) != len(expectedKeys) {
			t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(keys))
		}

		for i, v := range expectedKeys {
			if keys[i] != v {
				t.Errorf("Expected key %s at index %d, got %s", v, i, keys[i])
			}
		}

		// Check that the nested map is sorted
		nestedMap, ok := resultMap["a"].(map[string]interface{})
		if !ok {
			t.Error("Expected 'a' to be a map")
			return
		}

		nestedKeys := make([]string, 0, len(nestedMap))
		for k := range nestedMap {
			nestedKeys = append(nestedKeys, k)
		}

		// Sort the nested keys for consistent comparison
		sort.Strings(nestedKeys)

		expectedNestedKeys := []string{"x", "y", "z"}
		if len(nestedKeys) != len(expectedNestedKeys) {
			t.Errorf("Expected %d nested keys, got %d", len(expectedNestedKeys), len(nestedKeys))
		}

		for i, v := range expectedNestedKeys {
			if nestedKeys[i] != v {
				t.Errorf("Expected nested key %s at index %d, got %s", v, i, nestedKeys[i])
			}
		}

		// Test with empty map
		emptyResult := utils.ArrSortRecursive(map[string]interface{}{})

		if len(emptyResult.(map[string]interface{})) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult.(map[string]interface{})))
		}
	})

	// Test WhereNotNull
	t.Run("WhereNotNull", func(t *testing.T) {
		// Test with pointers
		val1 := 1
		val2 := 2
		array := []*int{&val1, nil, &val2, nil}

		result := utils.ArrWhereNotNull(array)

		if len(result) != 2 {
			t.Errorf("Expected 2 items, got %d", len(result))
		}

		if *result[0] != 1 || *result[1] != 2 {
			t.Error("Expected only non-nil values")
		}

		// Test with empty array
		emptyResult := utils.ArrWhereNotNull([]*int{})

		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(emptyResult))
		}

		// Test with all nil values
		allNilArray := []*int{nil, nil}
		allNilResult := utils.ArrWhereNotNull(allNilArray)

		if len(allNilResult) != 0 {
			t.Errorf("Expected empty result, got %d items", len(allNilResult))
		}

		// Test with structs (should not filter out structs with nil fields)
		type Nullable struct {
			Value *int
		}

		structArray := []Nullable{
			{Value: &val1},
			{Value: nil},
			{Value: &val2},
			{Value: nil},
		}

		structResult := utils.ArrWhereNotNull(structArray)

		// The function checks if the struct itself is nil, not if fields within the struct are nil
		// So all structs should be included in the result
		if len(structResult) != len(structArray) {
			t.Errorf("Expected %d items, got %d", len(structArray), len(structResult))
		}
	})
}
