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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentVariableAPI(t *testing.T) {
	app := Application{}
	require.Nil(t, app.Env)
	app.AddEnv("FOO", "bar")
	app.AddSecret("TOP", "secret", "/path/to/secret")
	assert.Equal(t, `bar`, (*app.Env)["FOO"].EnvVar)
	assert.Equal(t, EnvSecret{"secret"}, (*app.Env)["TOP"].EnvSecret)
	assert.Equal(t, Secret{"/path/to/secret"}, (*app.Secrets)["secret"])

	app.EmptyEnvs()
	require.NotNil(t, app.Env)
	assert.Equal(t, "", (*app.Env)["FOO"].EnvVar)

	app.EmptySecrets()
	require.NotNil(t, app.Secrets)
	assert.Equal(t, Secret{}, (*app.Secrets)["secret"])
}

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
	assert.Equal(t, `bar`, (*env)["FOO"].EnvVar)
	assert.Equal(t, EnvSecret{"secret"}, (*env)["TOP"].EnvSecret)
	assert.Equal(t, Secret{"/path/to/secret"}, (*secrets)["secret"])
}

func TestEnvironmentVariableMarshalIllegal(t *testing.T) {
	j := []byte(`{false}`)
	envvar := EnvironmentVariable{}
	assert.Error(t, envvar.UnmarshalJSON(j))
}

func TestEnvironmentVariableMarshal(t *testing.T) {
	tests := []struct {
		name     string
		env      EnvironmentVariable
		wantJSON string
	}{
		{
			name:     "env",
			env:      EnvironmentVariable{"bar", EnvSecret{}},
			wantJSON: `"bar"`,
		},
		{
			name:     "secret",
			env:      EnvironmentVariable{"", EnvSecret{"secret"}},
			wantJSON: `{"secret":"secret"}`,
		},
	}

	for _, test := range tests {
		label := fmt.Sprintf("test: %s", test.name)
		j, err := test.env.MarshalJSON()
		if assert.NoError(t, err, label) {
			assert.Equal(t, test.wantJSON, string(j), label)
		}
	}
}
