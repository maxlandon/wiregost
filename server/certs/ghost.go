package certs

const (
	// GhostCA - Directory containing sliver certificates
	GhostCA = "ghost"
)

// GhostGenerateECCCertificate - Generate a certificate signed with a given CA
func GhostGenerateECCCertificate(ghostName string) ([]byte, []byte, error) {
	cert, key := GenerateECCCertificate(GhostCA, ghostName, false, true)
	err := SaveCertificate(GhostCA, ECCKey, ghostName, cert, key)
	return cert, key, err
}

// GhostGenerateRSACertificate - Generate a certificate signed with a given CA
func GhostGenerateRSACertificate(ghostName string) ([]byte, []byte, error) {
	cert, key := GenerateRSACertificate(GhostCA, ghostName, false, true)
	err := SaveCertificate(GhostCA, RSAKey, ghostName, cert, key)
	return cert, key, err
}
