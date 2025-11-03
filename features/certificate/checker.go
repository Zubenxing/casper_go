package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"time"
)

// CertInfo 证书信息
type CertInfo struct {
	Domain       string
	Issuer       string
	Subject      string
	Organization string
	NotBefore    time.Time
	NotAfter     time.Time
	DaysLeft     int
	IsValid      bool
	Error        string
}

// Check 检查 URL 的证书信息
func Check(targetURL string) (*CertInfo, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("无效的 URL: %w", err)
	}

	host := parsedURL.Host
	if parsedURL.Port() == "" {
		host = parsedURL.Hostname() + ":443"
	}

	// 严格模式：正常验证证书
	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return &CertInfo{
			Domain:  parsedURL.Hostname(),
			IsValid: false,
			Error:   fmt.Sprintf("TLS连接失败: %v", err),
		}, nil
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("未找到证书")
	}

	cert := certs[0]
	daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)

	organization := ""
	if len(cert.Subject.Organization) > 0 {
		organization = cert.Subject.Organization[0]
	}

	// 严格验证：检查证书链和域名
	isValid := true
	errorMsg := ""

	opts := x509.VerifyOptions{
		DNSName: parsedURL.Hostname(),
	}
	if _, err := cert.Verify(opts); err != nil {
		isValid = false
		errorMsg = fmt.Sprintf("证书验证失败: %v", err)
	}

	// 检查证书时间有效性
	now := time.Now()
	if now.Before(cert.NotBefore) {
		isValid = false
		if errorMsg == "" {
			errorMsg = "证书尚未生效"
		}
	} else if now.After(cert.NotAfter) {
		isValid = false
		if errorMsg == "" {
			errorMsg = "证书已过期"
		}
	}

	return &CertInfo{
		Domain:       parsedURL.Hostname(),
		Issuer:       cert.Issuer.CommonName,
		Subject:      cert.Subject.CommonName,
		Organization: organization,
		NotBefore:    cert.NotBefore,
		NotAfter:     cert.NotAfter,
		DaysLeft:     daysLeft,
		IsValid:      isValid,
		Error:        errorMsg,
	}, nil
}

// GetStatus 检查证书状态
func GetStatus(daysLeft int, warningDays int) int {
	if daysLeft < 0 {
		return 3 // 过期
	} else if daysLeft <= warningDays {
		return 2 // 警告
	}
	return 1 // 正常
}
