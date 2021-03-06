package kit

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Shopify/themekit/kittest"
)

func resetConfig() {
	flagConfig = Configuration{}
	environmentConfig = Configuration{}
}

func TestSetFlagConfig(t *testing.T) {
	defer resetConfig()

	config, _ := NewConfiguration()
	assert.Equal(t, defaultConfig, *config)

	flagConfig := Configuration{
		Timeout: 2000000,
	}
	SetFlagConfig(flagConfig)

	config, _ = NewConfiguration()
	assert.Equal(t, flagConfig.Timeout, config.Timeout)
}

func TestConfiguration_Env(t *testing.T) {
	defer resetConfig()

	config, _ := NewConfiguration()
	assert.Equal(t, defaultConfig, *config)

	environmentConfig = Configuration{
		Password:     "password",
		ThemeID:      "themeid",
		Domain:       "nope.myshopify.com",
		Directory:    "../kit",
		IgnoredFiles: []string{"one", "two", "three"},
		Proxy:        ":3000",
		Ignores:      []string{"four", "five", "six"},
		Timeout:      40 * time.Second,
	}

	config, _ = NewConfiguration()
	assert.Equal(t, environmentConfig, *config)
}

func TestConfiguration_Precedence(t *testing.T) {
	defer resetConfig()

	config := &Configuration{Password: "file"}
	config, _ = config.compile(true)
	assert.Equal(t, "file", config.Password)

	environmentConfig = Configuration{Password: "environment"}
	config, _ = config.compile(true)
	assert.Equal(t, "environment", config.Password)

	flagConfig = Configuration{Password: "flag"}
	config, _ = config.compile(true)
	assert.Equal(t, "flag", config.Password)

	config = &Configuration{Password: "file"}
	config, _ = config.compile(false)
	assert.Equal(t, "file", config.Password)
}

func TestConfiguration_Validate(t *testing.T) {
	defer resetConfig()

	config := Configuration{Password: "file", ThemeID: "123", Domain: "test.myshopify.com"}
	assert.Nil(t, config.Validate())

	config = Configuration{Password: "file", ThemeID: "live", Domain: "test.myshopify.com"}
	assert.Nil(t, config.Validate())

	config = Configuration{ThemeID: "123", Domain: "test.myshopify.com"}
	err := config.Validate()
	if assert.NotNil(t, err) {
		assert.True(t, strings.Contains(err.Error(), "missing password"))
	}

	config = Configuration{Password: "test", ThemeID: "123", Domain: "test.nope.com"}
	err = config.Validate()
	if assert.NotNil(t, err) {
		assert.True(t, strings.Contains(err.Error(), "invalid store domain"))
	}

	config = Configuration{Password: "test", ThemeID: "123"}
	err = config.Validate()
	if assert.NotNil(t, err) {
		assert.True(t, strings.Contains(err.Error(), "missing store domain"))
	}

	config = Configuration{Password: "file", Domain: "test.myshopify.com"}
	err = config.Validate()
	if assert.NotNil(t, err) {
		assert.True(t, strings.Contains(err.Error(), "missing theme_id"))
	}

	config = Configuration{Password: "file", ThemeID: "abc", Domain: "test.myshopify.com"}
	err = config.Validate()
	if assert.NotNil(t, err) {
		assert.True(t, strings.Contains(err.Error(), "invalid theme_id"))
	}

	kittest.GenerateProject()
	defer kittest.Cleanup()
	config = Configuration{ThemeID: "123", Password: "abc123", Domain: "test.myshopify.com", Directory: kittest.SymlinkProjectPath}
	assert.Nil(t, config.Validate())
	assert.Equal(t, kittest.FixtureProjectPath, config.Directory)

	config = Configuration{ThemeID: "123", Password: "abc123", Domain: "test.myshopify.com", Directory: kittest.SymlinkProjectPath}
	os.Remove(kittest.SymlinkProjectPath)
	os.Symlink("nope", kittest.SymlinkProjectPath)
	assert.NotNil(t, config.Validate())

	config = Configuration{ThemeID: "123", Password: "abc123", Domain: "test.myshopify.com", Directory: kittest.SymlinkProjectPath}
	os.Remove(kittest.SymlinkProjectPath)
	assert.NotNil(t, config.Validate())
}

func TestConfiguration_IsLive(t *testing.T) {
	defer resetConfig()

	config := Configuration{ThemeID: "123"}
	assert.False(t, config.IsLive())

	config = Configuration{ThemeID: "live"}
	assert.True(t, config.IsLive())
}

func TestConfiguration_AsYaml(t *testing.T) {
	defer resetConfig()

	config := Configuration{Directory: defaultConfig.Directory}
	assert.Equal(t, "", config.asYAML().Directory)

	config = Configuration{Directory: "nope"}
	assert.Equal(t, "nope", config.asYAML().Directory)

	config = Configuration{Timeout: defaultConfig.Timeout}
	assert.Equal(t, time.Duration(0), config.asYAML().Timeout)

	config = Configuration{Timeout: 42}
	assert.Equal(t, time.Duration(42), config.asYAML().Timeout)
}
