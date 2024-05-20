package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestLoadByPath_HappyPath(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		fileName, err := createTempConfigFile(t)
		require.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(fileName)
		})

		cfg, err := loadByPath(fileName)

		require.NoError(t, err)
		assertConfig(t, cfg)
	})
}

func TestLoadFromEnv_HappyPath(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		setEnvVariables(t)

		cfg := mustLoadFromEnv()

		assertConfig(t, cfg)
	})
}

func TestMustLoad_FromYaml(t *testing.T) {
	t.Run("from yaml", func(t *testing.T) {
		fileName, err := createTempConfigFile(t)
		if err != nil {
			require.NoError(t, err)
		}
		t.Cleanup(func() {
			os.Remove(fileName)
		})

		os.Args = []string{"cmd", "-config", fileName}
		cfg := MustLoad()

		assertConfig(t, cfg)
	})
}

func TestMustLoad_FromEnv(t *testing.T) {
	t.Run("from env", func(t *testing.T) {
		setEnvVariables(t)

		cfg := MustLoad()

		assertConfig(t, cfg)
	})
}

func TestLoadByPath_FailPath_WrongPath(t *testing.T) {
	t.Run("wrong path", func(t *testing.T) {
		_, err := loadByPath("./some/wrong/path")

		require.Error(t, err, "have to return error")
	})
}

func createTempConfigFile(t *testing.T) (string, error) {
	t.Helper()
	f, err := os.Create("./temp_config.yaml")
	if err != nil {
		return "", err
	}
	defer f.Close()

	configData := []byte(`env: "test"
grpc:
  port: 44046
  timeout: "1h"
rabbitmq:
  user: "admin"
  password: "admin"
  host: "localhost"
  port: "5672"
mongodb:
  uri: "mongodb://admin:admin@localhost:27017"
  ping_timeout: "10s"
  database_name: "ucms-posts-dev"`)

	_, err = f.Write(configData)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func createTempConfigFileWithWrongData(t *testing.T) (string, error) {
	t.Helper()
	f, err := os.Create("./temp_config.yaml")
	if err != nil {
		return "", err
	}
	defer f.Close()

	configData := []byte(`env: "test"
grpc:
  port: 44046
  timeout: "1h"
rabbittttttmq:
  user: "admin"
  password: "admin"
mongodb:
uri: "mongodb://admin:aarstdmin@localhost:27017"
ping_timeout: "10s"
database_name: "ucms-posts-dev"`)

	_, err = f.Write(configData)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func setEnvVariables(t *testing.T) {
	t.Helper()
	t.Setenv("ENV", "test")

	t.Setenv("GRPC_PORT", "44046")
	t.Setenv("GRPC_TIMEOUT", "1h")

	t.Setenv("MONGODB_URI", "mongodb://admin:admin@localhost:27017")
	t.Setenv("MONGODB_PING_TIMEOUT", "10s")
	t.Setenv("MONGODB_DATABASE_NAME", "ucms-posts-dev")

	t.Setenv("RABBITMQ_USER", "admin")
	t.Setenv("RABBITMQ_PASSWORD", "admin")
	t.Setenv("RABBITMQ_HOST", "localhost")
	t.Setenv("RABBITMQ_PORT", "5672")
}

func assertConfig(t *testing.T, cfg *Config) {
	t.Helper()

	assert.Equal(t, "test", cfg.Env)

	assert.Equal(t, 44046, cfg.GRPC.Port)
	assert.Equal(t, time.Hour, cfg.GRPC.Timeout)

	assert.Equal(t, "mongodb://admin:admin@localhost:27017", cfg.MongoDB.URI)
	assert.Equal(t, 10*time.Second, cfg.MongoDB.PingTimeout)
	assert.Equal(t, "ucms-posts-dev", cfg.MongoDB.DatabaseName)

	assert.Equal(t, "admin", cfg.Rabbitmq.User)
	assert.Equal(t, "admin", cfg.Rabbitmq.Password)
	assert.Equal(t, "localhost", cfg.Rabbitmq.Host)
	assert.Equal(t, "5672", cfg.Rabbitmq.Port)
}
