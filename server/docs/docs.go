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
        "/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Index"
                ],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/deployments": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Deployments"
                ],
                "summary": "Deployment Index",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Deployments"
                ],
                "summary": "Create Deployment",
                "parameters": [
                    {
                        "description": "Create Deployment Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.CreateDeploymentReqDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/deployments/:id/env": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Deployments"
                ],
                "summary": "Update Deployment Env",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Deployment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/deployments/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Deployments"
                ],
                "summary": "Show Deployment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Deployment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Deployments"
                ],
                "summary": "Update Deployment",
                "parameters": [
                    {
                        "description": "Update Deployment Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.UpdateDeploymentReqDto"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Deployment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/repositories": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Repositories"
                ],
                "summary": "Get Repo Details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Repo Url",
                        "name": "repo_url",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "req.CreateDeploymentReqDto": {
            "type": "object",
            "required": [
                "branch_name",
                "repository_url",
                "title"
            ],
            "properties": {
                "branch_name": {
                    "type": "string"
                },
                "docker_file_path": {
                    "type": "string"
                },
                "env": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "repository_url": {
                    "type": "string"
                },
                "root_dir": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "req.UpdateDeploymentReqDto": {
            "type": "object",
            "properties": {
                "branch_name": {
                    "type": "string"
                },
                "docker_file_path": {
                    "type": "string"
                },
                "root_dir": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
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
