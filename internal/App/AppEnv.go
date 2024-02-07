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
	LISTEN_HOST string
	LISTEN_PORT int

	/* ... */
}

var AppEnv AppEnvironment

func init() {
	AppEnv = *loadAppEnvironment()
}

func loadAppEnvironment() *AppEnvironment {
	dotenvPosition := *flag.String("env", ".env", "The path of the environment file")
	fmt.Printf("Load environment from %s\n", dotenvPosition)
	godotenv.Load(dotenvPosition)

	getenv := func(key string) (string, error) {
		value := os.Getenv(key)
		if value == "" {
			return value, fmt.Errorf("Required value missing for %v", key)
		}
		return value, nil
	}

	env := AppEnvironment{}

	ts := reflect.TypeOf(env)
	ps := reflect.ValueOf(&env)

	s := ps.Elem()

	envErrors := []error{}

	if s.Kind() == reflect.Struct {

		for i := 0; i < s.NumField(); i++ {
			fname := ts.Field(i).Name
			f := s.FieldByName(fname)
			if f.IsValid() {
				if f.CanSet() {
					envValue, err := getenv(fname)
					if err != nil {
						envErrors = append(envErrors, err)
					} else if f.Kind() == reflect.String {
						f.SetString(envValue)
					} else if f.Kind() == reflect.Int {
						strValue := envValue
						intValue, err := strconv.Atoi(strValue)
						if err != nil {
							envErrors = append(envErrors, fmt.Errorf("env value(int): %s=%v: %v", fname, strValue, err))
						} else {
							f.SetInt(int64(intValue))
						}
					} else {
						envErrors = append(envErrors, fmt.Errorf("Bad env type for app env %v", fname))
					}
				}
			}
		}

	}

	if len(envErrors) > 0 {
		for _, err := range envErrors {
			fmt.Println(err)
		}
		panic("Incomplete env")
	}

	return &env
}

func (appEnv *AppEnvironment) GetListener() string {
	return fmt.Sprintf("%v:%v", appEnv.LISTEN_HOST, appEnv.LISTEN_PORT)
}
