package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/teagan42/snidemind/internal/config"
	"github.com/teagan42/snidemind/internal/llm"
	"github.com/teagan42/snidemind/internal/mcp"
	"github.com/teagan42/snidemind/internal/server"
	"github.com/teagan42/snidemind/internal/types"
)

var ctx = context.Background()

func main() {
	parser := argparse.NewParser("snidemind", "Snidemind AI Tool CLI Daemon")
	bind := parser.String("b", "bind", &argparse.Options{
		Required: false,
		Help:     "Address to bind the server to.",
	})
	port := parser.Int("p", "port", &argparse.Options{
		Required: false,
		Help:     "Port to run the server on.",
	})
	configPath := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     "Path to the configuration file.",
		Default:  "config.yaml",
	})
	if err := parser.Parse(os.Args); err != nil {
		log.Fatalf("[Argparse] Failed to parse arguments: %v", err)
		panic(err)
	}

	// redisClient = redis.NewClient(&redis.Options{Addr: *redisAddr})
	// if err := redisClient.Ping(ctx).Err(); err != nil {
	// 	log.Fatalf("[Redis] Connection failed: %v", err)
	// }

	var bindAddress *types.Host
	var portNumber *types.Port
	if bind != nil && *bind != "" {
		if bindHost, err := types.NewHost(*bind); err != nil {
			log.Fatalf("[Bind] Invalid bind address: %v", err)
			panic(err)
		} else {
			bindAddress = &bindHost
		}
	}
	if port != nil && *port > 0 {
		if portNum, err := types.NewPort(*port); err != nil {
			log.Fatalf("[Port] Invalid port number: %v", err)
			panic(err)
		} else {
			portNumber = &portNum
		}
	}
	config, err := config.LoadConfig(*configPath, bindAddress, portNumber)
	if err != nil {
		log.Fatalf("[Config] Failed to load configuration: %v", err)
		panic(err)
	} else {
		log.Printf("[Config] Loaded configuration from %s", *configPath)
	}

	fmt.Printf("[Config] %+v\n", config)

	for _, mcpServerConfig := range config.MCPServers {
		if client, err := mcp.NewClient(mcpServerConfig); err != nil {
			log.Fatalf("[MCP] Failed to create client for server %s: %v", mcpServerConfig.Name, err)
			panic(err)
		} else {
			log.Printf("[MCP] Created client for server %s at %s", mcpServerConfig.Name, mcpServerConfig.URL)
			client.Start(ctx)
			fmt.Printf("[MCP] Client for server %s started successfully\n", mcpServerConfig.Name)
		}
	}

	server := server.NewServer(config.Server)
	server.Start(llm.NewLLM(config.LLM))
}
