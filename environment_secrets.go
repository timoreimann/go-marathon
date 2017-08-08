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
		Env     map[string]interface{} `json:"env,omitempty"`
		Secrets map[string]TmpSecret   `json:"secrets,omitempty"`
	}{
		Alias: (*Alias)(app),
	}
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}
	env := make(map[string]string)
	secrets := make(map[string]Secret)

	for k, v := range aux.Env {
		if s, ok := v.(string); ok {
			env[k] = s
			continue
		}
		tmp, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("unrecognized environment variable type: %s", v)
		}
		s := new(TmpEnvSecret)
		if err := json.Unmarshal(tmp, &s); err != nil {
			return fmt.Errorf("unrecognized environment variable type: %s", v)
		}
		secrets[s.Secret] = Secret{EnvVar: k}
	}
	app.Env = env
	for k, v := range aux.Secrets {
		tmp := secrets[k]
		tmp.Source = v.Source
		secrets[k] = tmp
	}
	app.Secrets = secrets
	return nil
}

// MarshalJSON marshals secrets into its environment variable and secret pieces, and all other environment
// variables into into env
func (app *Application) MarshalJSON() ([]byte, error) {
	env := make(map[string]interface{})
	secrets := make(map[string]TmpSecret)

	for k, v := range app.Env {
		env[string(k)] = string(v)
	}

	for k, v := range app.Secrets {
		env[v.EnvVar] = TmpEnvSecret{Secret: k}
		secrets[k] = TmpSecret{v.Source}
	}
	aux := &struct {
		*Alias
		Env     map[string]interface{} `json:"env,omitempty"`
		Secrets map[string]TmpSecret   `json:"secrets,omitempty"`
	}{Alias: (*Alias)(app), Env: env, Secrets: secrets}

	return json.Marshal(aux)
}
