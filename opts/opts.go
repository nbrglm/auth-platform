package opts

// ConfigPath is the path to the configuration file.
var ConfigPath *string = new(string)

var Debug bool

// Default org data
const DefaultOrgSlug = "default"
const DefaultOrgName = "Default Organization"
const DefaultOrgId = "019735ab-f216-717f-9d12-e3915453c8d0"

// S3 Store Bucket Names
const S3StoreBucketName string = "nexeres"

// Admin OTP Config
const AdminOTPLength = 8

// Used for configuring everything, from metrics to logging.
// This file contains the version information for the application.
// This file is not meant to be modified.
const Name = "nexeres"
const FullName = "github.com/nbrglm/nexeres"
const Version = "0.0.1"
const VersionName = "v0.0.1"
const VersionDate = "2025-05-23"
const VersionDescription = "Initial release of the NBRGLM Nexeres CLI & API application."
