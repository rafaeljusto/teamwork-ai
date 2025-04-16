package company

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rafaeljusto/teamwork-ai/internal/config"
	twcompany "github.com/rafaeljusto/teamwork-ai/internal/teamwork/company"
)

var resourceList = mcp.NewResource("twapi://companies", "companies",
	mcp.WithResourceDescription("Companies, also know as clients, are organizations that the "+
		"customer offers services to."),
	mcp.WithMIMEType("application/json"),
)

var resourceItem = mcp.NewResourceTemplate("twapi://companies/{id}", "company",
	mcp.WithTemplateDescription("Company, also know as client, is an organization that the "+
		"customer offers services to."),
	mcp.WithTemplateMIMEType("application/json"),
)

func Register(mcpServer *server.MCPServer, resources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var companies twcompany.MultipleCompanies
			if err := resources.TeamworkEngine.Do(&companies); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, company := range companies {
				encoded, err := json.Marshal(company)
				if err != nil {
					return nil, err
				}
				resourceContents = append(resourceContents, mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://companies/%d", company.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				})
			}
			return resourceContents, nil
		},
	)

	reCompanyID := regexp.MustCompile(`twapi://companies/(\d+)`)
	mcpServer.AddResourceTemplate(resourceItem,
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			matches := reCompanyID.FindStringSubmatch(request.Params.URI)
			if len(matches) != 2 {
				return nil, fmt.Errorf("invalid company ID")
			}
			companyID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid company ID")
			}

			var company twcompany.SingleCompany
			company.ID = companyID
			if err := resources.TeamworkEngine.Do(&company); err != nil {
				return nil, err
			}

			encoded, err := json.Marshal(company)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      fmt.Sprintf("twapi://companies/%d", company.ID),
					MIMEType: "application/json",
					Text:     string(encoded),
				},
			}, nil
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("retrieve-companies",
			mcp.WithDescription("Retrieve multiple companies, also know as clients, in a customer site of Teamwork.com. "+
				"Companies, also know as clients, are organizations that the customer offers services to."),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var companies twcompany.MultipleCompanies
			if err := resources.TeamworkEngine.Do(&companies); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(companies)
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
			mcp.WithNumber("companyId",
				mcp.Required(),
				mcp.Description("The ID of the company."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var company twcompany.SingleCompany

			id, ok := request.Params.Arguments["companyId"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid companyId")
			} else if id == 0 {
				return nil, fmt.Errorf("companyId is required")
			}
			company.ID = int64(id)

			if err := resources.TeamworkEngine.Do(&company); err != nil {
				return nil, err
			}
			encoded, err := json.Marshal(company)
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
			mcp.WithArray("userIds",
				mcp.Description("List of user IDs assigned to the company. This is a JSON array of integers."),
				mcp.Items(map[string]any{
					"type": "number",
				}),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var company twcompany.CompanyCreation
			var ok bool

			company.Name, ok = request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid name")
			} else if company.Name == "" {
				return nil, fmt.Errorf("name is required")
			}

			if err := resources.TeamworkEngine.Do(&company); err != nil {
				return nil, err
			}
			return mcp.NewToolResultText("Company created successfully"), nil
		},
	)
}
