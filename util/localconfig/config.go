package localconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Test   string `mapstructure:"name"`
	Xendit Xendit `mapstructure:"xendit"`
}

type Xendit struct {
	EWallet Ewallet `mapstructure:"ewallet"`
}

type Ewallet struct {
	LegacyEnabled bool          `mapstructure:"enabled"`
	OVO           EwalletConfig `mapstructure:"ovo"`
	Dana          EwalletConfig `mapstructure:"dana"`
	LinkAja       EwalletConfig `mapstructure:"linkaja"`
}

type EwalletConfig struct {
	UseInvoice bool `mapstructure:"use_invoice"`
	UseLegacy  bool `mapstructure:"use_legacy"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadConfigFromBytes(data)
}

func LoadConfigFromBytes(data []byte) (*Config, error) {
	fang := viper.New()
	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()
	fang.SetEnvPrefix("GOPAYMENT")
	fang.SetConfigType("yaml")

	if err := fang.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	x := fang.GetString("name")
	fmt.Println(x)

	var cfg Config
	err := fang.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Error loading creds: %v", err)
	}

	return &cfg, nil
}
