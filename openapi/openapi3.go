package openapi

import (
	"net/http"

	"github.com/stevenferrer/invitesvc/authn"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// NewOpenAPI3 returns the OpenAPI3 spec for Invite Service API
func NewOpenAPI3() openapi3.T {
	spec := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Invite Service REST API",
			Description: "REST API endpoints used for generating invite tokens for the Catalyst Experience App.",
			Version:     "0.1.0",
		},
		Servers: openapi3.Servers{
			{
				Description: "Local development",
				URL:         "http://localhost:8000",
			},
		},
	}

	spec.Tags = openapi3.Tags{
		{
			Name:        "Admin",
			Description: "APIs used for generating and managing invite tokens.",
		},
		{
			Name:        "Public",
			Description: "Publicly accessible APIs used for redeeming tokens.",
		},
	}

	spec.Components.SecuritySchemes = openapi3.SecuritySchemes{
		"auth_key": &openapi3.SecuritySchemeRef{
			Value: openapi3.NewSecurityScheme().
				WithDescription("Authenticate with auth key").
				WithType("apiKey").
				WithName(authn.AuthKeyHeader).
				WithIn("header"),
		},
	}

	spec.Components.Schemas = openapi3.Schemas{
		"TokenString": openapi3.NewSchemaRef("", openapi3.NewStringSchema().
			WithLength(12).WithDefault("VxzUfkY36YQT")),
		"AuthKey": openapi3.NewSchemaRef("", openapi3.NewStringSchema().
			WithLength(32).WithDefault("0d8ee59c4c1f4571a61a887b28ef7612")),
		"Token": openapi3.NewSchemaRef("",
			openapi3.NewObjectSchema().
				WithPropertyRef("token", &openapi3.SchemaRef{
					Ref: "#/components/schemas/TokenString",
				}).
				WithProperty("redeemed", openapi3.NewBoolSchema()).
				WithProperty("expiration", openapi3.NewDateTimeSchema()).
				WithProperty("disabled", openapi3.NewBoolSchema())),
		"Tokens": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: "array",
				Items: &openapi3.SchemaRef{
					Ref: "#/components/schemas/Token",
				},
			},
		},
	}

	spec.Components.Responses = openapi3.Responses{
		"Error500Response": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Internal server error").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("message", openapi3.NewStringSchema()))),
		},

		"Error404Response": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Not found error").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("message", openapi3.NewStringSchema()))),
		},

		"Error422Response": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Unprocessable entity error").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("message", openapi3.NewStringSchema()))),
		},

		"Error429Response": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Too many request error").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("message", openapi3.NewStringSchema()))),
		},

		"GenerateTokenResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Generate token response").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("token", &openapi3.SchemaRef{
						Ref: "#/components/schemas/TokenString",
					}))),
		},

		"ListTokensResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("List tokens response").
				WithContent(openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
					Ref: "#/components/schemas/Tokens",
				})),
		},

		"GetTokenResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Get token response").
				WithContent(openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
					Ref: "#/components/schemas/Token",
				})),
		},

		"GenerateAuthKeyResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Generate auth key response").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("authKey", &openapi3.SchemaRef{
						Ref: "#/components/schemas/AuthKey",
					}))),
		},

		"RedeemTokenResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Redeem token response").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("message", openapi3.NewStringSchema().
						WithDefault("token successfully redeemed.")))),
		},

		"DisableTokenResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Disable token response").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("message", openapi3.NewStringSchema().
						WithDefault("token successfully disabled.")))),
		},
	}

	spec.Paths = openapi3.Paths{
		"/admin/authkey": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "GenerateAuthKey",
				Summary:     "Generate auth key",
				Description: "Generate auth key for accessing invite service API.",
				Responses: openapi3.Responses{
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/GenerateAuthKeyResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error500Response",
					},
				},
				Security: openapi3.NewSecurityRequirements().
					With(openapi3.NewSecurityRequirement().
						Authenticate("auth_key")),
				Tags: []string{"Admin"},
			},
		},

		"/admin/tokens": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "GenerateToken",
				Summary:     "Generate invite token",
				Description: "Generate invite tokens and share to your customers.",
				Responses: openapi3.Responses{
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/GenerateTokenResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error500Response",
					},
				},
				Security: openapi3.NewSecurityRequirements().
					With(openapi3.NewSecurityRequirement().
						Authenticate("auth_key")),
				Tags: []string{"Admin"},
			},

			Get: &openapi3.Operation{
				OperationID: "ListTokens",
				Summary:     "List invite tokens",
				Description: "Retrieve list of invite tokens.",
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/ListTokensResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error500Response",
					},
				},
				Security: openapi3.NewSecurityRequirements().
					With(openapi3.NewSecurityRequirement().
						Authenticate("auth_key")),
				Tags: []string{"Admin"},
			},
		},

		"/admin/tokens/{token}": &openapi3.PathItem{
			Get: &openapi3.Operation{
				OperationID: "GetToken",
				Summary:     "Retrieve invite token",
				Description: "Retrieve invite token details.",
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/GetTokenResponse",
					},
					"404": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error404Response",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error500Response",
					},
				},
				Security: openapi3.NewSecurityRequirements().
					With(openapi3.NewSecurityRequirement().
						Authenticate("auth_key")),
				Tags: []string{"Admin"},
			},
		},

		"/admin/tokens/{token}/disable": &openapi3.PathItem{
			Put: &openapi3.Operation{
				OperationID: "DisableToken",
				Summary:     "Disable invite token",
				Description: "Disable an invite token.",
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/DisableTokenResponse",
					},
					"404": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error404Response",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error500Response",
					},
				},
				Security: openapi3.NewSecurityRequirements().
					With(openapi3.NewSecurityRequirement().
						Authenticate("auth_key")),
				Tags: []string{"Admin"},
			},
		},

		"/tokens/{token}/redeem": &openapi3.PathItem{
			Put: &openapi3.Operation{
				OperationID: "RedeemToken",
				Summary:     "Redeem invite token",
				Description: "Redeem an invite token.",
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/RedeemTokenResponse",
					},
					"404": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error404Response",
					},
					"422": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error422Response",
					},
					"429": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error429Response",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/Error500Response",
					},
				},
				Tags: []string{"Public"},
			},
		},
	}

	return spec
}

// InitOpenAPI3 initializes route for openapi spec
func InitOpenAPI3Routes(e *echo.Echo) {
	spec := NewOpenAPI3()
	g := e.Group("/openapi3")
	g.GET("/spec.json", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, spec, " ")
	})
}
