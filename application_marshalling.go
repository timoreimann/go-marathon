/*
Copyright 2017 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marathon

import (
	"encoding/json"
	"fmt"
)

// Alias aliases the Application struct so that it will be marshaled/unmarshaled automatically
type Alias Application

// TmpEnvSecret holds the secret values deserialized from the environment variables field
type TmpEnvSecret struct {
	Secret string `json:"secret,omitempty"`
}

// TmpSecret holds the deserialized secrets field in a Marathon application configuration
type TmpSecret struct {
	Source string `json:"source,omitempty"`
}

// UnmarshalJSON unmarshals the given Application JSON into an environment variables and secrets.
func (app *Application) UnmarshalJSON(b []byte) error {
	aux := &struct {
		*Alias
		Env     map[string]interface{} `json:"env"`
		Secrets map[string]TmpSecret   `json:"secrets"`
	}{
		Alias: (*Alias)(app),
	}
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}
	env := &map[string]string{}
	secrets := &map[string]Secret{}

	for envName, genericEnvValue := range aux.Env {
		switch envValOrSecret := genericEnvValue.(type) {
		case string:
			(*env)[envName] = envValOrSecret
		case map[string]interface{}:
			for _, secretStore := range envValOrSecret {
				if secStore, ok := secretStore.(string); ok {
					(*secrets)[secStore] = Secret{EnvVar: envName}
					break
				}
				return fmt.Errorf("unexpected secret value type %T", envValOrSecret[envName])
			}
		default:
			return fmt.Errorf("unexpected environment variable type %T", envValOrSecret)
		}
	}
	app.Env = env
	for k, v := range aux.Secrets {
		tmp := (*secrets)[k]
		tmp.Source = v.Source
		(*secrets)[k] = tmp
	}
	app.Secrets = secrets
	return nil
}

// MarshalJSON marshals secrets into its environment variable and secret pieces, and all other environment
// variables into into env
func (app *Application) MarshalJSON() ([]byte, error) {
	env := make(map[string]interface{})
	secrets := make(map[string]TmpSecret)

	if app.Env != nil {
		for k, v := range *app.Env {
			env[string(k)] = string(v)
		}
	}
	if app.Secrets != nil {
		for k, v := range *app.Secrets {
			env[v.EnvVar] = TmpEnvSecret{Secret: k}
			secrets[k] = TmpSecret{v.Source}
		}
	}
	aux := &struct {
		*Alias
		Env     map[string]interface{} `json:"env,omitempty"`
		Secrets map[string]TmpSecret   `json:"secrets,omitempty"`
	}{Alias: (*Alias)(app), Env: env, Secrets: secrets}

	return json.Marshal(aux)
}
