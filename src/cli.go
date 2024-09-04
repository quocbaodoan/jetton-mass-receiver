package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/spf13/cobra"
)

// Configuration stores the user's settings
type Configuration struct {
	JettonMasterAddress  string `json:"jettonMasterAddress"`
	Commentary           string `json:"commentary"`
	MessageEntryFilename string `json:"messageEntryFilename"`
	ReceiverAddress      string `json:"receiverAddress"`
}

type MessageEntry struct {
	Seed string `json:"seed"`
}

const configFilePath = "config.json"

var rootCmd = &cobra.Command{
	Use:   "toncli",
	Short: "TON CLI application for managing transactions",
	Run:   runMainLogic,
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up or update the application configuration",
	Run:   func(cmd *cobra.Command, args []string) { setupConfiguration() },
}

func runMainLogic(cmd *cobra.Command, args []string) {
	config, err := loadConfiguration(configFilePath)
	if err != nil {
		fmt.Println("Configuration not found or invalid. Please run `setup` command.")
		return
	}
	fmt.Println("Running main logic with current configuration...")

	// Read message entries from file
	data, err := os.ReadFile(config.MessageEntryFilename)
	if err != nil {
		log.Fatal("Error reading file:", err.Error())
		return
	}

	var messages []MessageEntry
	err = json.Unmarshal(data, &messages)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err.Error())
		return
	}

	for i, message := range messages {
		log.Printf("Wallet %d", i+1)
		massSender(message.Seed, config.JettonMasterAddress, config.Commentary, config.ReceiverAddress)
	}

}

func setupConfiguration() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter Receive Address:")
	receiverAddress, _ := reader.ReadString('\n')

	fmt.Println("Enter Token Address:")
	jettonMasterAddress, _ := reader.ReadString('\n')

	fmt.Println("Enter Commentary:")
	commentary, _ := reader.ReadString('\n')

	fmt.Println("Enter Message Entry Filename:")
	messageEntryFilename, _ := reader.ReadString('\n')

	config := Configuration{
		ReceiverAddress:      receiverAddress,
		JettonMasterAddress:  jettonMasterAddress,
		Commentary:           commentary,
		MessageEntryFilename: messageEntryFilename,
	}
	if err := saveConfiguration(config, configFilePath); err != nil {
		log.Fatalf("Failed to save configuration: %v", err)
	}

	fmt.Println("Configuration saved successfully.")
}

func saveConfiguration(config Configuration, filePath string) error {
	// Trim whitespace from the configuration values
	config.ReceiverAddress = strings.TrimSpace(config.ReceiverAddress)
	config.JettonMasterAddress = strings.TrimSpace(config.JettonMasterAddress)
	config.Commentary = strings.TrimSpace(config.Commentary)
	config.MessageEntryFilename = strings.TrimSpace(config.MessageEntryFilename)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func loadConfiguration(filePath string) (Configuration, error) {
	var config Configuration
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
