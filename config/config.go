// Package config handles the configuration for the Nexeres authentication server.
//
// Note: The Config structs have two tags for json and yaml.
// JSON tags: Used to mark what data is visible to admins via the config json endpoint.
// YAML tags: Used to actually configure this auth server.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/nbrglm/nexeres/opts"
	"github.com/nbrglm/nexeres/utils"
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
	Multitenancy  = false // Default to false, can be set to true in the config file
	Notifications *NotificationsConfig
	Public        *PublicConfig
	Server        *ServerConfig
	Observability *ObservabilityConfig
	Password      *PasswordConfig
	JWT           *JWTConfig
	Branding      *BrandingConfig
	Security      *SecurityConfig
	Stores        *StoresConfig

	// Admins is a list of credentials for admin users
	Admins AdminConfig

	Config *CompleteConfig
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
	// The log level for nexeres.
	LogLevel string `json:"logLevel" yaml:"logLevel" validate:"required,oneof=debug info warn error dpanic panic fatal"`

	// OtelExporterEndpoint is the endpoint for the OpenTelemetry exporter
	OtelExporterEndpoint string `json:"otelExporterEndpoint" yaml:"otelExporterEndpoint" validate:"omitempty,url"`

	// OtelExporterProtocol is the protocol of the OpenTelemetry exporter
	OtelExporterProtocol string `json:"otelExporterProtocol" yaml:"otelExporterProtocol" validate:"required,oneof=http/protobuf grpc stdout"`
}

// PublicConfig holds the public server configuration options
//
// Note: If hosting Nexeres at "auth.example.com", then set the values as follows:
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

	// Domain is the domain at which Nexeres UI is being hosted.
	//
	// Eg. if hosting Nexeres at auth.example.com, provide "example.com" here.
	//
	// This setting is used as the Domain for all Cookies (except Refresh Token Cookie, that is only set at "subdomain.domain")
	Domain string `json:"domain" yaml:"domain" validate:"required"`

	// Subdomain at which Nexeres UI is being hosted.
	//
	// Eg. If hosting Nexeres at auth.example.com, provide "auth" here.
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

type AdminConfig struct {
	// Emails is a list of admin emails.
	//
	// For login, an email will be sent to the user's address if it exists in this list.
	// The user can then use the code to login.
	Emails []string `json:"-" yaml:"emails" validate:"required,dive,email"`

	// The session, without any calls to an admin api, will expire in this duration.
	//
	// In seconds, default 15 minutes
	SessionTimeout int `json:"-" yaml:"sessionTimeout" validate:"required,min=300"`
}

// TLSConfig holds the TLS configuration for the server
type TLSConfig struct {
	// CertFile is the path to the TLS certificate file
	CertFile string `json:"-" yaml:"certFile" validate:"required,file"`

	// KeyFile is the path to the TLS private key file
	KeyFile string `json:"-" yaml:"keyFile" validate:"required,file"`

	// CAFile is the path to the TLS CA certificate file
	// If provided, the server will use this CA to verify client certificates
	CAFile *string `json:"-" yaml:"caFile,omitempty" validate:"omitempty,file"`
}

// PasswordConfig holds the configuration for password policies
type PasswordConfig struct {
	// The algorithm used for hashing passwords, default: bcrypt
	Algorithm PasswordHashingAlgorithm `json:"-" yaml:"algorithm" validate:"required,oneof=bcrypt argon2id"`

	// Configuration for Bcrypt password hashing
	Bcrypt BcryptConfig `json:"-" yaml:"bcrypt" validate:"required_if=Algorithm bcrypt"`

	// Configuration for Argon2id password hashing
	Argon2id Argon2idConfig `json:"-" yaml:"argon2id" validate:"required_if=Algorithm argon2id"`
}

// BcryptConfig holds the configuration for Bcrypt password hashing
type BcryptConfig struct {
	// The cost factor for Bcrypt hashing (default is 10)
	Cost int `json:"-" yaml:"cost" validate:"required,min=10,max=31"`
}

// Argon2idConfig holds the configuration for Argon2id password hashing
type Argon2idConfig struct {
	// Memory in KiB (16MB minimum, and default)
	Memory int `json:"-" yaml:"memory" validate:"required,min=16384"`

	// Number of iterations (3-4 recommended, default 3)
	Iterations int `json:"-" yaml:"iterations" validate:"required,min=1"`

	// Number of parallel threads (1-4 recommended, default 1)
	Parallelism int `json:"-" yaml:"parallelism" validate:"required,min=1"`

	// Salt length in bytes (16-32 recommended, default 16)
	SaltLength int `json:"-" yaml:"saltLength" validate:"required,min=16,max=32"`

	// Key length in bytes (32-64 recommended, default 32)
	KeyLength int `json:"-" yaml:"keyLength" validate:"required,min=32,max=64"`
}

// JWTConfig holds the configuration for JWT tokens
type JWTConfig struct {
	// Session token expiration time in seconds (default: 1hr, 3600)
	SessionTokenExpiration int `json:"sessionTokenExpiration" yaml:"sessionTokenExpiration" validate:"required,min=60"`

	// Refresh token expiration time in seconds (default: 30d, 2592000)
	RefreshTokenExpiration int `json:"refreshTokenExpiration" yaml:"refreshTokenExpiration" validate:"required,min=86400"`

	// Path to the private key file for RS256 algorithm
	PrivateKeyFile string `json:"-" yaml:"privateKeyFile" validate:"required,file"`

	// Path to the public key file for RS256 algorithm
	PublicKeyFile string `json:"-" yaml:"publicKeyFile" validate:"required,file"`

	// Audiences claim for the JWT.
	//
	// NOTE: This is NOT for OIDC. This is for the session tokens! OIDC configuration is stored in the DB per tenant.
	//
	// NOTE: Do not add the subdomain or domain you added in the Public Config here, it is automatically added and will be duplicated if you add.
	// Eg. If domain has value example.com and subdomain has value auth,
	// then auth.example.com, and example.com are automatically added as audience claims.
	Audiences []string `json:"-" yaml:"audiences" validate:"required,dive,required"`
}

type NotificationsConfig struct {
	// Email configuration for sending notifications
	Email EmailNotificationConfig `json:"email" yaml:"email" validate:"required"`

	// SMS Configuration for sending notifications
	SMS *SMSNotificationConfig `json:"sms,omitempty" yaml:"sms,omitempty" validate:"omitempty"`
}

// EmailNotificationConfig holds the configuration for email notifications.
type EmailNotificationConfig struct {
	// Provider is the email provider to use for sending emails.
	Provider string `json:"provider" yaml:"provider" validate:"required,oneof=ses sendgrid smtp"`

	// SendGridConfig holds the configuration for SendGrid email provider.
	SendGrid *SendGridProviderConfig `json:"-" yaml:"sendgrid,omitempty" validate:"omitempty,required_if=Provider sendgrid"`

	// SMTPConfig holds the configuration for SMTP email provider.
	SMTP *SMTPProviderConfig `json:"-" yaml:"smtp,omitempty" validate:"omitempty,required_if=Provider smtp"`

	// SESConfig holds the configuration for AWS SES email provider.
	SES *SESProviderConfig `json:"-" yaml:"ses,omitempty" validate:"omitempty,required_if=Provider ses"`

	// Endpoints holds the configuration for URLs inside emails.
	Endpoints EmailEndpointsConfig `json:"endpoints" yaml:"endpoints" validate:"required"`

	// TemplatesDir is the directory where email templates are stored.
	TemplatesDir *string `json:"-" yaml:"templatesDir,omitempty" validate:"omitempty,file"`
}

// EmailEndpointsConfig holds the configuration for URLs inside emails.
type EmailEndpointsConfig struct {
	// VerificationEmail is the endpoint for the email verification link.
	// A `token` parameter will be passed to this URL.
	// Pass a full url, eg. https://auth.example.com/verify-email
	VerificationEmail string `json:"verificationEmail" yaml:"verificationEmail" validate:"required,url"`

	// PasswordReset is the endpoint for the password reset link.
	// A `token` parameter will be passed to this URL.
	// Pass a full url, eg. https://auth.example.com/password-reset
	PasswordReset string `json:"passwordReset" yaml:"passwordReset" validate:"required,url"`
}

// SendGridProviderConfig holds the configuration for SendGrid email provider.
type SendGridProviderConfig struct {
	// API key for SendGrid
	APIKey string `json:"-" yaml:"apiKey" validate:"required"`

	// Email address from which notifications are sent
	FromAddress string `json:"-" yaml:"fromAddress" validate:"required,email"`

	// Name associated with the FromAddress, optional
	FromName *string `json:"-" yaml:"fromName,omitempty"`
}

// SMTPProviderConfig holds the configuration for SMTP email provider.
//
// Note: SMTP is a generic email provider that can be used with any SMTP server.
// IT DOES NOT SUPPORT RETRIES, OR ANY ADVANCED FEATURES.
// It is recommended to use a more robust provider like SendGrid or AWS SES for production use.
type SMTPProviderConfig struct {
	// SMTP server host
	Host string `json:"-" yaml:"host" validate:"required"`

	// SMTP server port
	Port string `json:"-" yaml:"port"`

	// Email address from which notifications are sent
	FromAddress string `json:"-" yaml:"fromAddress" validate:"required,email"`

	// Password for SMTP authentication
	Password string `json:"-" yaml:"password" validate:"required"`
}

// SESProviderConfig holds the configuration for AWS SES email provider.
type SESProviderConfig struct {
	// AWS region where SES is configured
	Region string `json:"-" yaml:"region" validate:"required"`

	// AWS Access Key ID for SES
	AccessKeyID string `json:"-" yaml:"accessKeyID" validate:"required"`

	// AWS Secret Access Key for SES
	SecretAccessKey string `json:"-" yaml:"secretAccessKey" validate:"required"`

	// Email address from which notifications are sent
	FromAddress string `json:"-" yaml:"fromAddress" validate:"required,email"`

	// Optional: Name associated with the FromAddress, used in email headers, if not provided, AppName is used.
	FromName *string `json:"-" yaml:"fromName,omitempty"`
}

// SMSNotificationConfig holds the configuration for SMS notifications.
type SMSNotificationConfig struct {
	// TODO: Implement SMSNotificationConfig
	Provider string `json:"provider" yaml:"provider" validate:"required,oneof=twilio"`

	// TemplatesDir is the directory where SMS templates are stored.
	TemplatesDir *string `json:"-" yaml:"templatesDir,omitempty" validate:"omitempty,file"`
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
	AuditLogs AuditLogsConfig `json:"auditLogs" yaml:"auditLogs" validate:"required"`

	// The list of API Keys which are allowed to access the API endpoints.
	// Requests without an API key, or with a key not specified here, will be denied with 401.
	APIKeys []APIKeyConfig `json:"apiKeys" yaml:"apiKeys" validate:"required,dive"`

	// CORS configuration for the application.
	CORS CORSConfig `json:"-" yaml:"cors" validate:"required"`

	// Rate limiting configuration.
	RateLimit RateLimitConfig `json:"rateLimit" yaml:"rateLimit" validate:"required"`
}

type AuditLogsConfig struct {
	// Enable or disable audit logs.
	Enable bool `json:"enable" yaml:"enable"`
}

type APIKeyConfig struct {
	// The name of the API Key
	Name string `json:"name" yaml:"name" validate:"required"`

	// Description of the key
	Description string `json:"description" yaml:"description" validate:"required"`

	// The key itself. No restriction on the length, but keep it sensible please.
	Key string `json:"-" yaml:"key" validate:"required"`
}

// CORSConfig holds the configuration for Cross-Origin Resource Sharing (CORS).
type CORSConfig struct {
	// AllowedOrigins is a list of origins that are allowed to access the resources.
	AllowedOrigins []string `json:"-" yaml:"allowedOrigins" validate:"dive,url"`

	// AllowedMethods is a list of HTTP methods that are allowed.
	AllowedMethods []string `json:"-" yaml:"allowedMethods" validate:"required,dive,oneof=GET POST PUT DELETE OPTIONS PATCH HEAD"`

	// AllowedHeaders is a list of headers that are allowed in requests.
	AllowedHeaders []string `json:"-" yaml:"allowedHeaders" validate:"required"`
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
	PostgreSQL PostgreSQLConfig `json:"-" yaml:"postgres" validate:"required"`

	// Redis configuration
	Redis RedisConfig `json:"-" yaml:"redis" validate:"required"`

	S3 S3Config `json:"-" yaml:"s3" validate:"required"`
}

// S3Config holds the configuration for S3-like object storage.
type S3Config struct {
	// Endpoint is the S3-compatible storage endpoint.
	Endpoint string `json:"-" yaml:"endpoint" validate:"required,url"`

	// AccessKeyID is the access key ID for the S3 bucket.
	AccessKeyID string `json:"-" yaml:"accessKeyID" validate:"required"`

	// SecretAccessKey is the secret access key for the S3 bucket.
	SecretAccessKey string `json:"-" yaml:"secretAccessKey" validate:"required"`

	// Region is the region where the S3 bucket is located.
	Region string `json:"-" yaml:"region" validate:"required"`

	// UseSSL indicates whether to use SSL for S3 requests.
	// Set to true if using HTTPS
	UseSSL bool `json:"-" yaml:"useSSL"`
}

// PostgreSQLConfig holds the configuration for connecting to a PostgreSQL database
type PostgreSQLConfig struct {
	// Data Source Name for the database connection, in pgx format
	//
	// For connection pooling arguments, take a look at https://pkg.go.dev/github.com/jackc/pgx/v5@v5.7.5/pgxpool#ParseConfig
	DSN string `json:"-" yaml:"dsn" validate:"required"`
}

// RedisConfig holds the configuration for connecting to a Redis database
type RedisConfig struct {
	// Address of the Redis server, e.g., "localhost:6379"
	Address string `json:"-" yaml:"address" validate:"required"`

	// Password for the Redis server, if any
	Password *string `json:"-" yaml:"password,omitempty"`

	// Database index to use (default is 0)
	DB int `json:"-" yaml:"db" validate:"min=0"`
}

// This represents a temporary struct for configuration extraction from the config file.
type CompleteConfig struct {
	// Debug mode for the application
	Debug bool `json:"debug" yaml:"debug"`
	// Admins is a list of credentials
	Admins        AdminConfig         `json:"-" yaml:"admins" validate:"required"`
	Public        PublicConfig        `json:"public" yaml:"public" validate:"required"`
	Multitenancy  *bool               `json:"multitenancy" yaml:"multitenancy" validate:"required"`
	Server        ServerConfig        `json:"-" yaml:"server" validate:"required"`
	Observability ObservabilityConfig `json:"-" yaml:"observability" validate:"required"`
	Password      PasswordConfig      `json:"-" yaml:"password" validate:"required"`
	JWT           JWTConfig           `json:"jwt" yaml:"jwt" validate:"required"`
	Notifications NotificationsConfig `json:"notifications" yaml:"notifications" validate:"required"`
	Branding      BrandingConfig      `json:"branding" yaml:"branding" validate:"required"`
	Security      SecurityConfig      `json:"security" yaml:"security" validate:"required"`
	Stores        StoresConfig        `json:"-" yaml:"stores" validate:"required"`
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
	Config = new(CompleteConfig)
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

	err = yaml.Unmarshal(file, Config)
	if err != nil {
		return ConfigError{Message: fmt.Sprintf("Unable to read config file. Does the config file exist at %s?", configFile), UnderlyingError: err}
	}

	if err := setDefaults(); err != nil {
		return ConfigError{Message: "Invalid Base Configuration!", UnderlyingError: err}
	}

	// Set the debug mode value before validation
	opts.Debug = Config.Debug

	if err := utils.Validator.Struct(Config); err != nil {
		return ConfigError{Message: "Configuration validation failed", UnderlyingError: err}
	}

	// Assign the values to the global variables
	if Config.Multitenancy != nil {
		Multitenancy = *Config.Multitenancy
	}
	Admins = Config.Admins
	Server = &Config.Server
	Public = &Config.Public
	Observability = &Config.Observability
	Password = &Config.Password
	JWT = &Config.JWT
	Notifications = &Config.Notifications
	Branding = &Config.Branding
	Security = &Config.Security
	Stores = &Config.Stores

	return nil
}

func setDefaults() error {
	if strings.TrimSpace(Config.Server.Host) == "" {
		Config.Server.Host = "localhost"
	}
	if strings.TrimSpace(Config.Server.Port) == "" {
		Config.Server.Port = "3360" // Default port for the server
	}

	if strings.TrimSpace(Config.Server.InstanceID) == "" {
		return ConfigError{Message: "Server Instance ID cannot be empty"}
	}

	if Config.Server.TLSConfig != nil {
		if strings.TrimSpace(Config.Server.TLSConfig.CertFile) == "" {
			Config.Server.TLSConfig.CertFile = "/etc/nbrglm/workspace/nexeres/cert.pem"
		}
		if strings.TrimSpace(Config.Server.TLSConfig.KeyFile) == "" {
			Config.Server.TLSConfig.KeyFile = "/etc/nbrglm/workspace/nexeres/key.pem"
		}
	}

	if len(Config.Admins.Emails) == 0 {
		return ConfigError{Message: "At least one admin email must be provided, otherwise you may not be able to change ANY admin settings!"}
	}

	if Config.Admins.SessionTimeout == 0 || Config.Admins.SessionTimeout < 300 {
		Config.Admins.SessionTimeout = 900 // Default to 15 minutes
	}

	// Otel configuration
	if strings.TrimSpace(Config.Observability.LogLevel) == "" {
		if Config.Debug {
			Config.Observability.LogLevel = "debug" // Default to debug in debug mode
		} else {
			Config.Observability.LogLevel = "info" // Default to info in production mode
		}
	}

	if strings.TrimSpace(Config.Observability.OtelExporterProtocol) == "" {
		Config.Observability.OtelExporterProtocol = "stdout" // Default to stdout if not set
	}
	if Config.Observability.OtelExporterProtocol != "stdout" && strings.TrimSpace(Config.Observability.OtelExporterEndpoint) == "" {
		return ConfigError{Message: "OpenTelemetry exporter endpoint cannot be empty when using grpc or http/protobuf protocols"}
	}
	if Config.Observability.OtelExporterProtocol == "stdout" && strings.TrimSpace(Config.Observability.OtelExporterEndpoint) != "" {
		return ConfigError{Message: "OpenTelemetry exporter endpoint should not be set when using stdout protocol"}
	}

	if strings.TrimSpace(string(Config.Password.Algorithm)) == "" {
		Config.Password.Algorithm = "bcrypt"
	}

	if Config.Password.Bcrypt.Cost == 0 {
		Config.Password.Bcrypt.Cost = 10
	}

	if Config.Password.Argon2id.Memory == 0 {
		Config.Password.Argon2id.Memory = 64 * 1024 // 64MB, in KiB
	}

	if Config.Password.Argon2id.Iterations == 0 {
		Config.Password.Argon2id.Iterations = 1 // Recommended default
	}

	if Config.Password.Argon2id.Parallelism == 0 {
		Config.Password.Argon2id.Parallelism = 1 // Recommended default
	}

	if Config.Password.Argon2id.SaltLength == 0 {
		Config.Password.Argon2id.SaltLength = 16 // Recommended default
	}

	if Config.Password.Argon2id.KeyLength == 0 {
		Config.Password.Argon2id.KeyLength = 32 // Recommended default
	}

	if strings.TrimSpace(Config.Public.Scheme) == "" {
		if Config.Debug {
			Config.Public.Scheme = "http" // Default to http in debug mode
		} else {
			Config.Public.Scheme = "https" // Default to https in production mode
		}
	}

	if strings.TrimSpace(Config.Public.Domain) == "" {
		if !Config.Debug {
			return ConfigError{Message: "Public domain cannot be empty"}
		}
		Config.Public.Domain = "localhost" // Default domain for debug mode
	}
	if strings.TrimSpace(Config.Public.SubDomain) == "" {
		if !Config.Debug {
			return ConfigError{Message: "Public subdomain cannot be empty"}
		}
		Config.Public.SubDomain = "auth" // Default subdomain for debug mode
	}

	if strings.TrimSpace(Config.Public.DebugBaseURL) == "" {
		if Config.Debug {
			scheme := "http"
			if Config.Server.TLSConfig != nil {
				scheme = "https"
			}
			Config.Public.DebugBaseURL = fmt.Sprintf("%s://%s:%s", scheme, Config.Server.Host, Config.Server.Port)
		}
		// No debug base URL in production mode
	}

	if Config.JWT.SessionTokenExpiration == 0 {
		Config.JWT.SessionTokenExpiration = 3600 // Default to 1 hour
	}

	if Config.JWT.RefreshTokenExpiration == 0 {
		Config.JWT.RefreshTokenExpiration = 2592000 // Default to 30 days
	}

	if len(Config.JWT.Audiences) == 0 {
		Config.JWT.Audiences = []string{Config.Public.Domain, fmt.Sprintf("%s.%s", Config.Public.SubDomain, Config.Public.Domain)}
	} else {
		// Ensure the audiences contain the domain and subdomain
		if !slices.Contains(Config.JWT.Audiences, Config.Public.Domain) {
			Config.JWT.Audiences = append(Config.JWT.Audiences, Config.Public.Domain)
		}
		if !slices.Contains(Config.JWT.Audiences, fmt.Sprintf("%s.%s", Config.Public.SubDomain, Config.Public.Domain)) {
			Config.JWT.Audiences = append(Config.JWT.Audiences, fmt.Sprintf("%s.%s", Config.Public.SubDomain, Config.Public.Domain))
		}
	}

	if strings.TrimSpace(Config.JWT.PrivateKeyFile) == "" {
		return ConfigError{Message: "RS256 Private Key File cannot be empty"}
	}
	if strings.TrimSpace(Config.JWT.PublicKeyFile) == "" {
		return ConfigError{Message: "RS256 Public Key File cannot be empty"}
	}

	// No defaults for notifications configuration

	if strings.TrimSpace(Config.Branding.AppName) == "" {
		return ConfigError{Message: "Branding AppName cannot be empty"}
	}

	if strings.TrimSpace(Config.Branding.CompanyName) == "" {
		return ConfigError{Message: "Branding CompanyName cannot be empty"}
	}

	if strings.TrimSpace(Config.Branding.CompanyNameShort) == "" {
		Config.Branding.CompanyNameShort = Config.Branding.AppName // Default to AppName if not provided
	}

	if strings.TrimSpace(Config.Branding.SupportURL) == "" {
		return ConfigError{Message: "Branding SupportURL cannot be empty"}
	}

	if len(Config.Security.CORS.AllowedMethods) == 0 {
		// Default to common methods if not specified
		Config.Security.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	} else {
		// Add the common methods
		Config.Security.CORS.AllowedMethods = append(Config.Security.CORS.AllowedMethods, "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH")
	}

	if len(Config.Security.CORS.AllowedHeaders) == 0 {
		// Set to default headers
		Config.Security.CORS.AllowedHeaders = []string{"Content-Type", "X-Requested-With", "Accept", "Origin", "User-Agent", "X-NEXERES-Refresh-Token", "X-NEXERES-API-Key", "X-NEXERES-Session-Token"}
	} else {
		// Append default headers
		Config.Security.CORS.AllowedHeaders = append(Config.Security.CORS.AllowedHeaders, "Content-Type", "X-Requested-With", "Accept", "Origin", "User-Agent", "X-NEXERES-Refresh-Token", "X-NEXERES-API-Key", "X-NEXERES-Session-Token")
	}

	if slices.Contains(Config.Security.CORS.AllowedOrigins, "*") {
		return ConfigError{Message: "Invalid value: AllowedOrigins contains invalid '*' value!"}
	}

	if slices.Contains(Config.Security.CORS.AllowedMethods, "*") {
		return ConfigError{Message: "Invalid value: AllowedMethods contains invalid '*' value!"}
	}

	if slices.Contains(Config.Security.CORS.AllowedHeaders, "*") {
		return ConfigError{Message: "Invalid value: AllowedHeaders contains invalid '*' value!"}
	}

	if Config.Stores.PostgreSQL.DSN == "" {
		return ConfigError{Message: "PostgreSQL DSN cannot be empty"}
	}

	if strings.TrimSpace(Config.Stores.Redis.Address) == "" {
		return ConfigError{Message: "Redis address cannot be empty"}
	}
	if Config.Stores.Redis.DB < 0 {
		return ConfigError{Message: "Redis DB index cannot be negative"}
	}
	if Config.Stores.Redis.Password != nil && strings.TrimSpace(*Config.Stores.Redis.Password) == "" {
		return ConfigError{Message: "Redis password cannot be empty if provided"}
	}

	return nil
}
