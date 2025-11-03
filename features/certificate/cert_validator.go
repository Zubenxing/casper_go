package certificate

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

// ValidateCertRequest 证书验证请求
type ValidateCertRequest struct {
	CertContent       string `json:"cert_content" binding:"required"` // PEM 格式证书内容
	PrivateKeyContent string `json:"private_key_content"`             // PEM 格式私钥内容（可选）
}

// ValidateCertResponse 证书验证响应
type ValidateCertResponse struct {
	Valid              bool      `json:"valid"`
	Version            int       `json:"version"`
	SerialNumber       string    `json:"serial_number"`
	Issuer             string    `json:"issuer"`
	Subject            string    `json:"subject"`
	CommonName         string    `json:"common_name"`
	Organization       string    `json:"organization"`
	Country            string    `json:"country"`
	DNSNames           []string  `json:"dns_names"`
	IPAddresses        []string  `json:"ip_addresses"`
	EmailAddresses     []string  `json:"email_addresses"`
	NotBefore          time.Time `json:"not_before"`
	NotAfter           time.Time `json:"not_after"`
	DaysLeft           int       `json:"days_left"`
	IsExpired          bool      `json:"is_expired"`
	IsNotYetValid      bool      `json:"is_not_yet_valid"`
	SignatureAlgorithm string    `json:"signature_algorithm"`
	PublicKeyAlgorithm string    `json:"public_key_algorithm"`
	PublicKeySize      int       `json:"public_key_size"`
	KeyUsage           []string  `json:"key_usage"`
	ExtKeyUsage        []string  `json:"ext_key_usage"`
	IsCA               bool      `json:"is_ca"`
	KeyPairMatched     bool      `json:"key_pair_matched"` // 证书和私钥是否配对
	KeyPairChecked     bool      `json:"key_pair_checked"` // 是否进行了配对检查
	ErrorMessage       string    `json:"error_message,omitempty"`
}

// ValidateCertificate 验证证书文件（可选私钥配对验证）
func ValidateCertificate(certContent, privateKeyContent string) (*ValidateCertResponse, error) {
	// 1. 解码 PEM
	block, _ := pem.Decode([]byte(certContent))
	if block == nil {
		return &ValidateCertResponse{
			Valid:        false,
			ErrorMessage: "无效的 PEM 格式",
		}, nil
	}

	if block.Type != "CERTIFICATE" {
		return &ValidateCertResponse{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("无效的 PEM 类型: %s (期望 CERTIFICATE)", block.Type),
		}, nil
	}

	// 2. 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return &ValidateCertResponse{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("解析证书失败: %v", err),
		}, nil
	}

	// 3. 计算剩余天数
	now := time.Now()
	daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
	isExpired := now.After(cert.NotAfter)
	isNotYetValid := now.Before(cert.NotBefore)

	// 4. 提取主题信息
	country := ""
	if len(cert.Subject.Country) > 0 {
		country = cert.Subject.Country[0]
	}

	organization := ""
	if len(cert.Subject.Organization) > 0 {
		organization = cert.Subject.Organization[0]
	}

	// 5. 获取公钥算法和大小
	pubKeyAlgo := ""
	pubKeySize := 0

	switch pubKey := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		pubKeyAlgo = "RSA"
		pubKeySize = pubKey.N.BitLen()
	case *ecdsa.PublicKey:
		pubKeyAlgo = "ECDSA"
		pubKeySize = pubKey.Curve.Params().BitSize
	default:
		pubKeyAlgo = "Unknown"
	}

	// 6. 解析 KeyUsage
	keyUsage := []string{}
	if cert.KeyUsage&x509.KeyUsageDigitalSignature != 0 {
		keyUsage = append(keyUsage, "Digital Signature")
	}
	if cert.KeyUsage&x509.KeyUsageContentCommitment != 0 {
		keyUsage = append(keyUsage, "Content Commitment")
	}
	if cert.KeyUsage&x509.KeyUsageKeyEncipherment != 0 {
		keyUsage = append(keyUsage, "Key Encipherment")
	}
	if cert.KeyUsage&x509.KeyUsageDataEncipherment != 0 {
		keyUsage = append(keyUsage, "Data Encipherment")
	}
	if cert.KeyUsage&x509.KeyUsageKeyAgreement != 0 {
		keyUsage = append(keyUsage, "Key Agreement")
	}
	if cert.KeyUsage&x509.KeyUsageCertSign != 0 {
		keyUsage = append(keyUsage, "Certificate Sign")
	}
	if cert.KeyUsage&x509.KeyUsageCRLSign != 0 {
		keyUsage = append(keyUsage, "CRL Sign")
	}
	if cert.KeyUsage&x509.KeyUsageEncipherOnly != 0 {
		keyUsage = append(keyUsage, "Encipher Only")
	}
	if cert.KeyUsage&x509.KeyUsageDecipherOnly != 0 {
		keyUsage = append(keyUsage, "Decipher Only")
	}

	// 7. 解析 ExtKeyUsage
	extKeyUsage := []string{}
	for _, usage := range cert.ExtKeyUsage {
		switch usage {
		case x509.ExtKeyUsageServerAuth:
			extKeyUsage = append(extKeyUsage, "Server Authentication")
		case x509.ExtKeyUsageClientAuth:
			extKeyUsage = append(extKeyUsage, "Client Authentication")
		case x509.ExtKeyUsageCodeSigning:
			extKeyUsage = append(extKeyUsage, "Code Signing")
		case x509.ExtKeyUsageEmailProtection:
			extKeyUsage = append(extKeyUsage, "Email Protection")
		case x509.ExtKeyUsageTimeStamping:
			extKeyUsage = append(extKeyUsage, "Time Stamping")
		case x509.ExtKeyUsageOCSPSigning:
			extKeyUsage = append(extKeyUsage, "OCSP Signing")
		}
	}

	// 8. IP 地址转字符串
	ipAddresses := []string{}
	for _, ip := range cert.IPAddresses {
		ipAddresses = append(ipAddresses, ip.String())
	}

	// 9. 如果提供了私钥，验证证书和私钥是否配对
	keyPairMatched := false
	keyPairChecked := false

	if privateKeyContent != "" {
		keyPairChecked = true
		matched, err := verifyKeyPair(cert, privateKeyContent)
		if err == nil {
			keyPairMatched = matched
		}
	}

	return &ValidateCertResponse{
		Valid:              true,
		Version:            cert.Version,
		SerialNumber:       cert.SerialNumber.String(),
		Issuer:             cert.Issuer.String(),
		Subject:            cert.Subject.String(),
		CommonName:         cert.Subject.CommonName,
		Organization:       organization,
		Country:            country,
		DNSNames:           cert.DNSNames,
		IPAddresses:        ipAddresses,
		EmailAddresses:     cert.EmailAddresses,
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		DaysLeft:           daysLeft,
		IsExpired:          isExpired,
		IsNotYetValid:      isNotYetValid,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: pubKeyAlgo,
		PublicKeySize:      pubKeySize,
		KeyUsage:           keyUsage,
		ExtKeyUsage:        extKeyUsage,
		IsCA:               cert.IsCA,
		KeyPairMatched:     keyPairMatched,
		KeyPairChecked:     keyPairChecked,
	}, nil
}

// verifyKeyPair 验证证书和私钥是否配对
func verifyKeyPair(cert *x509.Certificate, privateKeyPEM string) (bool, error) {
	// 解码私钥 PEM
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return false, fmt.Errorf("无效的私钥 PEM 格式")
	}

	var privateKey interface{}
	var err error

	// 尝试不同的私钥格式
	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	default:
		return false, fmt.Errorf("不支持的私钥类型: %s", block.Type)
	}

	if err != nil {
		return false, fmt.Errorf("解析私钥失败: %v", err)
	}

	// 比较公钥
	switch certPubKey := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		privKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return false, nil
		}
		// 比较 RSA 公钥的模数和指数
		return certPubKey.N.Cmp(privKey.N) == 0 && certPubKey.E == privKey.E, nil

	case *ecdsa.PublicKey:
		privKey, ok := privateKey.(*ecdsa.PrivateKey)
		if !ok {
			return false, nil
		}
		// 比较 ECDSA 公钥的曲线和坐标
		return certPubKey.Curve == privKey.Curve &&
			certPubKey.X.Cmp(privKey.X) == 0 &&
			certPubKey.Y.Cmp(privKey.Y) == 0, nil

	default:
		return false, fmt.Errorf("不支持的公钥类型")
	}
}
