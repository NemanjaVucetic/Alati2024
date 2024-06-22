// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/configGroups/": {
            "get": {
                "description": "Retrieves all configuration groups",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configGroups"
                ],
                "summary": "Get all configuration groups",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.ConfigGroup"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new configuration group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configGroups"
                ],
                "summary": "Add a new configuration group",
                "parameters": [
                    {
                        "description": "Configuration group to add",
                        "name": "configGroup",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ConfigGroup"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.ConfigGroup"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "415": {
                        "description": "Unsupported media type",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/configGroups/{nameG}/{versionG}/configs/{labels}/{nameC}/{versionC}": {
            "get": {
                "description": "Retrieves configurations by labels within a group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configGroups"
                ],
                "summary": "Get configurations by labels",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the configuration group",
                        "name": "nameG",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration group",
                        "name": "versionG",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Labels of the configuration",
                        "name": "labels",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the configuration",
                        "name": "nameC",
                        "in": "path"
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration",
                        "name": "versionC",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Config"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "description": "Deletes configurations by labels within a group",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configGroups"
                ],
                "summary": "Delete configurations by labels",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the configuration group",
                        "name": "nameG",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration group",
                        "name": "versionG",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Labels of the configuration",
                        "name": "labels",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the configuration",
                        "name": "nameC",
                        "in": "path"
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration",
                        "name": "versionC",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/configGroups/{nameG}/{versionG}/configs/{nameC}/{versionC}": {
            "put": {
                "description": "Removes a configuration from a group by their names and versions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configGroups"
                ],
                "summary": "Remove a configuration from a group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the configuration group",
                        "name": "nameG",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration group",
                        "name": "versionG",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the configuration",
                        "name": "nameC",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration",
                        "name": "versionC",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success Put",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/configGroups/{name}/{version}": {
            "get": {
                "description": "Retrieves a configuration group by name and version",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configGroups"
                ],
                "summary": "Get a configuration group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the configuration group",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration group",
                        "name": "version",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ConfigGroup"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Configuration group not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a configuration group by name and version",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configGroups"
                ],
                "summary": "Delete a configuration group",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the configuration group",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration group",
                        "name": "version",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/configs/": {
            "get": {
                "description": "Retrieves all configurations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Get all configurations",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Config"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new configuration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Add a new configuration",
                "parameters": [
                    {
                        "description": "Configuration to add",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Config"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Config"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "415": {
                        "description": "Unsupported media type",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/configs/{name}/{version}": {
            "get": {
                "description": "Retrieves a configuration by name and version",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Get a configuration",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the configuration",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration",
                        "name": "version",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Config"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Configuration not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a configuration by name and version",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Delete a configuration",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the configuration",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Version of the configuration",
                        "name": "version",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Config": {
            "type": "object",
            "properties": {
                "labels": {
                    "description": "Labels are key-value pairs for configuration\nRequired: true",
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "description": "Name of the configuration\nRequired: true",
                    "type": "string"
                },
                "params": {
                    "description": "Params are key-value pairs for configuration\nRequired: true",
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "version": {
                    "description": "Version of the configuration\nRequired: true",
                    "type": "integer"
                }
            }
        },
        "model.ConfigGroup": {
            "type": "object",
            "properties": {
                "configs": {
                    "description": "Configs in the group\nRequired: true",
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/model.Config"
                    }
                },
                "name": {
                    "description": "Name of the configuration group\nRequired: true",
                    "type": "string"
                },
                "version": {
                    "description": "Version of the configuration group\nRequired: true",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Configuration API",
	Description:      "This is a sample server for a configuration service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}