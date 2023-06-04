package constant

import (
	"log"
	"regexp"

	"github.com/spf13/viper"

	"BackendServer/config"
)

var ConnectServer = []string{
	config.UserServerUrl,
}

type Message struct {
	EventType int         `json:"event_type"`
	Data      interface{} `json:"data"`
}

func ReadConfig(configPath string) {
	viper.SetConfigFile(configPath)
	viper.AddConfigPath(".")

	viper.SetDefault("RUN_MODE", "debug")

	envs := []string{
		"PORT",
		"RUN_MODE",
	}

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			log.Println(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	config.Initialize()
}

func ValidateStringChinese(input string) bool {
	for _, b := range input {
		if b >= 0x4e00 && b <= 0x9fff {
			return false
		}
	}
	return true
}

func ValiadateStringLength(input string, max int, min int) bool {
	return len(input) <= max && len(input) >= min
}

func ValidateStringSpecialCode(input string) bool {
	regex := "[!@#$%^&*()_+\\-=[\\]{}|\\\\;:'\",.<>/?]"
	re := regexp.MustCompile(regex)
	return !re.MatchString(input)
}

func ValidateString(input string, chinese, length, specialCode bool, max, min int) bool {
	check := true
	if chinese {
		check = ValidateStringChinese(input)
		if !check {
			return check
		}
	}

	if length {
		check = ValiadateStringLength(input, max, min)
		if !check {
			return check
		}
	}

	if specialCode {
		check = ValidateStringSpecialCode(input)
		if !check {
			return check
		}
	}
	return check
}
