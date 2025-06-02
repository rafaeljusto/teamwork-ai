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
	twcompany "github.com/rafaeljusto/teamwork-ai/internal/twapi/company"
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

func registerResources(mcpServer *server.MCPServer, configResources *config.Resources) {
	mcpServer.AddResource(resourceList,
		func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			var multiple twcompany.Multiple
			if err := configResources.TeamworkEngine.Do(ctx, &multiple); err != nil {
				return nil, err
			}
			var resourceContents []mcp.ResourceContents
			for _, company := range multiple.Response.Companies {
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

			var company twcompany.Single
			company.ID = companyID
			if err := configResources.TeamworkEngine.Do(ctx, &company); err != nil {
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
}
