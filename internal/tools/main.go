package tools

import (
	"bytes"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func InterfaceSliceToStringSlice(data []interface{}) []string {
	result := make([]string, len(data))
	for i, v := range data {
		result[i] = v.(string)
	}
	return result
}

func StringSliceToInterfaceSlice(strings []string) []interface{} {
	s := make([]interface{}, len(strings))
	for i, v := range strings {
		s[i] = string(v)
	}
	return s
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// Copied from https://siongui.github.io/2018/03/09/go-match-common-element-in-two-array/
func Intersection(a, b []int) (c []int) {
	m := make(map[int]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func GetIntPointer(number int) *int {
	return &number
}

func GetStringPointer(s string) *string {
	return &s
}

func GetBoolPointer(b bool) *bool {
	return &b
}

func ConvertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	sort.Strings(s)

	return s
}

func FormatPrivilegeActionsForConfig[T any](actions []T) string {
	var stringActions []string
	for _, action := range actions {
		stringActions = append(stringActions, fmt.Sprintf("\"%v\"", action))
	}
	return strings.Join(stringActions, ", ")
}

func TestGetCertificateFingerprint(cert *x509.Certificate) (string, error) {
	fingerprintBytes := sha1.Sum(cert.Raw)
	var fingerprint bytes.Buffer
	for i, f := range fingerprintBytes {
		if i > 0 {
			fmt.Fprintf(&fingerprint, ":")
		}
		fmt.Fprintf(&fingerprint, "%02X", f)
	}
	return fingerprint.String(), nil
}

func TestGenerateRandomCertificate() (string, *x509.Certificate, string, string, error) {
	certPEM, privateCertPEM, _ := acctest.RandTLSCert("acctest")
	certBlock, _ := pem.Decode([]byte(certPEM))
	cert, _ := x509.ParseCertificate(certBlock.Bytes)
	certFingerprint, _ := TestGetCertificateFingerprint(cert)
	return certPEM, cert, certFingerprint, privateCertPEM, nil
}

func TestRetrieveCert(url string) (*x509.Certificate, error) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, clientErr := client.Get(url)
	if clientErr != nil {
		return nil, clientErr
	}
	if resp.TLS != nil {
		certificates := resp.TLS.PeerCertificates
		if len(certificates) > 0 {
			return certificates[0], nil
		}
	}

	return nil, fmt.Errorf("failed retrieving certificate from '%s'", url)
}

func TestPemEncode(b []byte, block string) (string, error) {
	var buf bytes.Buffer
	pb := &pem.Block{Type: block, Bytes: b}
	if err := pem.Encode(&buf, pb); err != nil {
		return "", err
	}

	return buf.String(), nil
}
