package cmd

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	license "github.com/Clevermind-Think-Mint/buymint-cli-go/pkg"
)

var validateLicenseCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a license against a specific serial/metas",
	RunE:  validateLicense,
}

func validateLicense(cmd *cobra.Command, args []string) error {
	// Parsing meta JSON string
	metaToParse := viper.GetString("meta")
	var meta map[string]interface{}
	if err := json.Unmarshal([]byte(metaToParse), &meta); err != nil {
		return errors.Wrap(err, "Unable to parse meta from CLI argument")
	}
	// Building new license
	license, err := license.New(viper.GetString("license"), map[string]interface{}{
		"PublicKey":         viper.GetString("public_key"),
		"IgnoreInsecureSsl": viper.GetBool("self-signed"),
		"Token":             viper.GetString("token"),
	})
	if err != nil {
		return errors.Wrap(err, "Unable to initialize License")
	}
	// Validating license against desired data
	_, err = license.Validate(meta)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	validateLicenseCmd.Flags().StringP("license", "l", "", "The license")
	viper.BindPFlag("license", validateLicenseCmd.Flags().Lookup("license"))
	validateLicenseCmd.Flags().StringP("meta", "m", "{}", "The meta data to validate, written in JSON format (Eg: {\"foo\":\"test\"})")
	viper.BindPFlag("meta", validateLicenseCmd.Flags().Lookup("meta"))
	validateLicenseCmd.Flags().StringP("public_key", "p", "", "The public key to use to validate the license")
	viper.BindPFlag("public_key", validateLicenseCmd.Flags().Lookup("public_key"))
	rootCmd.AddCommand(validateLicenseCmd)
}
