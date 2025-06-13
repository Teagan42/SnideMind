package mcp

import (
	"context"
	"fmt"

	mcpClient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/teagan42/snidemind/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MCPClient struct {
	Client *mcpClient.Client
	Config config.MCPServerConfig
}

type Params struct {
	fx.In
	Log *zap.Logger
	Cfg config.MCPServerConfig
	Lc  fx.Lifecycle
}

type Result struct {
	fx.Out
	Client *MCPClient `name:"client"`
}

func NewMCPClient(p Params) (*Result, error) {
	var clientTransport transport.Interface
	switch p.Cfg.Type {
	case "sse":
		if sseTransport, err := transport.NewSSE(p.Cfg.URL); err != nil {
			return nil, err
		} else {
			clientTransport = sseTransport
		}
	case "http":
		if httpTransport, err := transport.NewStreamableHTTP(p.Cfg.URL); err != nil {
			return nil, err
		} else {
			clientTransport = httpTransport
		}
	}

	return &Result{
		Client: &MCPClient{
			Client: mcpClient.NewClient(
				clientTransport,
			),
			Config: p.Cfg,
		},
	}, nil

}

func (c *MCPClient) Start(context context.Context) {
	fmt.Printf("Starting MCP client for %s at %s\n", c.Config.Name, c.Config.URL)
	if err := c.Client.Start(context); err != nil {
		fmt.Printf("Failed to start MCP client for %s: %v\n", c.Config.Name, err)
		panic(err)
	} else {
		fmt.Printf("MCP client for %s started successfully\n", c.Config.Name)
		c.Client.OnNotification(c.HandleNotification)
		if result, err := c.Client.Initialize(context, mcp.InitializeRequest{}); err != nil {
			fmt.Printf("Failed to initialize MCP client for %s: %v\n", c.Config.Name, err)
			panic(err)
		} else {
			fmt.Printf("MCP client for %s initialized successfully with capabilities: %v\n", c.Config.Name, result.Capabilities)
		}
	}
}

func (c *MCPClient) HandleNotification(notification mcp.JSONRPCNotification) {
	fmt.Printf("Received notification: %s with %+v\n", notification.Method, notification.Params)
}

func (c *MCPClient) Close() error {
	if err := c.Client.Close(); err != nil {
		return err
	}
	return nil
}

func (c *MCPClient) ListTools(context context.Context, request mcp.ListToolsRequest) (*[]mcp.Tool, error) {
	tools, err := c.Client.ListTools(context, request)
	if err != nil {
		return nil, err
	}
	if tools == nil || len(tools.Tools) == 0 {
		// No tools available
		return nil, nil
	}
	var filteredTools []mcp.Tool
	if c.Config.Blacklist != nil {
		// Filter out blacklisted tools
		for _, tool := range tools.Tools {
			if !c.Config.Blacklist.IsPromptBlacklisted(tool.Name) {
				filteredTools = append(filteredTools, tool)
			}
		}
		tools.Tools = filteredTools
	} else {
		// No blacklist, return all tools
		filteredTools = tools.Tools
	}
	return &filteredTools, nil
}

func (c *MCPClient) ListPrompts(context context.Context, request mcp.ListPromptsRequest) (*[]mcp.Prompt, error) {
	prompts, err := c.Client.ListPrompts(context, request)
	if err != nil {
		return nil, err
	}
	if prompts == nil || len(prompts.Prompts) == 0 {
		// No prompts available
		return nil, nil
	}
	var filteredPrompts []mcp.Prompt
	if c.Config.Blacklist != nil {
		// Filter out blacklisted prompts
		for _, prompt := range prompts.Prompts {
			if !c.Config.Blacklist.IsPromptBlacklisted(prompt.Name) {
				filteredPrompts = append(filteredPrompts, prompt)
			}
		}
		prompts.Prompts = filteredPrompts
	} else {
		// No blacklist, return all prompts
		filteredPrompts = prompts.Prompts
	}
	return &filteredPrompts, nil
}

func (c *MCPClient) ListResources(context context.Context, request mcp.ListResourcesRequest) (*[]mcp.Resource, error) {
	resources, err := c.Client.ListResources(context, request)
	if err != nil {
		return nil, err
	}
	if resources == nil || len(resources.Resources) == 0 {
		// No resources available
		return nil, nil
	}
	var filteredResources []mcp.Resource
	if c.Config.Blacklist != nil {
		// Filter out blacklisted resources
		for _, resource := range resources.Resources {
			if !c.Config.Blacklist.IsResourceBlacklisted(resource.Name) {
				filteredResources = append(filteredResources, resource)
			}
		}
		resources.Resources = filteredResources
	} else {
		// No blacklist, return all resources
		filteredResources = resources.Resources
	}
	return &filteredResources, nil
}
