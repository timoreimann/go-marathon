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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentVariableUnmarshal(t *testing.T) {
	defaultConfig := NewDefaultConfig()
	configs := &configContainer{
		client: &defaultConfig,
		server: &serverConfig{
			scope: "environment-variables",
		},
	}

	endpoint := newFakeMarathonEndpoint(t, configs)
	defer endpoint.Close()

	application, err := endpoint.Client.Application(fakeAppName)
	require.NoError(t, err)

	env := application.Env
	secrets := application.Secrets

	require.NotNil(t, env)
	assert.Equal(t, "bar", (*env)["FOO"])
	assert.Equal(t, "TOP", (*secrets)["secret"].EnvVar)
	assert.Equal(t, "/path/to/secret", (*secrets)["secret"].Source)
}

func TestEnvironmentVaribleMarshal(t *testing.T) {
	testApp := new(Application)
	testApp.AddSecret("TOP", "secret1", "/path/to/secret")
	testApp.AddEnv("FOO", "bar")

	tmp, err := json.MarshalIndent(testApp, "", "  ")
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(testApp.String()), string(tmp))
}
