# 02. Environment Variables: Deep Dive & Best Practices

## The Problem

Environment variables are the standard for configuration in production systems, but there are many nuances, misconceptions, and best practices that aren't immediately obvious.

## Core Concepts

### What Environment Variables Actually Are

```bash
# These are OS-level variables, not application-level
export DATABASE_URL=postgresql://user:pass@localhost:5432/dbname
export JWT_SECRET=super-secret-key

# Any application can read these:
os.Getenv("DATABASE_URL")  // Go
process.env.DATABASE_URL   // Node.js
ENV['DATABASE_URL']        // Ruby
```

**Key Insight**: Environment variables exist at the operating system level, not application level. Any application running on the system can access them.

### The .env File Misconception

```bash
# .env file is NOT automatically available to applications
DATABASE_URL=postgresql://user:pass@localhost:5432/dbname
JWT_SECRET=dev-secret

# This file must be explicitly loaded and its contents
# converted to OS environment variables
```

**The Reality**: `.env` files are just text files. They only become "environment variables" when a library reads them and calls `os.Setenv()`.

## Environment Variable Loading Patterns

### Pattern 1: Simple Direct Loading (Production-Ready)

```go
func main() {
    // Direct environment variable access
    dbURL := os.Getenv("DATABASE_URL")
    jwtSecret := os.Getenv("JWT_SECRET")

    // Validate required variables
    if dbURL == "" {
        log.Fatal("DATABASE_URL is required")
    }

    startApp(dbURL, jwtSecret)
}
```

**Pros**:

- No external dependencies
- Works everywhere
- Production-friendly

**Cons**:

- No defaults
- No type conversion
- Manual validation required

### Pattern 2: Viper-Based Loading (Structured)

```go
func LoadConfig() (*Config, error) {
    viper.AutomaticEnv()  // Reads all OS environment variables
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

// Config struct
type Config struct {
    DatabaseURL string `mapstructure:"database_url"`
    JWTSecret   string `mapstructure:"jwt_secret"`
    ServerPort  int    `mapstructure:"server_port"`
}
```

**Pros**:

- Type conversion
- Structured configuration
- Validation support
- Default values

**Cons**:

- External dependency
- Learning curve

### Pattern 3: Hybrid Approach (Most Flexible)

```go
func main() {
    // Load .env in development only
    if os.Getenv("APP_ENV") != "production" {
        _ = godotenv.Load()
    }

    // Use Viper for structured config
    cfg, err := LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    startApp(cfg)
}
```

## Environment Variable Naming Conventions

### Industry Standards

#### 12-Factor App Style

```bash
DATABASE_URL=postgresql://user:pass@host:port/dbname
REDIS_URL=redis://host:port
SECRET_KEY_BASE=your-secret-key
PORT=3000
RACK_ENV=production
```

#### Descriptive Naming

```bash
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=secret
DATABASE_NAME=myapp

SERVER_HOST=0.0.0.0
SERVER_PORT=8080

JWT_SECRET=your-jwt-secret
JWT_EXPIRATION=24h
```

#### Service-Specific (Less Common)

```bash
MYAPP_DATABASE_HOST=localhost
MYAPP_SERVER_PORT=3000
USER_MANAGEMENT_API_JWT_SECRET=secret
```

### Best Practices

1. **Use SCREAMING_SNAKE_CASE**
2. **Be descriptive but concise**
3. **Group related variables with prefixes**
4. **Avoid unnecessary prefixes for common variables**
5. **Use standard names for common concepts**

```bash
# ✅ GOOD - Standard and descriptive
DATABASE_URL=postgresql://...
SERVER_PORT=8080
JWT_SECRET=secret

# ❌ AVOID - Unnecessary prefix
MYAPP_DATABASE_URL=postgresql://...
MYAPP_SERVER_PORT=8080
MYAPP_JWT_SECRET=secret

# ✅ GOOD - Grouped with purpose
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=redis-pass
REDIS_DB=0

# ✅ GOOD - Service-specific when needed
STRIPE_SECRET_KEY=sk_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

## Environment Variable Types

### 1. Required Variables

```bash
# Application cannot start without these
DATABASE_URL=postgresql://...
JWT_SECRET=super-secret-key
```

### 2. Optional with Defaults

```bash
# Application has sensible defaults
SERVER_PORT=3000  # Default: 3000
LOG_LEVEL=info    # Default: info
CACHE_TTL=3600    # Default: 3600
```

### 3. Feature Flags

```bash
# Control application behavior
ENABLE_METRICS=true
ENABLE_DEBUG_LOGGING=false
MAINTENANCE_MODE=false
```

### 4. External Service Configuration

```bash
# Third-party service settings
STRIPE_API_KEY=sk_test_...
SENDGRID_API_KEY=SG.xyz...
AWS_REGION=us-west-2
```

## Security Considerations

### 1. Never Commit Secrets

```bash
# ❌ NEVER commit to version control
DATABASE_PASSWORD=super-secret-password
JWT_SECRET=production-jwt-key

# ✅ Use .env files locally, OS env vars in production
# .env (local, not committed)
DATABASE_PASSWORD=dev-password

# Production (set by infrastructure)
export DATABASE_PASSWORD=super-secure-production-password
```

### 2. Use Different Values Per Environment

```bash
# Development
JWT_SECRET=dev-secret-not-for-production

# Staging
JWT_SECRET=staging-secret-different-from-prod

# Production
JWT_SECRET=super-secure-production-secret-32-chars-minimum
```

### 3. Validate at Startup

```go
func validateConfig(cfg *Config, env string) error {
    if cfg.JWT.Secret == "your-secret-key-change-in-production" {
        return errors.New("JWT_SECRET must be changed from default")
    }

    if len(cfg.JWT.Secret) < 32 {
        return errors.New("JWT_SECRET must be at least 32 characters")
    }

    if env == "production" && cfg.Database.Host == "localhost" {
        return errors.New("production database cannot be localhost")
    }

    return nil
}
```

### 4. Mask in Logs

```go
func logConfig(cfg *Config) {
    log.Printf("Database: %s", maskURL(cfg.Database.URL))
    log.Printf("JWT: %s", maskSecret(cfg.JWT.Secret))
}

func maskURL(url string) string {
    if len(url) > 30 {
        return url[:30] + "***"
    }
    return "***"
}
```

## Environment Variable Management Strategies

### 1. Docker Environment

```bash
# Command line
docker run -e DATABASE_URL=postgresql://... -e JWT_SECRET=... myapp

# Environment file
docker run --env-file .env.production myapp

# Docker Compose
services:
  app:
    environment:
      - DATABASE_URL=postgresql://...
      - JWT_SECRET=super-secret
```

### 2. Kubernetes ConfigMaps and Secrets

```yaml
# ConfigMap - Non-sensitive data
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_ENV: "production"
  LOG_LEVEL: "warn"
  SERVER_PORT: "8080"

---
# Secret - Sensitive data
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  DATABASE_URL: cG9zdGdyZXNxbDovLy4uLg== # base64 encoded
  JWT_SECRET: c3VwZXItc2VjdXJlLWp3dC1rZXk=
```

### 3. CI/CD Environment Variables

```yaml
# GitHub Actions
- name: Run Tests
  env:
    DATABASE_URL: postgresql://test:test@localhost:5432/test
    JWT_SECRET: test-jwt-secret
  run: go test ./...

# GitLab CI
variables:
  DATABASE_URL: postgresql://test:test@postgres:5432/test
  JWT_SECRET: test-jwt-secret
```

### 4. Cloud Provider Secrets

```bash
# AWS Systems Manager Parameter Store
aws ssm put-parameter \
  --name "/myapp/prod/database-url" \
  --value "postgresql://..." \
  --type SecureString

# AWS Secrets Manager
aws secretsmanager create-secret \
  --name "myapp/database-credentials" \
  --secret-string '{"username":"user","password":"pass"}'
```

## Common Pitfalls and Solutions

### Pitfall 1: Shell Commands in .env Files

```bash
# ❌ WRONG - .env files don't support shell commands
export DATABASE_URL=postgresql://user:pass@host:port/db
export JWT_SECRET=$(openssl rand -base64 32)

# ✅ CORRECT - Simple key=value pairs
DATABASE_URL=postgresql://user:pass@host:port/db
JWT_SECRET=your-secret-key
```

### Pitfall 2: Assuming Environment Variables Exist

```go
// ❌ WRONG - Can cause runtime panics
dbURL := os.Getenv("DATABASE_URL")
parseDatabaseURL(dbURL)  // Panics if dbURL is empty

// ✅ CORRECT - Validate first
dbURL := os.Getenv("DATABASE_URL")
if dbURL == "" {
    return errors.New("DATABASE_URL is required")
}
parseDatabaseURL(dbURL)
```

### Pitfall 3: Hardcoding Environment Names

```go
// ❌ WRONG - Brittle
if os.Getenv("RAILS_ENV") == "production" {
    // Rails-specific
}

// ✅ CORRECT - App-specific
if os.Getenv("APP_ENV") == "production" {
    // App-specific environment detection
}
```

### Pitfall 4: Not Using Type Conversion

```go
// ❌ WRONG - Always strings
port := os.Getenv("SERVER_PORT")  // "3000"
server.Listen(":" + port)

// ✅ CORRECT - Convert to proper types
port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
server.Listen(fmt.Sprintf(":%d", port))

// ✅ EVEN BETTER - Use Viper for type conversion
viper.SetDefault("server.port", 3000)
port := viper.GetInt("server.port")
```

## Advanced Patterns

### 1. Environment Variable Precedence

```go
func LoadConfig() (*Config, error) {
    // Load in order of precedence:

    // 1. Code defaults (lowest)
    viper.SetDefault("server.port", 3000)

    // 2. Config file (medium)
    viper.SetConfigFile("config.yaml")
    viper.ReadInConfig()

    // 3. Environment variables (highest)
    viper.AutomaticEnv()

    return &Config{}, viper.Unmarshal(&config)
}
```

### 2. Multi-Environment Configuration

```go
func LoadConfigForEnvironment(env string) (*Config, error) {
    // Load base config
    viper.SetConfigName("config.default")
    viper.ReadInConfig()

    // Override with environment-specific config
    viper.SetConfigName(fmt.Sprintf("config.%s", env))
    viper.ReadInConfig()

    // Override with environment variables
    viper.AutomaticEnv()

    return &Config{}, viper.Unmarshal(&config)
}
```

### 3. Environment Variable Validation

```go
type Config struct {
    DatabaseURL string `validate:"required,url"`
    JWTSecret   string `validate:"required,min=32"`
    ServerPort  int    `validate:"required,min=1,max=65535"`
}

func (c *Config) Validate() error {
    validate := validator.New()
    return validate.Struct(c)
}
```

## Testing with Environment Variables

### 1. Unit Tests

```go
func TestLoadConfig(t *testing.T) {
    // Set test environment variables
    os.Setenv("DATABASE_URL", "postgresql://test:test@localhost:5432/test")
    os.Setenv("JWT_SECRET", "test-jwt-secret-32-chars-long")
    defer func() {
        os.Unsetenv("DATABASE_URL")
        os.Unsetenv("JWT_SECRET")
    }()

    cfg, err := LoadConfig()
    require.NoError(t, err)
    assert.Equal(t, "postgresql://test:test@localhost:5432/test", cfg.DatabaseURL)
    assert.Equal(t, "test-jwt-secret-32-chars-long", cfg.JWT.Secret)
}
```

### 2. Integration Tests

```go
func TestProductionEnvironment(t *testing.T) {
    // Test production configuration
    os.Setenv("APP_ENV", "production")
    os.Setenv("DATABASE_URL", "postgresql://prod:pass@prod-db:5432/prod")
    os.Setenv("JWT_SECRET", "super-secure-production-key-32-chars")
    defer cleanup()

    cfg, err := LoadConfig()
    require.NoError(t, err)

    err = cfg.Validate("production")
    assert.NoError(t, err)
}
```

## Real-World Examples

### 1. Web Application

```bash
# Required
DATABASE_URL=postgresql://user:pass@host:port/dbname
JWT_SECRET=super-secret-key-32-chars
REDIS_URL=redis://host:port

# Optional
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
LOG_LEVEL=info
APP_ENV=production
```

### 2. Background Worker

```bash
# Required
DATABASE_URL=postgresql://user:pass@host:port/dbname
REDIS_URL=redis://host:port

# Worker-specific
WORKER_CONCURRENCY=10
QUEUE_NAME=default
JOB_TIMEOUT=300
```

### 3. Microservice

```bash
# Service discovery
SERVICE_NAME=user-api
SERVICE_PORT=8080
CONSUL_HOST=consul.example.com

# External dependencies
USER_SERVICE_URL=http://user-service:8080
NOTIFICATION_SERVICE_URL=http://notification-service:8080
```

## Related Learning

- **[01-configuration-patterns.md](./01-configuration-patterns.md)** - Configuration management patterns
- **[../guides/01-configuration-management.md](../guides/01-configuration-management.md)** - Implementation guide
- **[External Resources](#external-resources)** - Additional reading

## External Resources

- [12-Factor App - Configuration](https://12factor.net/config)
- [Docker Environment Variables](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e---env---env-file)
- [Kubernetes ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)
- [Viper Documentation](https://github.com/spf13/viper)

## Key Takeaways

1. **Environment variables are OS-level**, not application-level
2. **.env files are development convenience**, never used in production
3. **Use standard naming conventions** (SCREAMING_SNAKE_CASE, descriptive names)
4. **Validate at startup** to fail fast on misconfigurations
5. **Never commit secrets** to version control
6. **Use different values per environment**
7. **Mask sensitive data in logs**

---

_This entry captures the deep understanding of environment variables from basic concepts to production-ready patterns. The key is understanding that environment variables are a system-level configuration mechanism, not an application feature._
