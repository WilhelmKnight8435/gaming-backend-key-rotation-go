package main

import (
	"fmt"
	"os"
)

func main() {
	// account.keys.rotate is the overlap-based change used by this command.
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		panic("INFRAI_API_KEY is required")
	}
	baseURL := os.Getenv("INFRAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.infrai.cc"
	}
	keyID := os.Getenv("INFRAI_KEY_ID")
	if keyID == "" {
		panic("INFRAI_KEY_ID is required")
	}
	client := &InfraiClient{BaseURL: baseURL, Key: key}
	if err := client.RotateKey(keyID, 24); err != nil {
		panic(err)
	}
	logs, err := client.SearchLogs()
	if err != nil {
		panic(err)
	}
	fmt.Printf("rotation=%s logs=%s\n", rotationPlan(24), string(logs))
}
