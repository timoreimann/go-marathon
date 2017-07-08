/*
Copyright 2016 Rohith All rights reserved.

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

// EnvironmentVariable is an environment variable associated with an application
type EnvironmentVariable struct {
	EnvVar string `json:",omitempty"`
	EnvSecret
}

// EnvSecret carries environment secrets
type EnvSecret struct {
	Secret string `json:"secret,omitempty"`
}

// UnmarshalJSON unmarshals the given JSON into an EnvironmentVariable.  It will
// take both normal environment variables and secrets
func (ev *EnvironmentVariable) UnmarshalJSON(b []byte) error {
	var errIsEnvVar, errIsEnvSecret error
	if errIsEnvVar = json.Unmarshal(b, &ev.EnvVar); errIsEnvVar == nil {
		return nil
	}

	if errIsEnvSecret = json.Unmarshal(b, &ev.EnvSecret); errIsEnvSecret == nil {
		return nil
	}
	return fmt.Errorf("failed to unmarshal environment variable: unmarshaling into environment variable returned error '%s'; unmarshaling into secret returned error '%s'", errIsEnvVar, errIsEnvSecret)
}

// MarshalJSON marshals the environment variable into either a normal environment variable
// or a secret
func (ev EnvironmentVariable) MarshalJSON() ([]byte, error) {
	if (EnvSecret{}) == ev.EnvSecret {
		return json.Marshal(ev.EnvVar)
	}
	return json.Marshal(ev.EnvSecret)
}
