/*
 * OTUS Highload Architect
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.2.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DialogUserIdSendPost200Response struct {

	// Идентификатор сообщения
	MessageId string `json:"message_id,omitempty"`
}

// AssertDialogUserIdSendPost200ResponseRequired checks if the required fields are not zero-ed
func AssertDialogUserIdSendPost200ResponseRequired(obj DialogUserIdSendPost200Response) error {
	return nil
}

// AssertRecurseDialogUserIdSendPost200ResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DialogUserIdSendPost200Response (e.g. [][]DialogUserIdSendPost200Response), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDialogUserIdSendPost200ResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDialogUserIdSendPost200Response, ok := obj.(DialogUserIdSendPost200Response)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDialogUserIdSendPost200ResponseRequired(aDialogUserIdSendPost200Response)
	})
}