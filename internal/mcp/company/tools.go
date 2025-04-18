package company

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twmcp "github.com/rafaeljusto/teamwork-ai/internal/mcp"
	twcompany "github.com/rafaeljusto/teamwork-ai/internal/teamwork/company"
)

func registerTools(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddTool(
		mcp.NewTool("retrieve-companies",
			mcp.WithDescription("Retrieve multiple companies, also know as clients, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
			mcp.WithString("search-term",
				mcp.Description("A search term to filter companies by name."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to filter companies by tags"),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number for pagination of results."),
			),
			mcp.WithNumber("page-size",
				mcp.Description("Number of results per page for pagination."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var multiple twcompany.Multiple

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.OptionalParam(&multiple.Filters.SearchTerm, "search-term"),
				twmcp.OptionalNumericListParam(&multiple.Filters.TagIDs, "tag-ids"),
				twmcp.OptionalNumericParam(&multiple.Filters.Page, "page"),
				twmcp.OptionalNumericParam(&multiple.Filters.PageSize, "page-size"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(multiple.Result)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-company",
			mcp.WithDescription("Retrieve a specific company, also know as client, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
			mcp.WithNumber("company-id",
				mcp.Required(),
				mcp.Description("The ID of the company."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var single twcompany.Single

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&single.ID, "company-id"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &single); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(single)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(encoded)), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("create-company",
			mcp.WithDescription("Create a new company, also know as client, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the company."),
			),
			mcp.WithString("address-one",
				mcp.Description("The first line of the address of the company."),
			),
			mcp.WithString("address-two",
				mcp.Description("The second line of the address of the company."),
			),
			mcp.WithString("city",
				mcp.Description("The city of the company."),
			),
			mcp.WithString("state",
				mcp.Description("The state of the company."),
			),
			mcp.WithString("zip",
				mcp.Description("The ZIP or postal code of the company."),
			),
			mcp.WithString("country-code",
				mcp.Description("The country code of the company, e.g., 'US' for the United States."),
			),
			mcp.WithString("phone",
				mcp.Description("The phone number of the company."),
			),
			mcp.WithString("fax",
				mcp.Description("The fax number of the company."),
			),
			mcp.WithString("email-one",
				mcp.Description("The primary email address of the company."),
			),
			mcp.WithString("email-two",
				mcp.Description("The secondary email address of the company."),
			),
			mcp.WithString("email-three",
				mcp.Description("The tertiary email address of the company."),
			),
			mcp.WithString("website",
				mcp.Description("The website of the company."),
			),
			mcp.WithString("profile",
				mcp.Description("A profile description for the company."),
			),
			mcp.WithNumber("manager-id",
				mcp.Description("The ID of the user who manages the company."),
			),
			mcp.WithNumber("currency-id",
				mcp.Description("The ID of the currency used by the company."),
			),
			mcp.WithNumber("industry-id",
				mcp.Description("The ID of the industry the company belongs to."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the company."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var company twcompany.Creation

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredParam(&company.Name, "name"),
				twmcp.OptionalPointerParam(&company.AddressOne, "address-one"),
				twmcp.OptionalPointerParam(&company.AddressTwo, "address-two"),
				twmcp.OptionalPointerParam(&company.City, "city"),
				twmcp.OptionalPointerParam(&company.State, "state"),
				twmcp.OptionalPointerParam(&company.Zip, "zip"),
				twmcp.OptionalPointerParam(&company.CountryCode, "country-code"),
				twmcp.OptionalPointerParam(&company.Phone, "phone"),
				twmcp.OptionalPointerParam(&company.Fax, "fax"),
				twmcp.OptionalPointerParam(&company.EmailOne, "email-one"),
				twmcp.OptionalPointerParam(&company.EmailTwo, "email-two"),
				twmcp.OptionalPointerParam(&company.EmailThree, "email-three"),
				twmcp.OptionalPointerParam(&company.Website, "website"),
				twmcp.OptionalPointerParam(&company.Profile, "profile"),
				twmcp.OptionalNumericPointerParam(&company.ManagerID, "manager-id"),
				twmcp.OptionalNumericPointerParam(&company.CurrencyID, "currency-id"),
				twmcp.OptionalNumericPointerParam(&company.IndustryID, "industry-id"),
				twmcp.OptionalNumericListParam(&company.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &company); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Company created successfully"), nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("update-company",
			mcp.WithDescription("Update a company, also know as client, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
			mcp.WithNumber("company-id",
				mcp.Required(),
				mcp.Description("The ID of the company to update."),
			),
			mcp.WithString("name",
				mcp.Description("The name of the company."),
			),
			mcp.WithString("address-one",
				mcp.Description("The first line of the address of the company."),
			),
			mcp.WithString("address-two",
				mcp.Description("The second line of the address of the company."),
			),
			mcp.WithString("city",
				mcp.Description("The city of the company."),
			),
			mcp.WithString("state",
				mcp.Description("The state of the company."),
			),
			mcp.WithString("zip",
				mcp.Description("The ZIP or postal code of the company."),
			),
			mcp.WithString("country-code",
				mcp.Description("The country code of the company, e.g., 'US' for the United States."),
			),
			mcp.WithString("phone",
				mcp.Description("The phone number of the company."),
			),
			mcp.WithString("fax",
				mcp.Description("The fax number of the company."),
			),
			mcp.WithString("email-one",
				mcp.Description("The primary email address of the company."),
			),
			mcp.WithString("email-two",
				mcp.Description("The secondary email address of the company."),
			),
			mcp.WithString("email-three",
				mcp.Description("The tertiary email address of the company."),
			),
			mcp.WithString("website",
				mcp.Description("The website of the company."),
			),
			mcp.WithString("profile",
				mcp.Description("A profile description for the company."),
			),
			mcp.WithNumber("manager-id",
				mcp.Description("The ID of the user who manages the company."),
			),
			mcp.WithNumber("currency-id",
				mcp.Description("The ID of the currency used by the company."),
			),
			mcp.WithNumber("industry-id",
				mcp.Description("The ID of the industry the company belongs to."),
			),
			mcp.WithArray("tag-ids",
				mcp.Description("A list of tag IDs to associate with the company."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var companyUpdate twcompany.Update

			err := twmcp.ParamGroup(request.Params.Arguments,
				twmcp.RequiredNumericParam(&companyUpdate.ID, "company-id"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.Name, "name"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.AddressOne, "address-one"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.AddressTwo, "address-two"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.City, "city"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.State, "state"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.Zip, "zip"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.CountryCode, "country-code"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.Phone, "phone"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.Fax, "fax"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.EmailOne, "email-one"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.EmailTwo, "email-two"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.EmailThree, "email-three"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.Website, "website"),
				twmcp.OptionalPointerParam(&companyUpdate.Company.Profile, "profile"),
				twmcp.OptionalNumericPointerParam(&companyUpdate.Company.ManagerID, "manager-id"),
				twmcp.OptionalNumericPointerParam(&companyUpdate.Company.CurrencyID, "currency-id"),
				twmcp.OptionalNumericPointerParam(&companyUpdate.Company.IndustryID, "industry-id"),
				twmcp.OptionalNumericListParam(&companyUpdate.Company.TagIDs, "tag-ids"),
			)
			if err != nil {
				return nil, fmt.Errorf("invalid parameters: %w", err)
			}

			if err := configResources.TeamworkEngine.Do(ctx, &companyUpdate); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Company updated successfully"), nil
		},
	)
}
