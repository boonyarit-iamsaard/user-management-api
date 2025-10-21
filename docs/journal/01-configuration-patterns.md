# 01. Configuration Management Patterns

## The Problem

When implementing Task 1.3 (Configuration Management), we needed to create a system that works seamlessly across different environments while maintaining security and following Go community best practices.

## Initial Misconceptions

### Misconception 1: "UMA\_" Prefix is Standard

**What we thought**: Using a project-specific prefix like `UMA_` for environment variables was a good practice.

**Reality**: The Go community typically uses simple, descriptive names without prefixes. Common patterns are:

- `DATABASE_URL`, `PORT`, `SECRET_KEY` (12-factor app style)
- `SERVER_HOST`, `DATABASE_PASSWORD` (descriptive names)
- Rarely: `MYAPP_DATABASE_HOST` (full project name prefix)

### Misconception 2: ".env Files Work Everywhere"

**What we thought**: `.env` files are loaded automatically by applications and work in production.

**Reality**: `.env` files are:

- **Development convenience only** - Never used in production
- **Not automatically loaded** - Require libraries like `godotenv`
- **Security risk** - Should never be committed to version control
- **Infrastructure-managed** - Production uses real OS environment variables

### Misconception 3: "Viper Loads .env Files"

**What we thought**: Viper automatically reads `.env` files.

**Reality**: Viper only:

- Reads config files (YAML, JSON, etc.)
- Reads existing OS environment variables
- **Does NOT** parse `.env` files automatically

## The Insight: Production-First Configuration Pattern

### The "Aha!" Moment

The key insight was to **check the environment FIRST** before any file operations:

```go
// WRONG: Load files first, then check environment
godotenv.Load()  // This runs in production!
cfg := loadConfig()

// RIGHT: Check environment first
if os.Getenv("APP_ENV") == "production" {
    // Never load .env files in production
    setupProductionEnvironment()
} else {
    // Only load .env in non-production
    godotenv.Load()
    setupNonProductionEnvironment()
}
```

### Why This Pattern Works

1. **Security First**: Production never touches `.env` files
2. **Clear Intent**: Environment is explicitly declared
3. **Fail Fast**: Errors caught immediately
4. **Infrastructure Ready**: Works with Docker/K8s/CI-CD

## Real-World Flow

### Development Environment

```bash
# .env file exists
APP_ENV=development
DATABASE_PASSWORD=dev_password

# Flow:
1. os.Getenv("APP_ENV") == "" (not "production") ✅
2. godotenv.Load() → sets OS environment ✅
3. Load config with dev values ✅
```

### Production Environment

```bash
# No .env file, infra team sets OS env vars
APP_ENV=production
DATABASE_URL=postgresql://prod_user:secure_pass@prod_db:5432/proddb

# Flow:
1. os.Getenv("APP_ENV") == "production" ✅
2. Skip .env loading ✅
3. Load config from OS environment ✅
```

## Override Behavior: The Priority Pyramid

```text
Environment Variables (Highest Priority)
    ↓
Configuration File (Medium Priority)
    ↓
Code Defaults (Lowest Priority)
```

### JWT Secret Example

```go
// Code default (lowest priority)
viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")

// Override scenarios:
```

| Source      | Example                       | Result       | Security               |
| ----------- | ----------------------------- | ------------ | ---------------------- |
| Environment | `JWT_SECRET=super-secure-key` | Uses env var | ✅ Safe                |
| Config File | `jwt.secret: "config-secret"` | Uses config  | ⚠️ Less secure         |
| Nothing Set | (no variables)                | Uses default | ❌ Dangerous (blocked) |

## Key Patterns Discovered

### 1. Environment Detection Pattern

```go
func SetupEnvironment() string {
    env := os.Getenv("APP_ENV")

    if env == "production" {
        return setupProductionEnvironment()
    } else {
        return setupNonProductionEnvironment()
    }
}
```

### 2. Validation Pattern

```go
func (c *Config) Validate(env string) error {
    switch env {
    case "production":
        return c.validateProduction()  // Strict validation
    default:
        return c.validateNonProduction()  // Lenient validation
    }
}
```

### 3. DATABASE_URL Support Pattern

```go
// Support both DATABASE_URL and individual variables
if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
    parseDatabaseURL(dbURL)  // Extract individual settings
}
```

## Common Pitfalls and Solutions

### Pitfall 1: Using `export` in .env files

```bash
# ❌ WRONG in .env
export DATABASE_PASSWORD=secret

# ✅ CORRECT in .env
DATABASE_PASSWORD=secret
```

### Pitfall 2: Not validating production configuration

```go
// ❌ DANGEROUS - No validation
cfg := loadConfig()
startServer(cfg)

// ✅ SAFE - Production validation
if err := cfg.Validate("production"); err != nil {
    log.Fatal("Configuration validation failed:", err)
}
```

### Pitfall 3: Assuming .env files exist

```go
// ❌ WRONG - Assumes .env exists
godotenv.Load()  // Panics in production

// ✅ RIGHT - Graceful handling
if err := godotenv.Load(); err != nil {
    log.Println("No .env file - using OS environment variables")
}
```

## Production Deployment Checklist

### Infrastructure Team Responsibilities

```bash
# Must set these OS environment variables:
export APP_ENV=production
export DATABASE_URL=postgresql://user:pass@host:port/dbname
export JWT_SECRET=super-secure-key-min-32-chars
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
```

### Application Team Responsibilities

- Validate configuration at startup
- Never commit `.env` files
- Use different secrets per environment
- Implement proper logging (mask sensitive data)

## Related Files and References

### Implementation Guide

- **[Configuration Management Guide](../guides/01-configuration-management.md)** - Step-by-step implementation

### Code Files

- `internal/config/config.go` - Configuration structure and loading
- `internal/config/environment.go` - Environment detection and setup
- `cmd/main.go` - Application startup with environment handling

### External Resources

- [12-Factor App - Config](https://12factor.net/config)
- [Viper Documentation](https://github.com/spf13/viper)
- [godotenv Documentation](https://github.com/joho/godotenv)

## Key Takeaways

1. **Environment detection comes first** - Check `APP_ENV` before any file operations
2. **Production never uses .env files** - Always use OS environment variables
3. **Standard Go naming** - Use descriptive names without prefixes
4. **Validate configuration** - Fail fast on misconfigurations
5. **Security by default** - Default values should be safe and obvious

## Future Topics to Explore

- [ ] Secret rotation strategies
- [ ] Configuration hot-reloading
- [ ] Multi-environment config files (config.dev.yml, config.prod.yml)
- [ ] Integration with secret management systems (AWS Secrets Manager, Vault)
- [ ] Configuration validation libraries and best practices

---

_This entry captures the learning journey from misconceptions to production-ready configuration management patterns. The production-first approach ensures security while maintaining flexibility for development workflows._
