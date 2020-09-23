package types

type SwaggerType string

const (
	SWAGGER_TYPE__REF   SwaggerType = "$ref"
	SWAGGER_TYPE__ITEMS             = "items"
	SWAGGER_TYPE__TYPE              = "type"
)

type SwaggerPropertyType string

const (
	SWAGGER_PROPERTY_TYPE__OBJECT  SwaggerPropertyType = "object"
	SWAGGER_PROPERTY_TYPE__ARRAY                       = "array"
	SWAGGER_PROPERTY_TYPE__STRING                      = "string"
	SWAGGER_PROPERTY_TYPE__INTEGER                     = "integer"
	SWAGGER_PROPERTY_TYPE__BOOLEAN                     = "boolean"
)
