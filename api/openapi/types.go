// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package openapi

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// PostApiV1SecretJSONBody defines parameters for PostApiV1Secret.
type PostApiV1SecretJSONBody struct {
	Body *string                 `json:"body,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

// PostApiV1SecretJSONRequestBody defines body for PostApiV1Secret for application/json ContentType.
type PostApiV1SecretJSONRequestBody PostApiV1SecretJSONBody