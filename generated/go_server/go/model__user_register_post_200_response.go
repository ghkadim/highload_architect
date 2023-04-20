/*
 * OTUS Highload Architect
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.2.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type UserRegisterPost200Response struct {

	UserId string `json:"user_id,omitempty"`
}

// AssertUserRegisterPost200ResponseRequired checks if the required fields are not zero-ed
func AssertUserRegisterPost200ResponseRequired(obj UserRegisterPost200Response) error {
	return nil
}

// AssertRecurseUserRegisterPost200ResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of UserRegisterPost200Response (e.g. [][]UserRegisterPost200Response), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseUserRegisterPost200ResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aUserRegisterPost200Response, ok := obj.(UserRegisterPost200Response)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertUserRegisterPost200ResponseRequired(aUserRegisterPost200Response)
	})
}
