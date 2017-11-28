package cfsvcenv

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const vcapServicesEnv = "VCAP_SERVICES"

// Bind will take each of the Cloud Foundry service credentials found in
// the VCAP_SERVICES environment variable and put each services's credential
// value in the environment using setenv.
//
// In order to avoid conflicts across services that have a credential with the
// same name (e.g. `service-1` and `service-2` both have a credential named
// `API_KEY`), each environment variable is prefixed with a transformed service
// name. Credential name will also be transformed - this is because it is
// conventional for environment variables to be upper case. E.g. for the
// `service-1` credential `api-key`, the resulting environment variable will be
// named `SERVICE_1_API_KEY`.
//
// It is assumed that the service `credentials` is a map of key/value pairs. If
// it is not, that service will simply be ignored.
// Each credential value is converted to a string using the %v verb from fmt. In
// this way, it does not make sense to use anything except primitives for
// credential values.
func Bind() error {
	vcapServicesStr := os.Getenv(vcapServicesEnv)
	if vcapServicesStr == "" {
		return nil
	}
	var vcapServices map[string][]struct {
		Name        string      `json:"name"`
		Credentials interface{} `json:"credentials"`
	}
	if err := json.NewDecoder(strings.NewReader(vcapServicesStr)).Decode(&vcapServices); err != nil {
		return err
	}
	for _, services := range vcapServices {
		for _, service := range services {
			prefix := serviceEnvPrefix(service.Name)
			creds, ok := service.Credentials.(map[string]interface{})
			if !ok {
				continue
			}
			for k, v := range creds {
				name := fmt.Sprintf("%s%s", prefix, toEnvName(k))
				if _, ok := os.LookupEnv(name); ok {
					continue
				}
				val := fmt.Sprintf("%v", v)
				os.Setenv(name, val)
			}
		}
	}
	return nil
}

func toEnvName(s string) string {
	return strings.Replace(strings.ToUpper(s), "-", "_", -1)
}

func serviceEnvPrefix(name string) string {
	return fmt.Sprintf("%s_", toEnvName(name))
}
