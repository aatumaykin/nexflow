package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atumaikin/nexflow/internal/shared/config"
)

func main() {
	// Get directory of this file
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	// Build path to config_test_zai.yml
	configPath := filepath.Join(wd, "config_test_zai.yml")

	// Load and validate config_test_zai.yml
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load config_test_zai.yml: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ config_test_zai.yml loaded successfully")
	fmt.Printf("   Default LLM provider: %s\n", cfg.LLM.DefaultProvider)
	fmt.Printf("   Available providers: ")
	for name := range cfg.LLM.Providers {
		fmt.Printf("%s ", name)
	}
	fmt.Println()

	// Check if zai provider exists and is default
	if zaiConfig, ok := cfg.LLM.Providers["zai"]; ok {
		fmt.Println("✅ zai provider configured")
		fmt.Printf("   API Key format: %s\n", zaiConfig.APIKey)
		fmt.Printf("   Base URL: %s\n", zaiConfig.BaseURL)
		fmt.Printf("   Model: %s\n", zaiConfig.Model)
		fmt.Printf("   Temperature: %f\n", zaiConfig.Temperature)
		fmt.Printf("   Max Tokens: %d\n", zaiConfig.MaxTokens)

		if cfg.LLM.DefaultProvider == "zai" {
			fmt.Println("✅ zai is set as default provider")
		} else {
			fmt.Println("⚠️  zai is NOT the default provider")
		}
	} else {
		fmt.Println("❌ zai provider NOT found in config")
		os.Exit(1)
	}

	fmt.Println("\n✅ All checks passed!")
}
