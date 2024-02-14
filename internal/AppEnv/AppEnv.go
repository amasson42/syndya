package AppEnv

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppEnvironment struct {
	LISTEN_HOST string `default:"localhost"`
	LISTEN_PORT int    `default:"8080"`

	METADATAS_LIST         string `default:"_"`
	METADATAS_REVIVEPERIOD int    `default:"3000"`

	MATCHFINDER_LUASCRIPT    string `default:"_"`
	MATCHFINDER_TIMEINTERVAL int    `default:"5000"`
	MATCHFINDER_RESETSTATE   bool   `default:"false"`
}

var AppEnv AppEnvironment

func init() {
	AppEnv = *loadAppEnvironment()
}

func loadAppEnvironment() *AppEnvironment {
	envFilePath := *flag.String("env", ".env", "Path to the environment file")
	fmt.Printf("Loading environment from %s\n", envFilePath)
	godotenv.Load(envFilePath)

	var errors []error

	appEnv := AppEnvironment{}

	structType := reflect.TypeOf(appEnv)
	structValue := reflect.ValueOf(&appEnv).Elem()

	for i := 0; i < structValue.NumField(); i++ {

		field := structValue.Field(i)

		fieldName := structType.Field(i).Name
		defaultValue := structType.Field(i).Tag.Get("default")

		envValue := os.Getenv(fieldName)
		if envValue == "" {
			if defaultValue == "" && field.Kind() != reflect.Ptr {
				errors = append(errors, fmt.Errorf("required value missing for %v", fieldName))
			}
			envValue = defaultValue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				errors = append(errors, fmt.Errorf("env value(int): %s=%v: %v", fieldName, envValue, err))
			} else {
				field.SetInt(int64(intValue))
			}
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				errors = append(errors, fmt.Errorf("env value(bool): %s=%v: %v", fieldName, envValue, err))
			} else {
				field.SetBool(boolValue)
			}
		default:
			errors = append(errors, fmt.Errorf("bad env type for app env %v", fieldName))
		}
	}

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		panic("Invalid env")
	}

	return &appEnv
}

func (appEnv *AppEnvironment) GetListener() string {
	return fmt.Sprintf("%v:%v", appEnv.LISTEN_HOST, appEnv.LISTEN_PORT)
}

func (appEnv *AppEnvironment) GetMetadataList() []string {
	if appEnv.METADATAS_LIST == "_" {
		return []string{}
	}
	return strings.Split(appEnv.METADATAS_LIST, ",")
}

func (appEnv *AppEnvironment) HasMatchFinderScript() bool {
	return appEnv.MATCHFINDER_LUASCRIPT != "_"
}
