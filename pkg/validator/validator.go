package validator

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
)

// ValidateHost 验证主机名或IP地址
func ValidateHost(host string) error {
	if host == "" {
		return fmt.Errorf("主机名不能为空")
	}

	// 检查是否为有效的IP地址
	if net.ParseIP(host) != nil {
		return nil
	}

	// 检查是否为有效的主机名
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !hostnameRegex.MatchString(host) {
		return fmt.Errorf("无效的主机名格式")
	}

	return nil
}

// ValidatePort 验证端口号
func ValidatePort(port string) error {
	if port == "" {
		return fmt.Errorf("端口号不能为空")
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("端口号必须是数字")
	}

	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("端口号必须在1-65535之间")
	}

	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	return nil
}

// ValidatePassword 验证密码（可以为空，但如果提供则需要验证）
func ValidatePassword(password string) error {
	// 密码可以为空
	return nil
}

// ValidateDatabaseName 验证数据库名
func ValidateDatabaseName(dbName string) error {
	// 数据库名可以为空
	return nil
}
