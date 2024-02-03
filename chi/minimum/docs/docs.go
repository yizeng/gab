// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/countries/sum-population-by-state": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "countries"
                ],
                "summary": "Sum the total population by state",
                "parameters": [
                    {
                        "description": "request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SumPopulationByState"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.State"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Err"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.State": {
            "type": "object",
            "required": [
                "name",
                "population"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "population": {
                    "type": "integer"
                }
            }
        },
        "request.SumPopulationByState": {
            "type": "object",
            "required": [
                "states"
            ],
            "properties": {
                "states": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.State"
                    }
                }
            }
        },
        "response.Err": {
            "type": "object",
            "properties": {
                "error": {
                    "description": "user-facing error message",
                    "type": "string"
                },
                "error_code": {
                    "description": "application-specific error code",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
