package utils

import (
	"io"
	"os/exec"
)

//crypt byte in package pem PKCS7
func EncryptPackagePKCS7(data []byte, cert string, privateKey string, certPassword string) ([]byte, error) {
	path := ExistCliCommand("openssl")
	cmd := exec.Command(path, "smime", "-sign", "-signer", cert, "-inkey", privateKey, "-nochain", "-nocerts", "-outform", "PEM", "-nodetach", "-passin", "pass:", certPassword)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(data))
	}()

	return cmd.CombinedOutput()
}

// decrypt byte pkcs7 in byte
func DecryptPackagePKCS7(data []byte, pathCert string) ([]byte, error) {
	path := ExistCliCommand("openssl")
	cmd := exec.Command(path, "smime", "-verify", "-inform", "PEM", "-nointern", "-certfile", pathCert, "-CAfile", pathCert)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(data))
	}()

	return cmd.CombinedOutput()
}
