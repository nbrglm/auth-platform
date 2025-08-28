package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nbrglm/auth-platform/opts"
	"github.com/nbrglm/auth-platform/utils"
	"gopkg.in/yaml.v2"
)

// Type aliases for enums
type PasswordHashingAlgorithm string

const (
	// Bcrypt is the default password hashing algorithm
	BcryptPasswordHashingAlgorithm PasswordHashingAlgorithm = "bcrypt"
	// Argon2id is a more secure password hashing algorithm
	Argon2idPasswordHashingAlgorithm PasswordHashingAlgorithm = "argon2id"
)

// Variables for different configurations
var (
	NBRGLMBranding = true  // Default to true, can be set to false in the config file
	Multitenancy   = false // Default to false, can be set to true in the config file
	Public         *PublicConfig
	Server         *ServerConfig
	Observability  *ObservabilityConfig
	Password       *PasswordConfig
	JWT            *JWTConfig
	Notifications  *NotificationsConfig
	Branding       *BrandingConfig
	Security       *SecurityConfig
	Stores         *StoresConfig

	// Admins is a list of credentials for admin users
	Admins []AdminCredential
)

func Environment() string {
	if opts.Debug {
		return "development"
	}
	return "production"
}

// ObservabilityConfig holds the configuration for observability features
//
// Used for logs, traces and metrics (Recommended: SigNoz)
type ObservabilityConfig struct {
	// The log level for the auth platform.
	LogLevel string `json:"logLevel" yaml:"logLevel" validate:"required,oneof=debug info warn error dpanic panic fatal"`

	// OtelExporterEndpoint is the endpoint for the OpenTelemetry exporter
	OtelExporterEndpoint string `json:"otelExporterEndpoint" yaml:"otelExporterEndpoint" validate:"omitempty,url"`

	// OtelExporterProtocol is the protocol of the OpenTelemetry exporter
	OtelExporterProtocol string `json:"otelExporterProtocol" yaml:"otelExporterProtocol" validate:"required,oneof=http/protobuf grpc stdout"`
}

// PublicConfig holds the public server configuration options
//
// Note: If hosting NAP at "auth.example.com", then set the values as follows:
//
// domain: example.com
// subDomain: auth
//
// Since, the OIDC issuer URLs will be like:
// https://auth.example.com/.well-known/openid-configuration
//
// Your DNS should be configured to point to the server's IP address,
// for example:
// *.auth.example.com should point to the server's IP address.
type PublicConfig struct {
	// Scheme for public URLs (http/https)
	//
	// Use "https" even if TLS termination is handled by reverse proxy
	Scheme string `json:"scheme" yaml:"scheme" validate:"required,oneof=http https"`

	// Domain is the domain at which NAP is being hosted.
	//
	// Eg. if hosting NAP at auth.example.com, provide "example.com" here.
	//
	// This setting is used as the Domain for all Cookies (except Refresh Token Cookie, that is only set at "subdomain.domain")
	Domain string `json:"domain" yaml:"domain" validate:"required"`

	// Subdomain at which NAP is being hosted.
	//
	// Eg. If hosting NAP at auth.example.com, provide "auth" here.
	SubDomain string `json:"subDomain" yaml:"subDomain" validate:"required"`

	DebugBaseURL string `json:"debugBaseURL,omitempty" yaml:"debugBaseURL,omitempty" validate:"omitempty,url"`
}

func (p *PublicConfig) GetBaseURL() string {
	if opts.Debug && strings.TrimSpace(p.DebugBaseURL) != "" {
		return p.DebugBaseURL
	}
	return fmt.Sprintf("%s://%s.%s", p.Scheme, p.SubDomain, p.Domain)
}

func (p *PublicConfig) GetTenantBaseURL(tenant string) string {
	return fmt.Sprintf("%s://%s.%s.%s", p.Scheme, tenant, p.SubDomain, p.Domain)
}

type RedirectsConfig struct {
	// Domains is a list of that will be redirected to if the return_to parameter is provided.
	// If the return_to parameter is not provided, the user will be redirected to the redirects page.
	Domains []string `json:"domains" yaml:"domains" validate:"omitempty,dive,domain"`
}

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	// Host interface to listen on, Default: localhost
	Host string `json:"host" yaml:"host" validate:"required"`
	// Port to listen on, Default: 3360, must be between 1024 and 65535
	Port string `json:"port" yaml:"port" validate:"required"`

	// Instance ID for the metrics, logs, traces....
	InstanceID string `json:"instanceID" yaml:"instanceID" validate:"required"`

	// TLSConfig contains the TLS configuration for the server
	//
	// If TLS is enabled, the server will listen on HTTPS, else HTTP.
	TLSConfig *TLSConfig `json:"tls,omitempty" yaml:"tls,omitempty" validate:"omitempty"`
}

type AdminCredential struct {
	// Email of the admin user
	Email string `json:"email" yaml:"email" validate:"required,email"`
	// Password Hash of the admin user, Bcrypt only
	PasswordHash string `json:"passwordHash" yaml:"passwordHash" validate:"required"`
}

// TLSConfig holds the TLS configuration for the server
type TLSConfig struct {
	// CertFile is the path to the TLS certificate file
	CertFile string `json:"certFile" yaml:"certFile" validate:"required,file"`

	// KeyFile is the path to the TLS private key file
	KeyFile string `json:"keyFile" yaml:"keyFile" validate:"required,file"`

	// CAFile is the path to the TLS CA certificate file
	// If provided, the server will use this CA to verify client certificates
	CAFile *string `json:"caFile,omitempty" yaml:"caFile,omitempty" validate:"omitempty,file"`
}

// PasswordConfig holds the configuration for password policies
type PasswordConfig struct {
	// The algorithm used for hashing passwords, default: bcrypt
	Algorithm PasswordHashingAlgorithm `json:"algorithm" yaml:"algorithm" validate:"required,oneof=bcrypt argon2id"`

	// Configuration for Bcrypt password hashing
	Bcrypt BcryptConfig `json:"bcrypt" yaml:"bcrypt" validate:"required_if=Algorithm bcrypt"`

	// Configuration for Argon2id password hashing
	Argon2id Argon2idConfig `json:"argon2id" yaml:"argon2id" validate:"required_if=Algorithm argon2id"`
}

// BcryptConfig holds the configuration for Bcrypt password hashing
type BcryptConfig struct {
	// The cost factor for Bcrypt hashing (default is 10)
	Cost int `json:"cost" yaml:"cost" validate:"required,min=10,max=31"`
}

// Argon2idConfig holds the configuration for Argon2id password hashing
type Argon2idConfig struct {
	// Memory in KiB (16MB minimum, and default)
	Memory int `json:"memory" yaml:"memory" validate:"required,min=16384"`

	// Number of iterations (3-4 recommended, default 3)
	Iterations int `json:"iterations" yaml:"iterations" validate:"required,min=1"`

	// Number of parallel threads (1-4 recommended, default 1)
	Parallelism int `json:"parallelism" yaml:"parallelism" validate:"required,min=1"`

	// Salt length in bytes (16-32 recommended, default 16)
	SaltLength int `json:"saltLength" yaml:"saltLength" validate:"required,min=16,max=32"`

	// Key length in bytes (32-64 recommended, default 32)
	KeyLength int `json:"keyLength" yaml:"keyLength" validate:"required,min=32,max=64"`
}

// JWTConfig holds the configuration for JWT tokens
type JWTConfig struct {
	// Session token expiration time in seconds (default: 1hr, 3600)
	SessionTokenExpiration int `json:"sessionTokenExpiration" yaml:"sessionTokenExpiration" validate:"required,min=60"`

	// Refresh token expiration time in seconds (default: 30d, 2592000)
	RefreshTokenExpiration int `json:"refreshTokenExpiration" yaml:"refreshTokenExpiration" validate:"required,min=86400"`

	// Path to the private key file for RS256 algorithm
	PrivateKeyFile string `json:"privateKeyFile" yaml:"privateKeyFile" validate:"required,file"`

	// Path to the public key file for RS256 algorithm
	PublicKeyFile string `json:"publicKeyFile" yaml:"publicKeyFile" validate:"required,file"`

	// Audiences claim for the JWT.
	//
	// NOTE: This is NOT for OIDC. This is for the session tokens! OIDC configuration is stored in the DB per tenant.
	//
	// NOTE: Do not add the subdomain or domain you added in the Public Config here, it is automatically added and will be duplicated if you add.
	// Eg. If domain has value example.com and subdomain has value auth,
	// then auth.example.com, and example.com are automatically added as audience claims.
	Audiences []string `json:"audiences" yaml:"audiences" validate:"required,dive,required"`
}

type NotificationsConfig struct {
	// Email configuration for sending notifications
	Email *EmailNotificationConfig `json:"email,omitempty" yaml:"email,omitempty" validate:"omitempty"`

	// SMS Configuration for sending notifications
	SMS *SMSNotificationConfig `json:"sms,omitempty" yaml:"sms,omitempty" validate:"omitempty"`
}

// EmailNotificationConfig holds the configuration for email notifications.
type EmailNotificationConfig struct {
	// Provider is the email provider to use for sending emails.
	Provider string `json:"provider" yaml:"provider" validate:"required,oneof=ses sendgrid smtp"`

	// SendGridConfig holds the configuration for SendGrid email provider.
	SendGrid *SendGridProviderConfig `json:"sendgrid,omitempty" yaml:"sendgrid,omitempty" validate:"omitempty,required_if=Provider sendgrid"`

	// SMTPConfig holds the configuration for SMTP email provider.
	SMTP *SMTPProviderConfig `json:"smtp,omitempty" yaml:"smtp,omitempty" validate:"omitempty,required_if=Provider smtp"`

	// SESConfig holds the configuration for AWS SES email provider.
	SES *SESProviderConfig `json:"ses,omitempty" yaml:"ses,omitempty" validate:"omitempty,required_if=Provider ses"`

	// TemplatesDir is the directory where email templates are stored.
	TemplatesDir *string `json:"templatesDir,omitempty" yaml:"templatesDir,omitempty" validate:"omitempty,file"`
}

// SendGridProviderConfig holds the configuration for SendGrid email provider.
type SendGridProviderConfig struct {
	// API key for SendGrid
	APIKey string `json:"apiKey" yaml:"apiKey" validate:"required"`

	// Email address from which notifications are sent
	FromAddress string `json:"fromAddress" yaml:"fromAddress" validate:"required,email"`

	// Name associated with the FromAddress, optional
	FromName *string `json:"fromName,omitempty" yaml:"fromName,omitempty"`
}

// SMTPProviderConfig holds the configuration for SMTP email provider.
//
// Note: SMTP is a generic email provider that can be used with any SMTP server.
// IT DOES NOT SUPPORT RETRIES, OR ANY ADVANCED FEATURES.
// It is recommended to use a more robust provider like SendGrid or AWS SES for production use.
type SMTPProviderConfig struct {
	// SMTP server host
	Host string `json:"host" yaml:"host" validate:"required"`

	// SMTP server port
	Port string `json:"port" yaml:"port"`

	// Email address from which notifications are sent
	FromAddress string `json:"fromAddress" yaml:"fromAddress" validate:"required,email"`

	// Password for SMTP authentication
	Password string `json:"password" yaml:"password" validate:"required"`
}

// SESProviderConfig holds the configuration for AWS SES email provider.
type SESProviderConfig struct {
	// AWS region where SES is configured
	Region string `json:"region" yaml:"region" validate:"required"`

	// AWS Access Key ID for SES
	AccessKeyID string `json:"accessKeyID" yaml:"accessKeyID" validate:"required"`

	// AWS Secret Access Key for SES
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey" validate:"required"`

	// Email address from which notifications are sent
	FromAddress string `json:"fromAddress" yaml:"fromAddress" validate:"required,email"`

	// Optional: Name associated with the FromAddress, used in email headers, if not provided, AppName is used.
	FromName *string `json:"fromName,omitempty" yaml:"fromName,omitempty"`
}

// SMSNotificationConfig holds the configuration for SMS notifications.
type SMSNotificationConfig struct {
	// TODO: Implement SMSNotificationConfig

	// TemplatesDir is the directory where SMS templates are stored.
	TemplatesDir *string `json:"templatesDir,omitempty" yaml:"templatesDir,omitempty" validate:"omitempty,file"`
}

// BrandingConfig holds the configuration for branding elements such as names.
type BrandingConfig struct {
	// AppName is the name of the application, used in various places like email templates, UI, etc.
	AppName string `json:"appName" yaml:"appName" validate:"required"`

	// CompanyName is the name of the company, used in emails, UI, etc.
	CompanyName string `json:"companyName" yaml:"companyName" validate:"required"`

	// CompanyNameShort is a short version of the company name, used in places where space is limited.
	CompanyNameShort string `json:"companyNameShort" yaml:"companyNameShort" validate:"required"`

	// SupportURL is the URL for support, used in emails, UI, etc.
	SupportURL string `json:"supportURL" yaml:"supportURL" validate:"required,url"`
}

// SecurityConfig holds the security-related configurations for the application.
type SecurityConfig struct {
	// Enable or disable audit logs.
	EnableAuditLogs bool `json:"enableAuditLogs" yaml:"enableAuditLogs"`

	// The list of API Keys which are allowed to access the API endpoints.
	// Requests without an API key, or with a key not specified here, will be denied with 401.
	APIKeys []APIKeyConfig `json:"apiKeys" yaml:"apiKeys" validate:"required,dive"`

	// CORS configuration for the application.
	CORS CORSConfig `json:"cors" yaml:"cors" validate:"required"`

	// Rate limiting configuration.
	RateLimit RateLimitConfig `json:"rateLimit" yaml:"rateLimit" validate:"required"`
}

type APIKeyConfig struct {
	// The name of the API Key
	Name string `json:"name" yaml:"name" validate:"required"`

	// Description of the key
	Description string `json:"description" yaml:"description" validate:"required"`

	// The key itself. No restriction on the length, but keep it sensible please.
	Key string `json:"key" yaml:"key" validate:"required"`
}

// CORSConfig holds the configuration for Cross-Origin Resource Sharing (CORS).
type CORSConfig struct {
	// AllowedOrigins is a list of origins that are allowed to access the resources.
	AllowedOrigins []string `json:"allowedOrigins" yaml:"allowedOrigins" validate:"dive,url"`

	// AllowedMethods is a list of HTTP methods that are allowed.
	AllowedMethods []string `json:"allowedMethods" yaml:"allowedMethods" validate:"required,dive,oneof=GET POST PUT DELETE OPTIONS PATCH HEAD"`

	// AllowedHeaders is a list of headers that are allowed in requests.
	AllowedHeaders []string `json:"allowedHeaders" yaml:"allowedHeaders" validate:"required"`
}

// RateLimitConfig holds the configuration for rate limiting the API.
type RateLimitConfig struct {
	// Rate limit for API requests.
	// Format: "R-U", where R is requests and U is the time unit (s - per second, m - per minute, h - per hour, d - per day)
	Rate string `json:"rate" yaml:"rate" validate:"required"`
}

// StoresConfig holds the configuration for the different stores like postgres,redis, s3-like.
type StoresConfig struct {
	// PostgreSQL configuration
	PostgreSQL PostgreSQLConfig `json:"postgres" yaml:"postgres" validate:"required"`

	// Redis configuration
	Redis RedisConfig `json:"redis" yaml:"redis" validate:"required"`

	S3 S3Config `json:"s3" yaml:"s3" validate:"required"`
}

// S3Config holds the configuration for S3-like object storage.
type S3Config struct {
	// Endpoint is the S3-compatible storage endpoint.
	Endpoint string `json:"endpoint" yaml:"endpoint" validate:"required,url"`

	// AccessKeyID is the access key ID for the S3 bucket.
	AccessKeyID string `json:"accessKeyID" yaml:"accessKeyID" validate:"required"`

	// SecretAccessKey is the secret access key for the S3 bucket.
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey" validate:"required"`

	// Region is the region where the S3 bucket is located.
	Region string `json:"region" yaml:"region" validate:"required"`

	// UseSSL indicates whether to use SSL for S3 requests.
	// Set to true if using HTTPS
	UseSSL bool `json:"useSSL" yaml:"useSSL"`
}

// PostgreSQLConfig holds the configuration for connecting to a PostgreSQL database
type PostgreSQLConfig struct {
	// Data Source Name for the database connection, in pgx format
	//
	// For connection pooling arguments, take a look at https://pkg.go.dev/github.com/jackc/pgx/v5@v5.7.5/pgxpool#ParseConfig
	DSN string `json:"dsn" yaml:"dsn" validate:"required"`
}

// RedisConfig holds the configuration for connecting to a Redis database
type RedisConfig struct {
	// Address of the Redis server, e.g., "localhost:6379"
	Address string `json:"address" yaml:"address" validate:"required"`

	// Password for the Redis server, if any
	Password *string `json:"password,omitempty" yaml:"password,omitempty"`

	// Database index to use (default is 0)
	DB int `json:"db" yaml:"db" validate:"min=0"`
}

// This represents a temporary struct for configuration extraction from the config file.
type internalConfigStruct struct {
	// Debug mode for the application
	Debug bool `json:"debug" yaml:"debug"`
	// Admins is a list of credentials
	Admins         []AdminCredential   `json:"admins" yaml:"admins" validate:"required,min=1,dive,required"`
	Public         PublicConfig        `json:"public" yaml:"public" validate:"required"`
	NBRGLMBranding *bool               `json:"nbrglmBranding,omitempty" yaml:"nbrglmBranding,omitempty"` // NBRGLM branding flag
	Multitenancy   *bool               `json:"multitenancy" yaml:"multitenancy" validate:"required"`
	Server         ServerConfig        `json:"server" yaml:"server" validate:"required"`
	Observability  ObservabilityConfig `json:"observability" yaml:"observability" validate:"required"`
	Password       PasswordConfig      `json:"password" yaml:"password" validate:"required"`
	JWT            JWTConfig           `json:"jwt" yaml:"jwt" validate:"required"`
	Notifications  NotificationsConfig `json:"notifications" yaml:"notifications" validate:"required"`
	Branding       BrandingConfig      `json:"branding" yaml:"branding" validate:"required"`
	Security       SecurityConfig      `json:"security" yaml:"security" validate:"required"`
	Stores         StoresConfig        `json:"stores" yaml:"stores" validate:"required"`
}

// ConfigError represents an error that occurs during configuration initialization/reinitialization
type ConfigError struct {
	Message         string
	UnderlyingError error
}

func (c ConfigError) Error() string {
	if c.UnderlyingError != nil {
		return fmt.Sprintf("%v, UnderlyingError: %v", c.Message, c.UnderlyingError.Error())
	}
	return c.Message
}

func LoadConfigOptions(configFile string) error {
	var cfg *internalConfigStruct = new(internalConfigStruct)
	filePath, err := filepath.Abs(configFile)
	if err != nil {
		return ConfigError{
			Message: fmt.Sprintf("Unable to get convert given path (%s) to absolute path, Error: %v", configFile, err.Error()),
		}
	}

	file, err := os.ReadFile(filePath)

	if err != nil {
		return ConfigError{
			Message: fmt.Sprintf("Unable to read file %s, does the file exist and has correct permissions?", filePath),
		}
	}

	file = []byte(os.ExpandEnv(string(file)))

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		return ConfigError{Message: fmt.Sprintf("Unable to read config file. Does the config file exist at %s?", configFile), UnderlyingError: err}
	}

	if err := setDefaults(cfg); err != nil {
		return ConfigError{Message: "Invalid Base Configuration!", UnderlyingError: err}
	}

	// Set the debug mode value before validation
	opts.Debug = cfg.Debug

	if err := utils.Validator.Struct(cfg); err != nil {
		return ConfigError{Message: "Configuration validation failed", UnderlyingError: err}
	}

	// Assign the values to the global variables
	if cfg.NBRGLMBranding != nil {
		NBRGLMBranding = *cfg.NBRGLMBranding
	}
	if cfg.Multitenancy != nil {
		Multitenancy = *cfg.Multitenancy
	}
	Admins = cfg.Admins
	Server = &cfg.Server
	Public = &cfg.Public
	Observability = &cfg.Observability
	Password = &cfg.Password
	JWT = &cfg.JWT
	Notifications = &cfg.Notifications
	Branding = &cfg.Branding
	Security = &cfg.Security
	Stores = &cfg.Stores

	return nil
}

func setDefaults(cfg *internalConfigStruct) error {
	if strings.TrimSpace(cfg.Server.Host) == "" {
		cfg.Server.Host = "localhost"
	}
	if strings.TrimSpace(cfg.Server.Port) == "" {
		cfg.Server.Port = "3360" // Default port for the server
	}

	if strings.TrimSpace(cfg.Server.InstanceID) == "" {
		return ConfigError{Message: "Server Instance ID cannot be empty"}
	}

	if cfg.Server.TLSConfig != nil {
		if strings.TrimSpace(cfg.Server.TLSConfig.CertFile) == "" {
			cfg.Server.TLSConfig.CertFile = "/etc/nbrglm/auth-platform/cert.pem"
		}
		if strings.TrimSpace(cfg.Server.TLSConfig.KeyFile) == "" {
			cfg.Server.TLSConfig.KeyFile = "/etc/nbrglm/auth-platform/key.pem"
		}
	}

	// Otel configuration
	if strings.TrimSpace(cfg.Observability.LogLevel) == "" {
		if cfg.Debug {
			cfg.Observability.LogLevel = "debug" // Default to debug in debug mode
		} else {
			cfg.Observability.LogLevel = "info" // Default to info in production mode
		}
	}

	if strings.TrimSpace(cfg.Observability.OtelExporterProtocol) == "" {
		cfg.Observability.OtelExporterProtocol = "stdout" // Default to stdout if not set
	}
	if cfg.Observability.OtelExporterProtocol != "stdout" && strings.TrimSpace(cfg.Observability.OtelExporterEndpoint) == "" {
		return ConfigError{Message: "OpenTelemetry exporter endpoint cannot be empty when using grpc or http/protobuf protocols"}
	}
	if cfg.Observability.OtelExporterProtocol == "stdout" && strings.TrimSpace(cfg.Observability.OtelExporterEndpoint) != "" {
		return ConfigError{Message: "OpenTelemetry exporter endpoint should not be set when using stdout protocol"}
	}

	if strings.TrimSpace(string(cfg.Password.Algorithm)) == "" {
		cfg.Password.Algorithm = "bcrypt"
	}

	if cfg.Password.Bcrypt.Cost == 0 {
		cfg.Password.Bcrypt.Cost = 10
	}

	if cfg.Password.Argon2id.Memory == 0 {
		cfg.Password.Argon2id.Memory = 64 * 1024 // 64MB, in KiB
	}

	if cfg.Password.Argon2id.Iterations == 0 {
		cfg.Password.Argon2id.Iterations = 1 // Recommended default
	}

	if cfg.Password.Argon2id.Parallelism == 0 {
		cfg.Password.Argon2id.Parallelism = 1 // Recommended default
	}

	if cfg.Password.Argon2id.SaltLength == 0 {
		cfg.Password.Argon2id.SaltLength = 16 // Recommended default
	}

	if cfg.Password.Argon2id.KeyLength == 0 {
		cfg.Password.Argon2id.KeyLength = 32 // Recommended default
	}

	if strings.TrimSpace(cfg.Public.Scheme) == "" {
		if cfg.Debug {
			cfg.Public.Scheme = "http" // Default to http in debug mode
		} else {
			cfg.Public.Scheme = "https" // Default to https in production mode
		}
	}

	if strings.TrimSpace(cfg.Public.Domain) == "" {
		if !cfg.Debug {
			return ConfigError{Message: "Public domain cannot be empty"}
		}
		cfg.Public.Domain = "localhost" // Default domain for debug mode
	}
	if strings.TrimSpace(cfg.Public.SubDomain) == "" {
		if !cfg.Debug {
			return ConfigError{Message: "Public subdomain cannot be empty"}
		}
		cfg.Public.SubDomain = "auth" // Default subdomain for debug mode
	}

	if strings.TrimSpace(cfg.Public.DebugBaseURL) == "" {
		if cfg.Debug {
			scheme := "http"
			if cfg.Server.TLSConfig != nil {
				scheme = "https"
			}
			cfg.Public.DebugBaseURL = fmt.Sprintf("%s://%s:%s", scheme, cfg.Server.Host, cfg.Server.Port)
		}
		// No debug base URL in production mode
	}

	if cfg.JWT.SessionTokenExpiration == 0 {
		cfg.JWT.SessionTokenExpiration = 3600 // Default to 1 hour
	}

	if cfg.JWT.RefreshTokenExpiration == 0 {
		cfg.JWT.RefreshTokenExpiration = 2592000 // Default to 30 days
	}

	if len(cfg.JWT.Audiences) == 0 {
		cfg.JWT.Audiences = []string{cfg.Public.Domain, fmt.Sprintf("%s.%s", cfg.Public.SubDomain, cfg.Public.Domain)}
	} else {
		// Ensure the audiences contain the domain and subdomain
		if !slices.Contains(cfg.JWT.Audiences, cfg.Public.Domain) {
			cfg.JWT.Audiences = append(cfg.JWT.Audiences, cfg.Public.Domain)
		}
		if !slices.Contains(cfg.JWT.Audiences, fmt.Sprintf("%s.%s", cfg.Public.SubDomain, cfg.Public.Domain)) {
			cfg.JWT.Audiences = append(cfg.JWT.Audiences, fmt.Sprintf("%s.%s", cfg.Public.SubDomain, cfg.Public.Domain))
		}
	}

	if strings.TrimSpace(cfg.JWT.PrivateKeyFile) == "" {
		return ConfigError{Message: "RS256 Private Key File cannot be empty"}
	}
	if strings.TrimSpace(cfg.JWT.PublicKeyFile) == "" {
		return ConfigError{Message: "RS256 Public Key File cannot be empty"}
	}

	// No defaults for notifications configuration

	if strings.TrimSpace(cfg.Branding.AppName) == "" {
		return ConfigError{Message: "Branding AppName cannot be empty"}
	}

	if strings.TrimSpace(cfg.Branding.CompanyName) == "" {
		return ConfigError{Message: "Branding CompanyName cannot be empty"}
	}

	if strings.TrimSpace(cfg.Branding.CompanyNameShort) == "" {
		cfg.Branding.CompanyNameShort = cfg.Branding.AppName // Default to AppName if not provided
	}

	if strings.TrimSpace(cfg.Branding.SupportURL) == "" {
		return ConfigError{Message: "Branding SupportURL cannot be empty"}
	}

	if len(cfg.Security.CORS.AllowedMethods) == 0 {
		// Default to common methods if not specified
		cfg.Security.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	} else {
		// Add the common methods
		cfg.Security.CORS.AllowedMethods = append(cfg.Security.CORS.AllowedMethods, "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH")
	}

	if len(cfg.Security.CORS.AllowedHeaders) == 0 {
		// Set to default headers
		cfg.Security.CORS.AllowedHeaders = []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin", "User-Agent", "X-NAP-Refresh-Token", "X-CSRF-Token", "X-NAP-API-Key"}
	} else {
		// Append default headers
		cfg.Security.CORS.AllowedHeaders = append(cfg.Security.CORS.AllowedHeaders, "Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin", "User-Agent", "X-NAP-Refresh-Token", "X-CSRF-Token", "X-NAP-API-Key")
	}

	if slices.Contains(cfg.Security.CORS.AllowedOrigins, "*") {
		return ConfigError{Message: "Invalid value: AllowedOrigins contains invalid '*' value!"}
	}

	if slices.Contains(cfg.Security.CORS.AllowedMethods, "*") {
		return ConfigError{Message: "Invalid value: AllowedMethods contains invalid '*' value!"}
	}

	if slices.Contains(cfg.Security.CORS.AllowedHeaders, "*") {
		return ConfigError{Message: "Invalid value: AllowedHeaders contains invalid '*' value!"}
	}

	if cfg.Stores.PostgreSQL.DSN == "" {
		return ConfigError{Message: "PostgreSQL DSN cannot be empty"}
	}

	if strings.TrimSpace(cfg.Stores.Redis.Address) == "" {
		return ConfigError{Message: "Redis address cannot be empty"}
	}
	if cfg.Stores.Redis.DB < 0 {
		return ConfigError{Message: "Redis DB index cannot be negative"}
	}
	if cfg.Stores.Redis.Password != nil && strings.TrimSpace(*cfg.Stores.Redis.Password) == "" {
		return ConfigError{Message: "Redis password cannot be empty if provided"}
	}

	return nil
}
