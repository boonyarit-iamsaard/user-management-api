# Step-by-Step Guide: Configuration Management

## Overview

This guide will walk you through implementing production-ready configuration management for the User Management API using the Viper library. The configuration system will support secure loading from environment variables in production and .env files for development, following industry best practices.

## Prerequisites

- Go 1.25.0 or higher
- Basic understanding of Go modules
- Existing project structure from Task 1.2
- Understanding of production vs development environment patterns

## Step 1: Install Required Dependencies

### 1.1 Add Viper and dotenv libraries

```bash
cd /home/boonyarit-iamsaard/workspace/personal/user-management-api
go get github.com/spf13/viper
go get github.com/joho/godotenv
```

### 1.2 Verify installation

Check your `go.mod` file to ensure both libraries are added:

```bash
go mod tidy
```

The `godotenv` library will be used for local development (.env files) while `viper` handles the core configuration management.

## Step 2: Create Configuration Structure

### 2.1 Create config package

Create a new directory for configuration management:

```bash
mkdir -p internal/config
```

### 2.2 Define configuration structure

Create `internal/config/config.go`:

```go
package config

import (
    "fmt"
    "log"

    "github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
    Host         string `mapstructure:"host"`
    Port         string `mapstructure:"port"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
    IdleTimeout  int    `mapstructure:"idle_timeout"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
    Secret         string `mapstructure:"secret"`
    ExpirationTime int    `mapstructure:"expiration_time"`
}

// LoggerConfig holds logger-related configuration
type LoggerConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
    config := &Config{}

    // Set config file path and name
    viper.SetConfigFile(configPath)

    // Enable environment variable support
    viper.AutomaticEnv()

    // Set default values
    setDefaults()

    // Map environment variables to config struct (using standard Go patterns)
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // Try to read config file
    if err := viper.ReadInConfig(); err != nil {
        // If config file doesn't exist, log and continue with defaults and env vars
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Printf("Config file not found: %s. Using defaults and environment variables", configPath)
        } else {
            return nil, fmt.Errorf("error reading config file: %w", err)
        }
    }

    // Unmarshal config into struct
    if err := viper.Unmarshal(config); err != nil {
        return nil, fmt.Errorf("error unmarshaling config: %w", err)
    }

    return config, nil
}

// setDefaults defines default configuration values
func setDefaults() {
    // Server defaults
    viper.SetDefault("server.host", "localhost")
    viper.SetDefault("server.port", "3000")
    viper.SetDefault("server.read_timeout", 30)
    viper.SetDefault("server.write_timeout", 30)
    viper.SetDefault("server.idle_timeout", 60)

    // Database defaults
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", 5432)
    viper.SetDefault("database.user", "postgres")
    viper.SetDefault("database.password", "")
    viper.SetDefault("database.dbname", "user_management")
    viper.SetDefault("database.sslmode", "disable")

    // JWT defaults
    viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
    viper.SetDefault("jwt.expiration_time", 24) // hours

    // Logger defaults
    viper.SetDefault("logger.level", "info")
    viper.SetDefault("logger.format", "json")
}

// How JWT Secret Override Works:
//
// The default JWT secret "your-secret-key-change-in-production" will be overridden
// in this order (highest to lowest priority):
//
// 1. Environment Variables (ALWAYS overrides):
//    export JWT_SECRET=super-secure-production-key
//    OR in .env: JWT_SECRET=my-dev-secret
//
// 2. Config File (overrides default only):
//    jwt:
//      secret: "my-config-file-secret"
//
// 3. Code Default (used only if nothing else is set):
//    "your-secret-key-change-in-production" (SECURITY RISK!)
//
// Examples:
// - Development: .env file sets JWT_SECRET=my-dev-secret ✅
// - Production: OS env var sets JWT_SECRET=prod-secret ✅
// - Dangerous: Nothing set → uses default ❌ (caught by validation)

// ValidateConfig performs environment-specific validation
func ValidateConfig(cfg *Config, env string) error {
    switch env {
    case "production":
        return validateProductionConfig(cfg)
    case "development", "staging":
        return validateNonProductionConfig(cfg)
    default:
        return fmt.Errorf("unknown environment: %s", env)
    }
}

func validateProductionConfig(cfg *Config) error {
    // Production-specific validations
    if cfg.JWT.Secret == "your-secret-key-change-in-production" {
        return errors.New("JWT_SECRET must be changed from default in production")
    }

    if strings.Contains(cfg.Database.Host, "localhost") && cfg.Database.Host != "localhost" {
        return errors.New("production database should not use localhost")
    }

    if cfg.Logger.Level == "debug" {
        return errors.New("debug logging not recommended in production")
    }

    return nil
}

func validateNonProductionConfig(cfg *Config) error {
    // Non-production validations (more lenient)
    return nil
}
```

### 2.3 Add missing import

Add the missing import for strings:

```go
import (
    "fmt"
    "log"
    "strings"
    "errors"

    "github.com/spf13/viper"
)
```

## Step 3: Set Up Environment Variable Loading

### 3.1 Create .env.example file

Create `.env.example` in the project root:

```env
# Environment
APP_ENV=development

# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=3000
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30
SERVER_IDLE_TIMEOUT=60

# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password
DATABASE_NAME=user_management
DATABASE_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION_TIME=24

# Logger Configuration
LOGGER_LEVEL=info
LOGGER_FORMAT=json
```

**Important**: Note the `APP_ENV=development` line - this is critical for the production-first configuration pattern.

### 3.2 Create local .env file

Create `.env` file for local development:

```env
# Copy from .env.example and modify as needed
APP_ENV=development
SERVER_PORT=8080
DATABASE_PASSWORD=your_local_password
JWT_SECRET=your-local-development-secret
```

### 3.3 Update .gitignore

Add `.env` to `.gitignore`:

```gitignore
# Environment variables
.env
```

## Step 4: Create Configuration File Support

### 4.1 Create config.yaml file

Create `config.yaml` in the project root:

```yaml
server:
  host: "localhost"
  port: "3000"
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 60

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: ""
  dbname: "user_management"
  sslmode: "disable"

jwt:
  secret: "your-secret-key-change-in-production"
  expiration_time: 24

logger:
  level: "info"
  format: "json"
```

### 4.2 Create config.yaml.example

Create `config.yaml.example` as a template:

```yaml
server:
  host: "localhost"
  port: "3000"
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 60

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your_password_here"
  dbname: "user_management"
  sslmode: "disable"

jwt:
  secret: "your-secret-key-change-in-production"
  expiration_time: 24

logger:
  level: "info"
  format: "json"
```

## Step 5: Update Main Application

### 5.1 Modify cmd/main.go

Update the main.go file to use the configuration:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/boonyarit-iamsaard/user-management-api/internal/config"
    "github.com/boonyarit-iamsaard/user-management-api/internal/handler"
    "github.com/gofiber/fiber/v3"
    "github.com/joho/godotenv"
)

func main() {
    // Step 1: Environment detection and setup (production-first approach)
    env := os.Getenv("APP_ENV")

    if env == "production" {
        setupProductionEnvironment()
    } else {
        setupNonProductionEnvironment()
        // Reload env after .env loading
        env = os.Getenv("APP_ENV")
    }

    // Step 2: Load configuration
    configPath := getConfigPath()
    cfg, err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Step 3: Validate configuration based on environment
    if err := config.ValidateConfig(cfg, env); err != nil {
        log.Fatalf("Configuration validation failed for %s: %v", env, err)
    }

    // Step 4: Create and configure application
    app := fiber.New(fiber.Config{
        // Server configuration can be customized based on cfg
        ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
        IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
    })

    // Setup middleware
    handler.SetupMiddleware(app)

    // Setup routes
    handler.SetupRoutes(app)

    // Step 5: Start server
    serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
    log.Printf("🚀 Starting User Management API in %s mode", env)
    log.Printf("📡 Server listening on %s", serverAddr)

    if err := app.Listen(serverAddr); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}

func setupProductionEnvironment() {
    log.Println("🚀 Production mode - using OS environment variables only")

    // Validate critical production requirements
    requiredEnvVars := []string{
        "DATABASE_HOST",
        "DATABASE_PASSWORD",
        "JWT_SECRET",
    }

    for _, envVar := range requiredEnvVars {
        if os.Getenv(envVar) == "" {
            log.Panicf("❌ Production requires %s to be set as OS environment variable", envVar)
        }
    }
}

func setupNonProductionEnvironment() {
    log.Println("🛠️  Non-production mode - loading .env file")

    // Must have .env file for non-production
    if err := godotenv.Load(); err != nil {
        log.Panic("❌ Non-production environment requires .env file")
    }

    env := os.Getenv("APP_ENV")
    if env == "" {
        log.Panic("❌ APP_ENV must be set in .env file")
    }

    log.Printf("🔧 Loaded %s environment from .env", env)
}

// getConfigPath determines the configuration file path
func getConfigPath() string {
    // Check for custom config path from environment
    if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
        return configPath
    }

    // Default config file locations
    locations := []string{
        "config.yaml",
        "config.yml",
        "configs/config.yaml",
        "configs/config.yml",
    }

    // Check if any config file exists
    for _, location := range locations {
        if _, err := os.Stat(location); err == nil {
            return location
        }
    }

    // Return default if none found
    return "config.yaml"
}
```

## Step 6: Test Configuration Loading

### 6.1 Create a test file

Create `internal/config/config_test.go`:

```go
package config

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
    // Test with default configuration
    cfg, err := LoadConfig("config.yaml")
    require.NoError(t, err)
    assert.NotNil(t, cfg)

    // Test default values
    assert.Equal(t, "localhost", cfg.Server.Host)
    assert.Equal(t, "3000", cfg.Server.Port)
    assert.Equal(t, 5432, cfg.Database.Port)
    assert.Equal(t, "postgres", cfg.Database.User)
    assert.Equal(t, "user_management", cfg.Database.DBName)
    assert.Equal(t, "info", cfg.Logger.Level)
}

func TestLoadConfigWithEnvVars(t *testing.T) {
    // Set environment variables
    os.Setenv("SERVER_PORT", "8080")
    os.Setenv("DATABASE_PORT", "5433")
    defer func() {
        os.Unsetenv("SERVER_PORT")
        os.Unsetenv("DATABASE_PORT")
    }()

    cfg, err := LoadConfig("config.yaml")
    require.NoError(t, err)
    assert.Equal(t, "8080", cfg.Server.Port)
    assert.Equal(t, 5433, cfg.Database.Port)
}

func TestValidateProductionConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {
            name: "valid production config",
            config: Config{
                JWT: JWTConfig{Secret: "secure-production-secret"},
                Database: DatabaseConfig{Host: "prod-db.example.com"},
                Logger: LoggerConfig{Level: "warn"},
            },
            wantErr: false,
        },
        {
            name: "invalid JWT secret",
            config: Config{
                JWT: JWTConfig{Secret: "your-secret-key-change-in-production"},
                Database: DatabaseConfig{Host: "prod-db.example.com"},
                Logger: LoggerConfig{Level: "warn"},
            },
            wantErr: true,
        },
        {
            name: "debug logging in production",
            config: Config{
                JWT: JWTConfig{Secret: "secure-production-secret"},
                Database: DatabaseConfig{Host: "prod-db.example.com"},
                Logger: LoggerConfig{Level: "debug"},
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateConfig(&tt.config, "production")
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestLoadConfigFileNotFound(t *testing.T) {
    cfg, err := LoadConfig("nonexistent.yaml")
    require.NoError(t, err) // Should not error when file doesn't exist
    assert.NotNil(t, cfg)
}
```

### 6.2 Run tests

```bash
go test ./internal/config -v
```

## Step 7: Update Documentation

### 7.1 Update README.md

Add configuration section to README.md:

````markdown
## Configuration

The application supports configuration through multiple sources with production-first security:

1. **OS Environment Variables**: Production environment (highest priority)
2. **.env Files**: Development and staging environments
3. **Configuration File**: `config.yaml` (optional, for defaults)
4. **Code Defaults**: Built-in fallbacks (lowest priority)

### Environment Detection

The application uses a production-first approach:

- **Production**: Reads OS environment variables only (never loads .env files)
- **Development/Staging**: Requires `.env` file with `APP_ENV` set

### Development Setup

1. Copy the template:

```bash
cp .env.example .env
```

1. Edit `.env` with your local values:

```bash
APP_ENV=development
SERVER_PORT=3000
DATABASE_PASSWORD=your_local_password
```

1. Run the application:

```bash
go run cmd/main.go
```

### Production Setup

Set OS environment variables via your deployment method:

**Docker:**

```bash
docker run -e APP_ENV=production -e DATABASE_URL=... -e JWT_SECRET=... myapp
```

**Kubernetes:**

```yaml
env:
  - name: APP_ENV
    value: "production"
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: url
```

**Systemd Service:**

```bash
Environment=APP_ENV=production
Environment=DATABASE_URL=postgresql://...
```

### Required Environment Variables

**Production:**

- `APP_ENV=production`
- `DATABASE_HOST`
- `DATABASE_PASSWORD`
- `JWT_SECRET`

**Development/Staging:**

- All variables in `.env.example`

### Configuration Priority

1. OS Environment Variables (production) or .env file (development)
2. Configuration file (`config.yaml`)
3. Code defaults
````

## Step 8: Verification

### 8.1 Test the application

Run the application with different configuration scenarios:

```bash
# Test with default config
go run cmd/main.go

# Test with environment variables
SERVER_PORT=8080 go run cmd/main.go

# Test with custom config file
CONFIG_PATH=custom.yaml go run cmd/main.go
```

### 8.2 Verify configuration loading

Check that the application starts and uses the correct configuration values by observing the log output.

## Troubleshooting

### Common Issues

1. **"Non-production environment requires .env file"**: Create `.env` file with required variables
2. **"Production requires X to be set"**: Set OS environment variables before starting
3. **"Configuration validation failed"**: Check production-specific requirements (JWT secret, logging level)
4. **"JWT_SECRET must be changed from default in production"**: Set a secure JWT_SECRET environment variable
5. **Config file not found**: Optional - application works with just environment variables
6. **Invalid YAML syntax**: Validate YAML syntax using online validators
7. **Environment variables not working**: Ensure no `export` commands in `.env` file

### JWT Secret Override Scenarios

| Scenario    | How it's Set                               | Result             | Security               |
| ----------- | ------------------------------------------ | ------------------ | ---------------------- |
| Development | `.env` file: `JWT_SECRET=dev-secret`       | Uses dev-secret    | ✅ Safe                |
| Production  | OS env: `export JWT_SECRET=prod-secret`    | Uses prod-secret   | ✅ Safe                |
| Config Only | `config.yaml`: `jwt.secret: config-secret` | Uses config-secret | ⚠️ Less secure         |
| Nothing Set | No variables at all                        | Uses default value | ❌ DANGEROUS (blocked) |

**Important**: The production validation will prevent startup if the default JWT secret is used, ensuring this security vulnerability never reaches production.

### Debug Configuration

Add debug logging to see what configuration values are being loaded:

```go
log.Printf("Loaded configuration: %+v", cfg)
```

## Production Deployment Checklist

Before deploying to production, ensure:

- [ ] `APP_ENV=production` is set as OS environment variable
- [ ] All required production environment variables are set
- [ ] JWT_SECRET uses a strong, unique value (not the default)
- [ ] Database credentials are properly configured
- [ ] No `.env` file is present in production
- [ ] Application starts successfully and validates configuration

## Security Best Practices

1. **Never commit `.env` files** to version control
2. **Use different secrets** for each environment
3. **Rotate secrets regularly** in production
4. **Use secret management systems** (AWS Secrets Manager, Kubernetes Secrets)
5. **Validate configuration at startup** to fail fast on misconfigurations

## Next Steps

After completing this task:

1. The application has production-ready configuration management
2. Environment variables follow Go community standards
3. Security-first approach prevents .env file usage in production
4. Automatic validation ensures production requirements are met
5. Clear separation between development and production environments
6. Ready for production deployment with proper environment variable management

Proceed to **Task 2.1: Database Connection Setup** to implement database connectivity using the production-ready configuration values.
