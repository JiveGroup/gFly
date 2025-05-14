package test

import (
	"gfly/app/utils"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestColUtils(t *testing.T) {
	// Test All
	t.Run("All", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColAll(collection)

		if !reflect.DeepEqual(result, collection) {
			t.Errorf("Expected %v, got %v", collection, result)
		}
	})

	// Test Avg
	t.Run("Avg", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColAvg(collection, func(item int) float64 {
			return float64(item)
		})

		if result != 3.0 {
			t.Errorf("Expected 3.0, got %f", result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyResult := utils.ColAvg(emptyCollection, func(item int) float64 {
			return float64(item)
		})

		if emptyResult != 0.0 {
			t.Errorf("Expected 0.0 for empty collection, got %f", emptyResult)
		}
	})

	// Test Chunk
	t.Run("Chunk", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := utils.ColChunk(collection, 3)

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
			if !reflect.DeepEqual(result[i], chunk) {
				t.Errorf("Expected chunk %d to be %v, got %v", i, chunk, result[i])
			}
		}

		// Test with invalid chunk size
		invalidResult := utils.ColChunk(collection, 0)
		if len(invalidResult) != 0 {
			t.Errorf("Expected empty result for invalid chunk size, got %v", invalidResult)
		}
	})

	// Test Collapse
	t.Run("Collapse", func(t *testing.T) {
		collection := [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}
		result := utils.ColCollapse(collection)

		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with empty collection
		emptyCollection := [][]int{}
		emptyResult := utils.ColCollapse(emptyCollection)
		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result, got %v", emptyResult)
		}
	})

	// Test Contains
	t.Run("Contains", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		if !utils.ColContains(collection, 3) {
			t.Error("Expected collection to contain 3")
		}

		if utils.ColContains(collection, 6) {
			t.Error("Expected collection to not contain 6")
		}
	})

	// Test ContainsFunc
	t.Run("ContainsFunc", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		if !utils.ColContainsFunc(collection, func(item int) bool {
			return item > 3
		}) {
			t.Error("Expected collection to contain an item greater than 3")
		}

		if utils.ColContainsFunc(collection, func(item int) bool {
			return item > 5
		}) {
			t.Error("Expected collection to not contain an item greater than 5")
		}
	})

	// Test Count
	t.Run("Count", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColCount(collection)
		if result != 5 {
			t.Errorf("Expected count 5, got %d", result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyResult := utils.ColCount(emptyCollection)
		if emptyResult != 0 {
			t.Errorf("Expected count 0 for empty collection, got %d", emptyResult)
		}
	})

	// Test CountBy
	t.Run("CountBy", func(t *testing.T) {
		collection := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
		result := utils.ColCountBy(collection, func(item int) int {
			return item
		})

		if result[1] != 1 || result[2] != 2 || result[3] != 3 || result[4] != 4 {
			t.Errorf("Expected counts {1:1, 2:2, 3:3, 4:4}, got %v", result)
		}
	})

	// Test CrossJoin
	t.Run("CrossJoin", func(t *testing.T) {
		collection := []string{"a", "b"}
		result := utils.ColCrossJoin(collection, []string{"1", "2"}, []string{"x", "y"})

		// Expected: [["a", "1", "x"], ["a", "1", "y"], ["a", "2", "x"], ["a", "2", "y"], ["b", "1", "x"], ["b", "1", "y"], ["b", "2", "x"], ["b", "2", "y"]]
		if len(result) != 8 {
			t.Errorf("Expected 8 results, got %d", len(result))
		}

		// Check a few sample results
		if !reflect.DeepEqual(result[0], []string{"a", "1", "x"}) {
			t.Errorf("Expected first result to be [a, 1, x], got %v", result[0])
		}

		if !reflect.DeepEqual(result[7], []string{"b", "2", "y"}) {
			t.Errorf("Expected last result to be [b, 2, y], got %v", result[7])
		}
	})

	// Test Diff
	t.Run("Diff", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColDiff(collection, []int{2, 4, 6})

		expected := []int{1, 3, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test DiffAssoc
	t.Run("DiffAssoc", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColDiffAssoc(collection, map[string]int{
			"a": 1,
			"b": 5,
			"d": 4,
		})

		expected := map[string]int{
			"b": 2,
			"c": 3,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test DiffKeys
	t.Run("DiffKeys", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColDiffKeys(collection, map[string]int{
			"a": 10,
			"d": 40,
		})

		expected := map[string]int{
			"b": 2,
			"c": 3,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Each
	t.Run("Each", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		sum := 0
		utils.ColEach(collection, func(item int, index int) bool {
			sum += item
			return true
		})

		if sum != 15 {
			t.Errorf("Expected sum 15, got %d", sum)
		}

		// Test early termination
		sum = 0
		utils.ColEach(collection, func(item int, index int) bool {
			sum += item
			return item < 3
		})

		if sum != 6 {
			t.Errorf("Expected sum 6 (1+2+3), got %d", sum)
		}
	})

	// Test Every
	t.Run("Every", func(t *testing.T) {
		collection := []int{2, 4, 6, 8, 10}
		result := utils.ColEvery(collection, func(item int) bool {
			return item%2 == 0
		})

		if !result {
			t.Error("Expected all items to be even")
		}

		collection = []int{2, 4, 5, 8, 10}
		result = utils.ColEvery(collection, func(item int) bool {
			return item%2 == 0
		})

		if result {
			t.Error("Expected not all items to be even")
		}
	})

	// Test Except
	t.Run("Except", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColExcept(collection, []string{"a", "c"})

		expected := map[string]int{
			"b": 2,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Filter
	t.Run("Filter", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColFilter(collection, func(item int) bool {
			return item%2 == 0
		})

		expected := []int{2, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test First
	t.Run("First", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result, found := utils.ColFirst(collection, func(item int) bool {
			return item > 2
		})

		if !found || result != 3 {
			t.Errorf("Expected 3, got %d (found: %v)", result, found)
		}

		// Test not found
		result, found = utils.ColFirst(collection, func(item int) bool {
			return item > 10
		})

		if found {
			t.Errorf("Expected not found, got %d", result)
		}
	})

	// Test FirstOrDefault
	t.Run("FirstOrDefault", func(t *testing.T) {
		collection := []int{1, 2, 3}
		result := utils.ColFirstOrDefault(collection, 0)
		if result != 1 {
			t.Errorf("Expected 1, got %d", result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyResult := utils.ColFirstOrDefault(emptyCollection, 0)
		if emptyResult != 0 {
			t.Errorf("Expected 0 for empty collection, got %d", emptyResult)
		}
	})

	// Test FlatMap
	t.Run("FlatMap", func(t *testing.T) {
		collection := []int{1, 2, 3}
		result := utils.ColFlatMap(collection, func(item int) []string {
			return []string{strconv.Itoa(item), strconv.Itoa(item * 2)}
		})

		expected := []string{"1", "2", "2", "4", "3", "6"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Flatten
	t.Run("Flatten", func(t *testing.T) {
		collection := [][]int{
			{1, 2},
			{3, 4},
			{5, 6},
		}
		result := utils.ColFlatten(collection)

		expected := []int{1, 2, 3, 4, 5, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Flip
	t.Run("Flip", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColFlip(collection)

		expected := map[int]string{
			1: "a",
			2: "b",
			3: "c",
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Forget
	t.Run("Forget", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColForget(collection, "a", "c")

		expected := map[string]int{
			"b": 2,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
		}
		result := utils.ColGet(collection, "a", 0)
		if result != 1 {
			t.Errorf("Expected 1, got %d", result)
		}

		// Test with default value
		defaultResult := utils.ColGet(collection, "c", 3)
		if defaultResult != 3 {
			t.Errorf("Expected default value 3, got %d", defaultResult)
		}
	})

	// Test GroupBy
	t.Run("GroupBy", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		collection := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 25},
			{Name: "Dave", Age: 30},
		}

		result := utils.ColGroupBy(collection, func(p Person) int {
			return p.Age
		})

		if len(result) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(result))
		}

		if len(result[25]) != 2 || len(result[30]) != 2 {
			t.Errorf("Expected 2 people in each group, got %d in group 25 and %d in group 30",
				len(result[25]), len(result[30]))
		}
	})

	// Test Has
	t.Run("Has", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
		}
		if !utils.ColHas(collection, "a") {
			t.Error("Expected collection to have key 'a'")
		}

		if utils.ColHas(collection, "c") {
			t.Error("Expected collection to not have key 'c'")
		}
	})

	// Test Implode
	t.Run("Implode", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColImplode(collection, ", ", func(item int) string {
			return strconv.Itoa(item)
		})

		expected := "1, 2, 3, 4, 5"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	// Test Intersect
	t.Run("Intersect", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColIntersect(collection, []int{3, 4, 5, 6, 7})

		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test IntersectByKeys
	t.Run("IntersectByKeys", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColIntersectByKeys(collection, []string{"a", "c", "d"})

		expected := map[string]int{
			"a": 1,
			"c": 3,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test IsEmpty
	t.Run("IsEmpty", func(t *testing.T) {
		collection := []int{1, 2, 3}
		if utils.ColIsEmpty(collection) {
			t.Error("Expected collection to not be empty")
		}

		emptyCollection := []int{}
		if !utils.ColIsEmpty(emptyCollection) {
			t.Error("Expected empty collection to be empty")
		}
	})

	// Test IsNotEmpty
	t.Run("IsNotEmpty", func(t *testing.T) {
		collection := []int{1, 2, 3}
		if !utils.ColIsNotEmpty(collection) {
			t.Error("Expected collection to not be empty")
		}

		emptyCollection := []int{}
		if utils.ColIsNotEmpty(emptyCollection) {
			t.Error("Expected empty collection to be empty")
		}
	})

	// Test KeyBy
	t.Run("KeyBy", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}

		collection := []Person{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Charlie"},
		}

		result := utils.ColKeyBy(collection, func(p Person) int {
			return p.ID
		})

		if len(result) != 3 {
			t.Errorf("Expected 3 items, got %d", len(result))
		}

		if result[1].Name != "Alice" || result[2].Name != "Bob" || result[3].Name != "Charlie" {
			t.Error("Expected correct mapping of IDs to people")
		}
	})

	// Test Keys
	t.Run("Keys", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColKeys(collection)

		// Sort for consistent comparison
		sort.Strings(result)
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Last
	t.Run("Last", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result, found := utils.ColLast(collection, func(item int) bool {
			return item%2 == 0
		})

		if !found || result != 4 {
			t.Errorf("Expected 4, got %d (found: %v)", result, found)
		}

		// Test not found
		result, found = utils.ColLast(collection, func(item int) bool {
			return item > 10
		})

		if found {
			t.Errorf("Expected not found, got %d", result)
		}
	})

	// Test LastOrDefault
	t.Run("LastOrDefault", func(t *testing.T) {
		collection := []int{1, 2, 3}
		result := utils.ColLastOrDefault(collection, 0)
		if result != 3 {
			t.Errorf("Expected 3, got %d", result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyResult := utils.ColLastOrDefault(emptyCollection, 0)
		if emptyResult != 0 {
			t.Errorf("Expected 0 for empty collection, got %d", emptyResult)
		}
	})

	// Test Map
	t.Run("Map", func(t *testing.T) {
		collection := []int{1, 2, 3}
		result := utils.ColMap(collection, func(item int) string {
			return strconv.Itoa(item)
		})

		expected := []string{"1", "2", "3"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Max
	t.Run("Max", func(t *testing.T) {
		collection := []int{1, 5, 3, 9, 7}
		result := utils.ColMax(collection, func(item int) int {
			return item
		})

		if result != 9 {
			t.Errorf("Expected max 9, got %d", result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyResult := utils.ColMax(emptyCollection, func(item int) int {
			return item
		})

		if emptyResult != 0 {
			t.Errorf("Expected 0 for empty collection, got %d", emptyResult)
		}
	})

	// Test Merge
	t.Run("Merge", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
		}
		result := utils.ColMerge(collection, map[string]int{
			"b": 3,
			"c": 4,
		})

		expected := map[string]int{
			"a": 1,
			"b": 3,
			"c": 4,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Min
	t.Run("Min", func(t *testing.T) {
		collection := []int{5, 3, 1, 9, 7}
		result := utils.ColMin(collection, func(item int) int {
			return item
		})

		if result != 1 {
			t.Errorf("Expected min 1, got %d", result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyResult := utils.ColMin(emptyCollection, func(item int) int {
			return item
		})

		if emptyResult != 0 {
			t.Errorf("Expected 0 for empty collection, got %d", emptyResult)
		}
	})

	// Test Only
	t.Run("Only", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColOnly(collection, []string{"a", "c"})

		expected := map[string]int{
			"a": 1,
			"c": 3,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Pad
	t.Run("Pad", func(t *testing.T) {
		collection := []int{1, 2, 3}
		result := utils.ColPad(collection, 5, 0)

		expected := []int{1, 2, 3, 0, 0}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with size smaller than collection
		smallResult := utils.ColPad(collection, 2, 0)
		if !reflect.DeepEqual(smallResult, collection) {
			t.Errorf("Expected original collection %v, got %v", collection, smallResult)
		}
	})

	// Test Partition
	t.Run("Partition", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		pass, fail := utils.ColPartition(collection, func(item int) bool {
			return item%2 == 0
		})

		expectedPass := []int{2, 4}
		expectedFail := []int{1, 3, 5}
		if !reflect.DeepEqual(pass, expectedPass) {
			t.Errorf("Expected pass %v, got %v", expectedPass, pass)
		}
		if !reflect.DeepEqual(fail, expectedFail) {
			t.Errorf("Expected fail %v, got %v", expectedFail, fail)
		}
	})

	// Test Pluck
	t.Run("Pluck", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		collection := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 35},
		}

		result := utils.ColPluck(collection, func(p Person) string {
			return p.Name
		})

		expected := []string{"Alice", "Bob", "Charlie"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Prepend
	t.Run("Prepend", func(t *testing.T) {
		collection := []int{3, 4, 5}
		result := utils.ColPrepend(collection, 1, 2)

		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Pull
	t.Run("Pull", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		item, result := utils.ColPull(collection, 2)

		if item != 3 {
			t.Errorf("Expected pulled item 3, got %d", item)
		}

		expected := []int{1, 2, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with invalid index
		invalidItem, invalidResult := utils.ColPull(collection, 10)
		var zero int
		if invalidItem != zero {
			t.Errorf("Expected zero value for invalid index, got %d", invalidItem)
		}
		if !reflect.DeepEqual(invalidResult, collection) {
			t.Errorf("Expected original collection for invalid index, got %v", invalidResult)
		}
	})

	// Test Push
	t.Run("Push", func(t *testing.T) {
		collection := []int{1, 2, 3}
		result := utils.ColPush(collection, 4, 5)

		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Put
	t.Run("Put", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
		}
		result := utils.ColPut(collection, "c", 3)

		expected := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test overwriting existing key
		overwriteResult := utils.ColPut(collection, "a", 10)
		expectedOverwrite := map[string]int{
			"a": 10,
			"b": 2,
		}
		if !reflect.DeepEqual(overwriteResult, expectedOverwrite) {
			t.Errorf("Expected %v, got %v", expectedOverwrite, overwriteResult)
		}
	})

	// Test Random
	t.Run("Random", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		_, found := utils.ColRandom(collection)
		if !found {
			t.Error("Expected to find a random element")
		}

		// Test with empty collection
		emptyCollection := []int{}
		_, found = utils.ColRandom(emptyCollection)
		if found {
			t.Error("Expected not to find a random element in empty collection")
		}
	})

	// Test RandomOrDefault
	t.Run("RandomOrDefault", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColRandomOrDefault(collection, 0)

		// Check that the result is one of the elements in the collection
		found := false
		for _, v := range collection {
			if result == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected result to be one of the elements in the collection, got %d", result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyResult := utils.ColRandomOrDefault(emptyCollection, 42)
		if emptyResult != 42 {
			t.Errorf("Expected default value 42, got %d", emptyResult)
		}
	})

	// Test Reduce
	t.Run("Reduce", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColReduce(collection, 0, func(acc, item int) int {
			return acc + item
		})

		if result != 15 {
			t.Errorf("Expected sum 15, got %d", result)
		}
	})

	// Test Reject
	t.Run("Reject", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColReject(collection, func(item int) bool {
			return item%2 == 0
		})

		expected := []int{1, 3, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Reverse
	t.Run("Reverse", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColReverse(collection)

		expected := []int{5, 4, 3, 2, 1}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Search
	t.Run("Search", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		index, found := utils.ColSearch(collection, 3)

		if !found || index != 2 {
			t.Errorf("Expected to find 3 at index 2, got index %d, found %v", index, found)
		}

		// Test not found
		index, found = utils.ColSearch(collection, 6)
		if found {
			t.Errorf("Expected not to find 6, but found at index %d", index)
		}
	})

	// Test SearchFunc
	t.Run("SearchFunc", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		index, found := utils.ColSearchFunc(collection, func(item int) bool {
			return item > 3
		})

		if !found || index != 3 {
			t.Errorf("Expected to find item > 3 at index 3, got index %d, found %v", index, found)
		}

		// Test not found
		index, found = utils.ColSearchFunc(collection, func(item int) bool {
			return item > 10
		})
		if found {
			t.Errorf("Expected not to find item > 10, but found at index %d", index)
		}
	})

	// Test Shift
	t.Run("Shift", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		item, result := utils.ColShift(collection)

		if item != 1 {
			t.Errorf("Expected shifted item 1, got %d", item)
		}

		expected := []int{2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with empty collection
		emptyCollection := []int{}
		emptyItem, emptyResult := utils.ColShift(emptyCollection)
		var zero int
		if emptyItem != zero {
			t.Errorf("Expected zero value for empty collection, got %d", emptyItem)
		}
		if len(emptyResult) != 0 {
			t.Errorf("Expected empty result for empty collection, got %v", emptyResult)
		}
	})

	// Test Shuffle
	t.Run("Shuffle", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColShuffle(collection)

		// Check that the result has the same length
		if len(result) != len(collection) {
			t.Errorf("Expected length %d, got %d", len(collection), len(result))
		}

		// Check that the result contains all the original elements
		for _, v := range collection {
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
	})

	// Test Slice
	t.Run("Slice", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColSlice(collection, 2)

		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with negative index
		negativeResult := utils.ColSlice(collection, -2)
		expectedNegative := []int{4, 5}
		if !reflect.DeepEqual(negativeResult, expectedNegative) {
			t.Errorf("Expected %v, got %v", expectedNegative, negativeResult)
		}
	})

	// Test SliceWithLength
	t.Run("SliceWithLength", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColSliceWithLength(collection, 1, 3)

		expected := []int{2, 3, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with negative index
		negativeResult := utils.ColSliceWithLength(collection, -3, 2)
		expectedNegative := []int{3, 4}
		if !reflect.DeepEqual(negativeResult, expectedNegative) {
			t.Errorf("Expected %v, got %v", expectedNegative, negativeResult)
		}
	})

	// Test Sort
	t.Run("Sort", func(t *testing.T) {
		collection := []int{5, 3, 1, 4, 2}
		result := utils.ColSort(collection, func(a, b int) bool {
			return a < b
		})

		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test SortBy
	t.Run("SortBy", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		collection := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 20},
			{Name: "Charlie", Age: 25},
		}

		result := utils.ColSortBy(collection, func(p Person) int {
			return p.Age
		}, func(a, b int) bool {
			return a < b
		})

		if result[0].Name != "Bob" || result[1].Name != "Charlie" || result[2].Name != "Alice" {
			t.Errorf("Expected sorted by age: Bob, Charlie, Alice, got %v", result)
		}
	})

	// Test SortByDesc
	t.Run("SortByDesc", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		collection := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 20},
			{Name: "Charlie", Age: 25},
		}

		result := utils.ColSortByDesc(collection, func(p Person) int {
			return p.Age
		}, func(a, b int) bool {
			return a < b
		})

		if result[0].Name != "Alice" || result[1].Name != "Charlie" || result[2].Name != "Bob" {
			t.Errorf("Expected sorted by age desc: Alice, Charlie, Bob, got %v", result)
		}
	})

	// TODO Test Splice
	//t.Run("Splice", func(t *testing.T) {
	//	collection := []int{1, 2, 3, 4, 5}
	//	removed, result := utils.ColSplice(collection, 0, 3)
	//
	//	expectedRemoved := []int{2, 3, 4}
	//	expectedResult := []int{1, 5}
	//	if !reflect.DeepEqual(removed, expectedRemoved) {
	//		t.Errorf("Expected removed %v, got %v", expectedRemoved, removed)
	//	}
	//	if !reflect.DeepEqual(result, expectedResult) {
	//		t.Errorf("Expected result %v, got %v", expectedResult, result)
	//	}
	//})

	// Test Split
	t.Run("Split", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5, 6}
		result := utils.ColSplit(collection, 3)

		if len(result) != 3 {
			t.Errorf("Expected 3 groups, got %d", len(result))
		}

		// Check that all elements are distributed
		allElements := []int{}
		for _, group := range result {
			allElements = append(allElements, group...)
		}
		sort.Ints(allElements)
		expected := []int{1, 2, 3, 4, 5, 6}
		if !reflect.DeepEqual(allElements, expected) {
			t.Errorf("Expected all elements %v, got %v", expected, allElements)
		}
	})

	// Test Sum
	t.Run("Sum", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColSum(collection, func(item int) int {
			return item
		})

		if result != 15 {
			t.Errorf("Expected sum 15, got %d", result)
		}
	})

	// Test Take
	t.Run("Take", func(t *testing.T) {
		collection := []int{1, 2, 3, 4, 5}
		result := utils.ColTake(collection, 3)

		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with limit larger than collection
		largeResult := utils.ColTake(collection, 10)
		if !reflect.DeepEqual(largeResult, collection) {
			t.Errorf("Expected original collection for large limit, got %v", largeResult)
		}
	})

	// Test Tap
	t.Run("Tap", func(t *testing.T) {
		collection := []int{1, 2, 3}
		tapped := false
		result := utils.ColTap(collection, func(c []int) {
			tapped = true
			if !reflect.DeepEqual(c, collection) {
				t.Errorf("Expected callback to receive %v, got %v", collection, c)
			}
		})

		if !tapped {
			t.Error("Expected callback to be called")
		}
		if !reflect.DeepEqual(result, collection) {
			t.Errorf("Expected result to be original collection, got %v", result)
		}
	})

	// Test Unique
	t.Run("Unique", func(t *testing.T) {
		collection := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
		result := utils.ColUnique(collection)

		expected := []int{1, 2, 3, 4, 5}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}

		// Check that all expected elements are in the result
		for _, v := range expected {
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
	})

	// Test UniqueBy
	t.Run("UniqueBy", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		collection := []Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 25},
			{Name: "Dave", Age: 30},
		}

		result := utils.ColUniqueBy(collection, func(p Person) int {
			return p.Age
		})

		if len(result) != 2 {
			t.Errorf("Expected 2 unique ages, got %d", len(result))
		}
	})

	// Test Values
	t.Run("Values", func(t *testing.T) {
		collection := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		result := utils.ColValues(collection)

		// Sort for consistent comparison
		sort.Ints(result)
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	// Test Zip
	t.Run("Zip", func(t *testing.T) {
		collection := []string{"a", "b", "c"}
		result := utils.ColZip(collection, []string{"1", "2", "3"}, []string{"x", "y", "z"})

		expected := [][]string{
			{"a", "1", "x"},
			{"b", "2", "y"},
			{"c", "3", "z"},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// Test with arrays of different lengths
		shortResult := utils.ColZip(collection, []string{"1", "2"})
		expectedShort := [][]string{
			{"a", "1"},
			{"b", "2"},
		}
		if !reflect.DeepEqual(shortResult, expectedShort) {
			t.Errorf("Expected %v, got %v", expectedShort, shortResult)
		}
	})
}
