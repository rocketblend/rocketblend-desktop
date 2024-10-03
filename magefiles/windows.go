package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// buildWindowsAMD64 builds the Windows AMD64 version of the project.
func buildWindowsAMD64(ldFlags, appVersion string, skipFrontend bool, sign bool) error {
	outputFileName := fmt.Sprintf("rocketblend-desktop-windows-amd64-%s.exe", appVersion)
	skipBindingsFlag, skipFrontendFlag := "", ""
	if skipFrontend {
		skipBindingsFlag, skipFrontendFlag = "-skipbindings", "-s"
	}

	crossCompileFlags := map[string]string{
		"GOOS":   "windows",
		"GOARCH": "amd64",
		"CC":     "x86_64-w64-mingw32-gcc",
		"CXX":    "x86_64-w64-mingw32-g++",
	}

	err := sh.RunWithV(crossCompileFlags, "wails", "build", "-m", "-nosyncgomod", "-ldflags", ldFlags, "-nsis", "-platform", "windows/amd64", "-o", outputFileName, skipBindingsFlag, skipFrontendFlag)
	if err != nil {
		return fmt.Errorf("error building Windows AMD64: %v", err)
	}

	if sign {
		// Recommended to just use Azure Trusted Certificate for signing instead.
		// Keeping for reference.
		fmt.Println("Importing Code Signing Certificates")
		certFilePath := "certificate/certificate.pfx"
		pemFilePath := "certificate/certificate.pem"
		signCert := os.Getenv("WIN_CERTIFICATE")
		signCertPassword := os.Getenv("WIN_CERTIFICATE_PASSWORD")

		if signCert == "" || signCertPassword == "" {
			return fmt.Errorf("missing required environment variables for code-signing")
		}

		if _, err := os.Stat("certificate"); os.IsNotExist(err) {
			if err := os.Mkdir("certificate", os.ModePerm); err != nil {
				return fmt.Errorf("error creating certificate directory: %v", err)
			}
		}

		if err := os.WriteFile("certificate/certificate.txt", []byte(signCert), 0600); err != nil {
			return fmt.Errorf("error writing base64 certificate to file: %v", err)
		}

		if err := sh.Run("certutil", "-decode", "certificate/certificate.txt", certFilePath); err != nil {
			return fmt.Errorf("error decoding certificate: %v", err)
		}

		if err := sh.Run("openssl", "pkcs12", "-in", certFilePath, "-out", pemFilePath, "-nodes", "-passin", fmt.Sprintf("pass:%s", signCertPassword)); err != nil {
			return fmt.Errorf("error converting PFX to PEM: %v", err)
		}

		fmt.Println("Signing Build")
		if err := sh.Run("osslsigncode", "sign", "-certs", pemFilePath, "-key", pemFilePath, "-pass", signCertPassword, "-in", outputFileName, "-out", outputFileName, "-t", "http://timestamp.digicert.com", "-h", "sha256"); err != nil {
			return fmt.Errorf("error signing executable: %v", err)
		}
	}

	return nil
}
