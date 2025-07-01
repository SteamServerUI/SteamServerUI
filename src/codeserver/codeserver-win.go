//go:build windows

package codeserver

import "fmt"

func InitCodeServer() error {
	return fmt.Errorf("Code server can only be used on Linux")
}
