package transformers

import "gfly/app/utils"

// ToListResponse generic function takes a list of records, and their transformer function,
// process then return a slice of response data
func ToListResponse[T any, R any](records []T, transformerFn func(T) R) []R {
	return utils.TransformList(records, transformerFn)
}
