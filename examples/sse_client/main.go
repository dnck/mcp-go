package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// Helper function to print tool results
func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}

func ping(messageURL string, pingID float64) {
	if messageURL == "" {
		fmt.Printf("Did not receive message endpoint URL\n")
		return
	}

	pingResponse := map[string]any{
		"jsonrpc": "2.0",
		"id":      pingID,
		"result":  map[string]any{},
	}

	requestBody, err := json.Marshal(pingResponse)
	if err != nil {
		fmt.Printf("Failed to marshal ping response: %v\n", err)
		return
	}

	resp, err := http.Post(
		messageURL,
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		fmt.Printf("Failed to send ping response: %v\b", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		fmt.Printf("Expected status 202 for ping response, got %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	if len(body) > 0 {
		var response map[string]any
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Printf("Failed to parse response body: %v\n", err)
			return
		}

		if response["error"] != nil {
			fmt.Printf("Expected no error in response, got %v\n", response["error"])
			return
		}
	}
	fmt.Println("Ping response sent successfully")
	listTools()
}

func listTools() {
	c, _ := client.NewSSEMCPClient("http://localhost:8080/sse", client.WithHeaders(map[string]string{}))
	if e := c.Start(context.Background()); e != nil {
		fmt.Printf("Failed to start client: %v\n", e)
		return
	}
	ir := mcp.InitializeRequest{}
	ir.Request = mcp.Request{
		Method: "initialize",
	}
	ir.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	ir.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}
	_, err := c.Initialize(context.Background(), ir)
	if err != nil {
		fmt.Println("failed to init", err.Error())
		return
	}
	fmt.Println("Listing available tools...")
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := c.ListTools(context.Background(), toolsRequest)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
		b, _ := tool.MarshalJSON()
		fmt.Printf("```\n%s\n```\n", string(b))
	}

	call := mcp.CallToolRequest{}
	call.Params.Arguments = map[string]any{
		"a": 1,
		"b": 2,
	}
	//call.Method = "add"
	call.Params.Name = "add"
	result, err := c.CallTool(context.Background(), call)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}
	printToolResult(result)
}

func main() {
	listTools()
	//sseResp, err := http.Get(fmt.Sprintf("%s/sse", "http://127.0.0.1:8080"))
	//if err != nil {
	//	fmt.Printf("Failed to connect to SSE endpoint: %v\n", err)
	//}
	//defer sseResp.Body.Close()
	//
	//reader := bufio.NewReader(sseResp.Body)
	//
	//var messageURL string
	//var pingID float64
	//
	//for {
	//	line, err := reader.ReadString('\n')
	//	if err != nil {
	//		fmt.Printf("Failed to read SSE event: %v\n", err)
	//		return
	//	}
	//	fmt.Printf("Received line: %s\n", line)
	//
	//	if strings.HasPrefix(line, "event: endpoint") {
	//		dataLine, err := reader.ReadString('\n')
	//		if err != nil {
	//			fmt.Printf("Failed to read endpoint data: %v\n", err)
	//			return
	//		}
	//		messageURL = strings.TrimSpace(strings.TrimPrefix(dataLine, "data: "))
	//
	//		fmt.Printf("messageURL: %s\n", messageURL)
	//
	//		_, err = reader.ReadString('\n')
	//		if err != nil {
	//			fmt.Printf("Failed to read blank line: %v\n", err)
	//			return
	//		}
	//
	//		go ping(messageURL, pingID)
	//
	//	}
	//
	//	if strings.HasPrefix(line, "event: message") {
	//		dataLine, err := reader.ReadString('\n')
	//		if err != nil {
	//			fmt.Printf("Failed to read message data: %v\n", err)
	//			return
	//		}
	//
	//		pingData := strings.TrimSpace(strings.TrimPrefix(dataLine, "data:"))
	//		var pingMsg mcp.JSONRPCRequest
	//		if err := json.Unmarshal([]byte(pingData), &pingMsg); err != nil {
	//			fmt.Printf("Failed to parse ping message: %v\n", err)
	//			return
	//		}
	//
	//		if pingMsg.Method == "ping" {
	//			pingID = pingMsg.ID.(float64)
	//			fmt.Printf("Received ping with ID: %f\n", pingID)
	//			break // We got the ping, exit the loop
	//		}
	//
	//		_, err = reader.ReadString('\n')
	//		if err != nil {
	//			fmt.Printf("Failed to read blank line: %v\n", err)
	//			return
	//		}
	//	}
	//
	//	if messageURL != "" && pingID != 0 {
	//		break
	//	}
	//}
}
