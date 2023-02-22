/*
 * OTUS Highload Architect
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.2.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type PostUpdatePutRequest struct {

	// Идентификатор поста
	Id string `json:"id"`

	// Текст поста
	Text string `json:"text"`
}

// AssertPostUpdatePutRequestRequired checks if the required fields are not zero-ed
func AssertPostUpdatePutRequestRequired(obj PostUpdatePutRequest) error {
	elements := map[string]interface{}{
		"id":   obj.Id,
		"text": obj.Text,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecursePostUpdatePutRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of PostUpdatePutRequest (e.g. [][]PostUpdatePutRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecursePostUpdatePutRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aPostUpdatePutRequest, ok := obj.(PostUpdatePutRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertPostUpdatePutRequestRequired(aPostUpdatePutRequest)
	})
}