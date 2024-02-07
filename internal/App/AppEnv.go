package App

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

type AppEnvironment struct {
	LISTEN_HOST string `default:"localhost"`
	LISTEN_PORT int    `default:"8080"`

	// Add more environment variables here with default values
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
