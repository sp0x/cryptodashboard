package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

var (
	Get struct {
		CoinGeckoAPIKey     string   `yaml:"coingecko_api_key"`
		UpdateIntervalInSec int      `yaml:"update_interval_in_sec"`
		TargetCurrency      string   `yaml:"target_currency"`
		CryptoCurrencyIds   []string `yaml:"crypto_currency_ids"`
		Port                int      `yaml:"port"`
		Host                string   `yaml:"host"`
	}
)

var (
	defaultUpdateIntervalInSec = 5
	defaultTargetCurrency      = "usd"
	defaultCryptoCurrencyIds   = []string{"bitcoin", "ethereum", "cardano"}
	defaultPort                = 8080
	defaultHost                = "127.0.0.1"
)

func init() {
	Get.CoinGeckoAPIKey = ""
	Get.UpdateIntervalInSec = defaultUpdateIntervalInSec
	Get.TargetCurrency = defaultTargetCurrency
	Get.CryptoCurrencyIds = defaultCryptoCurrencyIds
	Get.Port = defaultPort
	Get.Host = defaultHost
}

// Load the configuration file.
func Load(filename string) (err error) {
	if filename == "" {
		overrideWithEnvVars()
		return validate()
	}

	var contents []byte
	contents, err = os.ReadFile(filename)
	if nil == err {
		yaml.Unmarshal(contents, &Get)
	}

	overrideWithEnvVars()
	return validate()
}

func overrideWithEnvVars() {
	Get.CoinGeckoAPIKey = getEnvStr("COINGECKO_API_KEY", Get.CoinGeckoAPIKey)
	Get.TargetCurrency = getEnvStr("TARGET_CURRENCY", Get.TargetCurrency)
	Get.UpdateIntervalInSec = getEnvInt("UPDATE_INTERVAL_IN_SEC", Get.UpdateIntervalInSec)
	Get.CryptoCurrencyIds = getEnvStrings("CRYPTO_CURRENCY_IDS", Get.CryptoCurrencyIds)
	Get.Port = getEnvInt("PORT", Get.Port)
	Get.Host = getEnvStr("HOST", Get.Host)
}

func validate() error {
	if Get.CoinGeckoAPIKey == "" {
		return fmt.Errorf("coingecko API key is required")
	}
	return nil
}

func getEnvStr(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v, ok := os.LookupEnv(key); ok {
		number, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return defaultValue
		}
		return int(number)
	}
	return defaultValue
}

func getEnvStrings(key string, defaultValue []string) []string {
	if v, ok := os.LookupEnv(key); ok {
		return strings.Split(v, ",")
	}
	return defaultValue
}
