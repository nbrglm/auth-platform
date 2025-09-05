package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/nbrglm/nexeres/utils"
	"github.com/spf13/cobra"
)

type keygenConfig struct {
	// Mode, one of "cookie-signing", "rs256", or "csrf"
	Mode string `validate:"required,oneof=cookie-signing rs256 csrf"`

	// PrivateKeyPath is the path to the private key file for RS256 algorithm
	//
	// Defaults to "/etc/nbrglm/workspace/nexeres/keys/jwt/private.pem"
	// This field is required if Mode is "rs256".
	PrivateKeyPath string `validate:"file,required_if=Mode rs256"`

	// PublicKeyPath is the path to the public key file for RS256 algorithm
	// Defaults to "/etc/nbrglm/workspace/nexeres/keys/jwt/public.pem"
	// This field is required if Mode is "rs256".
	PublicKeyPath string `validate:"file,required_if=Mode rs256"`

	// CookieSigningSecret is the path to the secret key file for Cookie Signing algorithm
	// Defaults to "/etc/nbrglm/workspace/nexeres/keys/cookie/signing-secret"
	// This field is required if Mode is "cookie-signing".
	CookieSigningSecret string `validate:"file,required_if=Mode cookie-signing"`

	// CSRFSecretKeyPath is the path to the CSRF secret key file
	// Defaults to "/etc/nbrglm/workspace/nexeres/keys/csrf/secret"
	// This field is required if Mode is "csrf".
	CSRFSecretKeyPath string `validate:"file,required_if=Mode csrf"`

	// Force option, to overwrite existing keys, default: false
	Force bool
}

var keygenCfg = keygenConfig{}

func initKeygenCommand() {
	keygenCmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate keys for Password Hashing, Cookie Signing, and CSRF Protection.",
		Long:  "Generate keys for Password Hashing, Cookie Signing, and CSRF Protection. This command generates the necessary keys based on the specified mode.",
		Run: func(cmd *cobra.Command, args []string) {
			keygen(cmd)
		},
	}

	keygenCmd.Flags().StringVar(&keygenCfg.Mode, "mode", "csrf", "Mode to generate key for. One of 'cookie-signing', 'rs256', 'csrf'")
	keygenCmd.Flags().StringVar(&keygenCfg.PrivateKeyPath, "private-key-path", "/etc/nbrglm/workspace/nexeres/keys/jwt/private.pem", "Path to the private key file (for RS256 algorithm)")
	keygenCmd.Flags().StringVar(&keygenCfg.PublicKeyPath, "public-key-path", "/etc/nbrglm/workspace/nexeres/keys/jwt/public.pem", "Path to the public key file (for RS256 algorithm)")
	keygenCmd.Flags().StringVar(&keygenCfg.CookieSigningSecret, "cookie-signing-secret-path", "/etc/nbrglm/workspace/nexeres/keys/cookie/signing-secret", "Path to the secret key file (for Cookie Signing)")
	keygenCmd.Flags().StringVar(&keygenCfg.CSRFSecretKeyPath, "csrf-secret-key-path", "/etc/nbrglm/workspace/nexeres/keys/csrf/secret", "Path to the CSRF secret key file")
	keygenCmd.Flags().BoolVar(&keygenCfg.Force, "force", false, "Force overwrite existing keys")

	rootCmd.AddCommand(keygenCmd)
}

func keygen(cmd *cobra.Command) {
	// Initialize the validator before everything else, since validation is used by the config file loader.
	utils.InitValidator()
	if err := utils.Validator.Struct(keygenCfg); err != nil {
		cmd.PrintErrln("Validation error:", err)
		return
	}

	switch keygenCfg.Mode {
	case "cookie-signing":
		err := generateCookieSigningSecret(keygenCfg.CookieSigningSecret, keygenCfg.Force)
		if err != nil {
			cmd.PrintErrf("Error generating Cookie Signing Secret: %v\n", err)
			return
		}
	case "rs256":
		err := generateRS256Keys(keygenCfg.PrivateKeyPath, keygenCfg.PublicKeyPath, keygenCfg.Force)
		if err != nil {
			cmd.PrintErrf("Error generating RS256 keys: %v\n", err)
			return
		}
	case "csrf":
		err := generateCSRFSecretKey(keygenCfg.CSRFSecretKeyPath, keygenCfg.Force)
		if err != nil {
			cmd.PrintErrf("Error generating CSRF secret key: %v\n", err)
			return
		}
	default:
		cmd.PrintErrf("Unsupported mode: %s\n", keygenCfg.Mode)
		return
	}
	cmd.Println("Key generation completed successfully.")
	cmd.Println("You can now use the generated keys in the configuration file.")
	cmd.Println("For more information, refer to the documentation.")
	cmd.Println("Thank you for using the Nexeres CLI!")
	cmd.Println("If you have any questions or issues, please reach out to the support team.")
}

// generateCSRFSecretKey creates a new CSRF secret key at the specified path
func generateCSRFSecretKey(csrfSecretKeyPath string, force bool) error {
	// Check if the CSRF secret key already exists
	if utils.FileExists(csrfSecretKeyPath) && !force {
		return nil // CSRF secret key already exists, no need to generate
	}

	// Generate a new CSRF secret key
	csrfSecretKey := make([]byte, 32) // 32 bytes for CSRF secret key
	if _, err := rand.Read(csrfSecretKey); err != nil {
		return err
	}

	// os.Create creates the file if it does not exist, or truncates it to zero-length if it does.
	file, err := os.Create(csrfSecretKeyPath)
	if err != nil {
		return err
	}
	defer file.Close()

	b64CSRFSecretKey := utils.EncodeB64Key(csrfSecretKey)

	// Write the CSRF secret key to the file
	if _, err := file.WriteString(b64CSRFSecretKey); err != nil {
		return err
	}
	return nil
}

// generateCookieSigningSecret creates a new secret key for Cookie Signing at the specified path
func generateCookieSigningSecret(secretKeyPath string, force bool) error {
	// Check if the secret key already exists
	if utils.FileExists(secretKeyPath) && !force {
		return nil // Secret key already exists, no need to generate
	}

	// Generate a new secret key
	secretKey := make([]byte, 32) // 32 bytes for Cookie Signing Secret key
	if _, err := rand.Read(secretKey); err != nil {
		return err
	}

	// os.Create creates the file if it does not exist, or truncates it to zero-length if it does.
	file, err := os.Create(secretKeyPath)
	if err != nil {
		return err
	}
	defer file.Close()

	b64SecretKey := utils.EncodeB64Key(secretKey)

	// Write the secret key to the file
	if _, err := file.WriteString(b64SecretKey); err != nil {
		return err
	}
	return nil
}

// generateRS256Keys creates a new private key and its corresponding public key
// at the specified paths if they do not already exist.
func generateRS256Keys(privateKeyPath, publicKeyPath string, force bool) error {
	// Check if the private key already exists
	if utils.FileExists(privateKeyPath) && !force {
		return nil // Private key already exists, no need to generate
	}

	// We don't check for the public key existence here because it will be generated from the private key.

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return err
	}

	// os.Create creates the file if it does not exist, or truncates it to zero-length if it does.
	privFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privFile.Close()

	// Save the private key in PEM format
	pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	pubFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	// Write to file
	pem.Encode(pubFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return nil
}
