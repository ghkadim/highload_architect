/*
 * OTUS Highload Architect
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.2.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DialogUserIdSendPostRequest struct {

	// Текст сообщения
	Text string `json:"text"`
}

// AssertDialogUserIdSendPostRequestRequired checks if the required fields are not zero-ed
func AssertDialogUserIdSendPostRequestRequired(obj DialogUserIdSendPostRequest) error {
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

// AssertRecurseDialogUserIdSendPostRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DialogUserIdSendPostRequest (e.g. [][]DialogUserIdSendPostRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDialogUserIdSendPostRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDialogUserIdSendPostRequest, ok := obj.(DialogUserIdSendPostRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDialogUserIdSendPostRequestRequired(aDialogUserIdSendPostRequest)
	})
}
