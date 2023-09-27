package license

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/Clevermind-Think-Mint/buymint-cli-go/internal/logger"
	"github.com/Clevermind-Think-Mint/buymint-cli-go/internal/rest"
	"github.com/pkg/errors"
)

type License struct {
	Signature string                 `json:"serial"`
	Message   string                 `json:"message"`
	Meta      map[string]interface{} `json:"meta"`
	publicKey *rsa.PublicKey         `json:"-"`
}

func New(license string, options map[string]interface{}) (*License, error) {
	if options == nil {
		options = map[string]interface{}{}
	}
	if options["PublicKey"] == nil {
		options["PublicKey"] = "https://buy.bmint.studio/api/v1/service/microservice/licensor/key"
	}
	byteLicense, err := parseArgument(license, options)
	if err != nil {
		return nil, errors.Wrap(err, `Unable to parse license`)
	}
	bytePublicKey, err := parseArgument(options["PublicKey"].(string), options)
	if err != nil {
		return nil, errors.Wrap(err, `Unable to parse public key`)
	}
	signature, message, meta, err := extractLicenseData(byteLicense, bytePublicKey)
	if err != nil {
		return nil, errors.Wrap(err, `Unable to extract license data`)
	}
	// Converting PublicKey
	publicKey, err := convertPublicKey(bytePublicKey)
	if err != nil {
		return nil, errors.Wrap(err, `Unable to convert public key`)
	}
	return &License{
		Signature: signature,
		Message:   message,
		Meta:      meta,
		publicKey: publicKey,
	}, nil
}

// Validating if desired serial/metas are contained in current license
// See this simple explanation if you wish to understand what's happening: https://www.sohamkamani.com/golang/rsa-encryption/#signing-and-verification
func (t *License) Validate(meta map[string]interface{}) (bool, error) {
	logger.Debug("Verifying...\n\n\t.::Message::.\n\n%s\n\n\t.::Signature::.\n\n%s", t.Message, t.Signature)
	// Building verifiable message
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(t.Message))
	if err != nil {
		return false, errors.Wrap(err, `Unable to hash verifiable message`)
	}
	msgHashSum := msgHash.Sum(nil)
	// Since signature is encoded into base64 we decode it
	signature, err := base64.StdEncoding.DecodeString(t.Signature)
	if err != nil {
		return false, errors.Wrap(err, `Unable to decode base54 signature`)
	}
	// Verifying signature
	err = rsa.VerifyPKCS1v15(t.publicKey, crypto.SHA256, msgHashSum, signature)
	if err != nil {
		//return false, errors.New(`Unable to verify license`)
		return false, errors.Wrap(err, `Unable to verify license`)
	}
	if meta != nil {
		logger.Debug("Checking desired meta: %v against license meta: %v", meta, t.Meta)
		for key, value := range meta {
			logger.Debug("Checking %q: %v...", key, value)
			if t.Meta[key] != value {
				return false, errors.New(`Unable to verify the presence for metadata "` + key + `" into the license`)
			}
		}
	}
	logger.Info("License validated correctly")
	return true, nil
}

// Verifiyng license data and Extracting serial and meta
func extractLicenseData(license []byte, publicKey []byte) (string, string, map[string]interface{}, error) {
	logger.Debug("Extracting data from:\n\n\t.::License::.\n\n%s\n\n\t.::Public Key::.\n\n%s", license, publicKey)
	// Getting the signature from license
	reSignature := regexp.MustCompile(`(?s)====BEGIN SIGNATURE====(.*)====END SIGNATURE====`)
	signatureMatches := reSignature.FindStringSubmatch(string(license))
	if len(signatureMatches) != 2 {
		return "", "", nil, errors.New(`Invalid license signature format`)
	}
	signature := strings.Trim(signatureMatches[1], "\n")
	// Getting the message from license
	reMessage := regexp.MustCompile(`(?s)====BEGIN LICENSE====(.*)=====END LICENSE=====`)
	messageMatches := reMessage.FindStringSubmatch(string(license))
	if len(messageMatches) != 2 {
		return signature, "", nil, errors.New(`Invalid license message format`)
	}
	message := strings.Trim(messageMatches[0], "\n")
	// Getting the metadata from license
	reMetadata := regexp.MustCompile(`Metadata: (.*)`)
	metadataMatches := reMetadata.FindStringSubmatch(string(license))
	if len(messageMatches) != 2 {
		return signature, "", nil, errors.New(`Invalid license metadata format`)
	}
	metadataToParse := strings.Trim(metadataMatches[1], "\n")
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(metadataToParse), &metadata); err != nil {
		return signature, message, nil, errors.Wrap(err, "Unable to parse metadata from license")
	}
	// Offering extracted data from license
	return signature, message, metadata, nil
}

func parseArgument(arg string, options map[string]interface{}) ([]byte, error) {
	if isURL(arg) {
		_, content, _, err := rest.Get(arg, map[string]string{
			"Authorization": "Bearer " + options["Token"].(string),
		}, options)
		return content, err
	}
	if isPath(arg) {
		return os.ReadFile(arg)
	}
	// Otherwise arg is a string
	return []byte(arg), nil
}

// Checking if a string is an URL
func isURL(stringToCheck string) bool {
	_, err := url.ParseRequestURI(stringToCheck)
	if err != nil {
		return false
	}
	return true
}

// Checking if string is a valid path
func isPath(stringToCheck string) bool {
	_, err := os.Stat(stringToCheck)
	if err != nil {
		return false
	}
	return true
}

func convertPublicKey(key []byte) (*rsa.PublicKey, error) {
	data, _ := pem.Decode(key)
	if data == nil {
		return nil, errors.New(`Unable to decode RSA public key; ensure the URL/path/string is correct`)
	}
	var parsedKey *rsa.PublicKey
	keyInterface, _ := x509.ParsePKIXPublicKey(data.Bytes)
	parsedKey = keyInterface.(*rsa.PublicKey)
	return parsedKey, nil
}

func convertPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	data, _ := pem.Decode(key)
	if data == nil {
		return nil, errors.New(`Unable to decode RSA private key; ensure the URL/path/string is correct`)
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, `Unable to parse RSA private key by x509`)
	}
	return parsedKey, nil
}
