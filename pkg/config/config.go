package config

import (
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"

	"slack-bot/pkg/utils"
)

var (
	SLACK_BOT_TOKEN string
	SLACK_APP_TOKEN string
)

var (
	FileCache = cache.New(30*time.Minute, 10*time.Minute)
)

var (
	AZURE_STORAGE_ACCOUNT     string
	AZURE_STORAGE_ACCESS_KEY  string
	AZURE_STORAGE_ACCOUNT_URL string
	CONTAINER                 = "slack-bot"

	AzClient          *azblob.ServiceClient
	AzContainerClient *azblob.ContainerClient
)

func init() {
	loadEnv()
	setEnv()
	setUpAzure()
}

func setUpAzure() {
	credential, err := azblob.NewSharedKeyCredential(AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_ACCESS_KEY)
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to create Azure credentials - %s\n", err.Error())
	}

	AzClient, err = azblob.NewServiceClientWithSharedKey(AZURE_STORAGE_ACCOUNT_URL, credential, nil)
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to create Azure client - %s\n", err.Error())
	}

	AzContainerClient, err = AzClient.NewContainerClient(CONTAINER)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create container client - %s\n", err.Error())
	}

	utils.InfoLogger.Println("Azure Clients are initialized")
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Fatalln("Failed to load .env file")
	}
}

func setEnv() {
	SLACK_APP_TOKEN = os.Getenv("SLACK_APP_TOKEN")
	SLACK_BOT_TOKEN = os.Getenv("SLACK_BOT_TOKEN")

	if SLACK_APP_TOKEN == "" || SLACK_BOT_TOKEN == "" {
		utils.ErrorLogger.Fatalln("SLACK_BOT_TOKEN or SLACK_APP_TOKEN are not set in .env file")
	}

	AZURE_STORAGE_ACCOUNT = os.Getenv("AZURE_STORAGE_ACCOUNT")
	AZURE_STORAGE_ACCESS_KEY = os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	AZURE_STORAGE_ACCOUNT_URL = os.Getenv("AZURE_STORAGE_ACCOUNT_URL")

	if AZURE_STORAGE_ACCOUNT == "" || AZURE_STORAGE_ACCESS_KEY == "" || AZURE_STORAGE_ACCOUNT_URL == "" {
		utils.ErrorLogger.Fatalln("AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY or AZURE_STORAGE_ACCOUNT_URL are not set in .env file")
	}

	utils.InfoLogger.Println("Environment and Config variables are set")
}
