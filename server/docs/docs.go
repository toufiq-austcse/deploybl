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
        "/api/v1/deployments": {
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/res.DeploymentRes"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
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
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.DeploymentRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/check-deploying-cron": {
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
                "summary": "Check Deploying state Deployments",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/check-stopped-cron": {
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
                "summary": "Check Stopped Deployments",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/latest-status": {
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
                "summary": "Deployments Latest Status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Deployment ID",
                        "name": "ids",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/res.DeploymentLatestStatusRes"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/{id}": {
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.DeploymentDetailsRes"
                                        }
                                    }
                                }
                            ]
                        }
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.DeploymentDetailsRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/{id}/env": {
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.DeploymentDetailsRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/{id}/events": {
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
                "summary": "Deployment Events",
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/res.EventRes"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/{id}/rebuild-and-redeploy": {
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
                "summary": "Rebuild and Deploy Deployment",
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.DeploymentDetailsRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/{id}/restart": {
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
                "summary": "Restart Deployment",
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.DeploymentDetailsRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/deployments/{id}/stop": {
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
                "summary": "Stop Deployment",
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.DeploymentDetailsRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/repositories": {
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/res.RepoDetailsRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/repositories/branches": {
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
                "summary": "Get Repo Branches",
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
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api_response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/res.RepoBranchRes"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api_response.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "This is Name",
                    "type": "integer"
                },
                "data": {},
                "errors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
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
        },
        "res.DeploymentDetailsRes": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "branch_name": {
                    "type": "string"
                },
                "container_id": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "docker_file_path": {
                    "type": "string"
                },
                "docker_image_tag": {
                    "type": "string"
                },
                "domain_url": {
                    "type": "string"
                },
                "env": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "last_deployed_at": {
                    "type": "string"
                },
                "latest_status": {
                    "type": "string"
                },
                "repository_name": {
                    "type": "string"
                },
                "repository_provider": {
                    "type": "string"
                },
                "repository_url": {
                    "type": "string"
                },
                "root_directory": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "res.DeploymentLatestStatusRes": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "domain_url": {
                    "type": "string"
                },
                "last_deployed_at": {
                    "type": "string"
                },
                "latest_status": {
                    "type": "string"
                }
            }
        },
        "res.DeploymentRes": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "branch_name": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "domain_url": {
                    "type": "string"
                },
                "last_deployed_at": {
                    "type": "string"
                },
                "latest_status": {
                    "type": "string"
                },
                "repository_provider": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "res.EventRes": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deployment_id": {
                    "type": "string"
                },
                "event_log_file_url": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "triggered_by": {
                    "type": "string"
                },
                "triggered_value": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "res.RepoBranchRes": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "res.RepoDetailsRes": {
            "type": "object",
            "properties": {
                "default_branch": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "svn_url": {
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
