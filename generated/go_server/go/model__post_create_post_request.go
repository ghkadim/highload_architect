/*
 * OTUS Highload Architect
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.2.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type PostCreatePostRequest struct {

	// Текст поста
	Text string `json:"text"`
}

// AssertPostCreatePostRequestRequired checks if the required fields are not zero-ed
func AssertPostCreatePostRequestRequired(obj PostCreatePostRequest) error {
	elements := map[string]interface{}{
		"text": obj.Text,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecursePostCreatePostRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of PostCreatePostRequest (e.g. [][]PostCreatePostRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecursePostCreatePostRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aPostCreatePostRequest, ok := obj.(PostCreatePostRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertPostCreatePostRequestRequired(aPostCreatePostRequest)
	})
}
