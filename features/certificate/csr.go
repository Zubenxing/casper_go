package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
)

// CSRInfo CSR 信息结构
type CSRInfo struct {
	Subject            pkix.Name
	DNSNames           []string
	EmailAddresses     []string
	IPAddresses        []string
	SignatureAlgorithm string
	PublicKeyAlgorithm string
	PublicKeySize      int
}

// GenerateCSRRequest CSR 生成请求
type GenerateCSRRequest struct {
	CommonName         string   `json:"common_name" binding:"required"`   // 域名/通用名
	Country            string   `json:"country"`                          // 国家代码，如 CN
	Province           string   `json:"province"`                         // 省份
	Locality           string   `json:"locality"`                         // 城市
	Organization       string   `json:"organization"`                     // 组织名
	OrganizationalUnit string   `json:"organizational_unit"`              // 部门
	EmailAddress       string   `json:"email_address"`                    // 邮箱
	DNSNames           []string `json:"dns_names"`                        // SAN - DNS 名称
	IPAddresses        []string `json:"ip_addresses"`                     // SAN - IP 地址
	KeyAlgorithm       string   `json:"key_algorithm" binding:"required"` // RSA 或 ECDSA
	KeySize            int      `json:"key_size"`                         // RSA: 2048/4096, ECDSA: 256/384
}

// GenerateCSRResponse CSR 生成响应
type GenerateCSRResponse struct {
	PrivateKey string `json:"private_key"` // PEM 格式私钥
	CSR        string `json:"csr"`         // PEM 格式 CSR
}

// ValidateCSRRequest CSR 校验请求
type ValidateCSRRequest struct {
	CSRContent string `json:"csr_content" binding:"required"` // PEM 格式 CSR 内容
}

// ValidateCSRResponse CSR 校验响应
type ValidateCSRResponse struct {
	Valid              bool     `json:"valid"`
	CommonName         string   `json:"common_name"`
	Country            string   `json:"country"`
	Province           string   `json:"province"`
	Locality           string   `json:"locality"`
	Organization       string   `json:"organization"`
	OrganizationalUnit string   `json:"organizational_unit"`
	DNSNames           []string `json:"dns_names"`
	EmailAddresses     []string `json:"email_addresses"`
	SignatureAlgorithm string   `json:"signature_algorithm"`
	PublicKeyAlgorithm string   `json:"public_key_algorithm"`
	PublicKeySize      int      `json:"public_key_size"`
	ErrorMessage       string   `json:"error_message,omitempty"`
}

// GenerateCSR 生成 CSR 和私钥
func GenerateCSR(req *GenerateCSRRequest) (*GenerateCSRResponse, error) {
	// 1. 生成私钥
	var privateKey interface{}
	var err error

	switch req.KeyAlgorithm {
	case "RSA":
		keySize := req.KeySize
		if keySize != 2048 && keySize != 4096 {
			keySize = 2048 // 默认 2048
		}
		privateKey, err = rsa.GenerateKey(rand.Reader, keySize)
		if err != nil {
			return nil, fmt.Errorf("生成 RSA 私钥失败: %w", err)
		}

	case "ECDSA":
		var curve elliptic.Curve
		switch req.KeySize {
		case 256:
			curve = elliptic.P256()
		case 384:
			curve = elliptic.P384()
		default:
			curve = elliptic.P256() // 默认 P-256
		}
		privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("生成 ECDSA 私钥失败: %w", err)
		}

	default:
		return nil, fmt.Errorf("不支持的密钥算法: %s (仅支持 RSA 或 ECDSA)", req.KeyAlgorithm)
	}

	// 2. 构建证书请求主题
	subject := pkix.Name{
		CommonName: req.CommonName,
	}
	if req.Country != "" {
		subject.Country = []string{req.Country}
	}
	if req.Province != "" {
		subject.Province = []string{req.Province}
	}
	if req.Locality != "" {
		subject.Locality = []string{req.Locality}
	}
	if req.Organization != "" {
		subject.Organization = []string{req.Organization}
	}
	if req.OrganizationalUnit != "" {
		subject.OrganizationalUnit = []string{req.OrganizationalUnit}
	}

	// 3. 创建 CSR 模板
	template := x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames:           req.DNSNames,
	}

	// 设置签名算法
	if req.KeyAlgorithm == "ECDSA" {
		template.SignatureAlgorithm = x509.ECDSAWithSHA256
	}

	// 添加邮箱
	if req.EmailAddress != "" {
		template.EmailAddresses = []string{req.EmailAddress}
	}

	// 4. 生成 CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return nil, fmt.Errorf("创建 CSR 失败: %w", err)
	}

	// 5. 编码为 PEM 格式
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	// 6. 编码私钥为 PEM 格式
	var privateKeyPEM []byte
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		privateKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
	case *ecdsa.PrivateKey:
		keyBytes, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("编码 ECDSA 私钥失败: %w", err)
		}
		privateKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: keyBytes,
		})
	}

	return &GenerateCSRResponse{
		PrivateKey: string(privateKeyPEM),
		CSR:        string(csrPEM),
	}, nil
}

// ValidateCSR 校验 CSR 文件
func ValidateCSR(csrContent string) (*ValidateCSRResponse, error) {
	// 1. 解码 PEM
	block, _ := pem.Decode([]byte(csrContent))
	if block == nil {
		return &ValidateCSRResponse{
			Valid:        false,
			ErrorMessage: "无效的 PEM 格式",
		}, nil
	}

	if block.Type != "CERTIFICATE REQUEST" && block.Type != "NEW CERTIFICATE REQUEST" {
		return &ValidateCSRResponse{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("无效的 PEM 类型: %s (期望 CERTIFICATE REQUEST)", block.Type),
		}, nil
	}

	// 2. 解析 CSR
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return &ValidateCSRResponse{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("解析 CSR 失败: %v", err),
		}, nil
	}

	// 3. 验证签名
	if err := csr.CheckSignature(); err != nil {
		return &ValidateCSRResponse{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("CSR 签名验证失败: %v", err),
		}, nil
	}

	// 4. 提取信息
	country := ""
	if len(csr.Subject.Country) > 0 {
		country = csr.Subject.Country[0]
	}

	province := ""
	if len(csr.Subject.Province) > 0 {
		province = csr.Subject.Province[0]
	}

	locality := ""
	if len(csr.Subject.Locality) > 0 {
		locality = csr.Subject.Locality[0]
	}

	organization := ""
	if len(csr.Subject.Organization) > 0 {
		organization = csr.Subject.Organization[0]
	}

	organizationalUnit := ""
	if len(csr.Subject.OrganizationalUnit) > 0 {
		organizationalUnit = csr.Subject.OrganizationalUnit[0]
	}

	// 获取公钥算法和大小
	pubKeyAlgo := ""
	pubKeySize := 0

	switch pubKey := csr.PublicKey.(type) {
	case *rsa.PublicKey:
		pubKeyAlgo = "RSA"
		pubKeySize = pubKey.N.BitLen()
	case *ecdsa.PublicKey:
		pubKeyAlgo = "ECDSA"
		pubKeySize = pubKey.Curve.Params().BitSize
	default:
		pubKeyAlgo = "Unknown"
	}

	return &ValidateCSRResponse{
		Valid:              true,
		CommonName:         csr.Subject.CommonName,
		Country:            country,
		Province:           province,
		Locality:           locality,
		Organization:       organization,
		OrganizationalUnit: organizationalUnit,
		DNSNames:           csr.DNSNames,
		EmailAddresses:     csr.EmailAddresses,
		SignatureAlgorithm: csr.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: pubKeyAlgo,
		PublicKeySize:      pubKeySize,
	}, nil
}
