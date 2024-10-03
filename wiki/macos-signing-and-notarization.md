# Signing and Notarizing (MacOS)

To ensure the CI pipeline can correctly sign and notarize macOS applications, we need to set up Apple Developer environment variables and obtain a Developer ID Application certificate. This setup is necessary for releasing macOS versions outside the Mac App Store and complying with Apple's security requirements.

## Apple Developer Environment Variables

To sign and notarize macOS applications, the following environment variables need to be configured:

| Variable                  | Description                                                                                                   |
|---------------------------|---------------------------------------------------------------------------------------------------------------|
| `AC_CERTIFICATE`           | Base64-encoded `.p12` certificate used for code signing macOS applications.                                   |
| `AC_CERTIFICATE_PASSWORD`  | Password for the `.p12` certificate used for signing.                                                         |
| `AC_APPLE_ID`              | Apple Developer account email used for notarization.                                                          |
| `AC_PASSWORD`              | App-specific password for your Apple Developer account (generated via your Apple ID settings).                |
| `AC_TEAM_ID`               | Apple Developer Team ID, which identifies your team in the Apple Developer Program.                           |
| `AC_DEVELOPER_ID`          | The name or ID of the "Developer ID Application" certificate used to sign the macOS application.              |

## Creating a Developer ID Application Certificate

To sign and notarize macOS applications, you'll need a **Developer ID Application Certificate** from Apple. Here's how to create and export this certificate:

### Step 1: Enroll in the Apple Developer Program

You must be a member of the [Apple Developer Program](https://developer.apple.com/programs/), which requires a paid membership.

### Step 2: Create a Certificate Signing Request (CSR)

1. Open **Keychain Access** on your Mac (located in `/Applications/Utilities/Keychain Access.app`).
2. From the **Keychain Access** menu, select **Certificate Assistant -> Request a Certificate from a Certificate Authority**.
3. Fill in the required details:
   - **User Email Address**: Your Apple Developer account email.
   - **Common Name**: Your name or your organization’s name.
   - **CA Email Address**: Leave this blank.
   - **Request Is**: Select **Saved to Disk**.
4. Save the certificate request (CSR) file to your local machine.

### Step 3: Request the Developer ID Application Certificate

1. Log in to the [Apple Developer Certificates Portal](https://developer.apple.com/account/resources/certificates/list) using your Apple Developer account.
2. In the **Certificates, Identifiers & Profiles** section, click the **"+"** button to create a new certificate.
3. Choose **Developer ID Application** under the **Production** section.
4. Upload the **CSR** file you saved in the previous step.
5. Download the certificate (`.cer` file) provided by Apple and double-click it to install it into your Keychain.

### Step 4: Export the Developer ID Certificate as `.p12`

1. Open **Keychain Access** again and find your **Developer ID Application** certificate under **My Certificates**.
2. Right-click on the certificate and choose **Export**.
3. Save the file as a **.p12** (PKCS12) file. You’ll be prompted to set a password for the file — remember this password as it will be used in the CI pipeline for signing.
4. Base64-encode the `.p12` file so that it can be used in CI:
    ```bash
    openssl base64 -in your_cert.p12 -out cert_base64.txt
    ```

5. Copy the contents of `cert_base64.txt` to use as the value for the **`AC_CERTIFICATE`** environment variable.

## GitHub Actions Secrets Configuration

Once you've created and exported the Developer ID certificate, you need to add the following environment variables as GitHub Secrets:

| Secret Name               | Value (Example)                                           |
|---------------------------|-----------------------------------------------------------|
| `AC_CERTIFICATE`           | Base64-encoded certificate (`cert_base64.txt` content).   |
| `AC_CERTIFICATE_PASSWORD`  | Your `.p12` file password.                               |
| `AC_APPLE_ID`              | Your Apple Developer email (e.g., `you@example.com`).     |
| `AC_PASSWORD`              | App-specific password for your Apple ID.                 |
| `AC_TEAM_ID`               | Your 10-character Apple Developer Team ID (e.g., `ABCD123456`). |
| `AC_DEVELOPER_ID`          | Developer ID Application certificate name (e.g., `Developer ID Application: Your Name (TEAM_ID)`). |

### Notes

- Ensure all environment variables are set correctly before running macOS builds in your CI pipeline.
- If you are using multiple Apple Developer teams, specify the correct `AC_TEAM_ID` and `AC_DEVELOPER_ID` for the application being signed and notarized.
- PR CI builds are not signed or notarized as notarization can take between 5 - 15+ minutes.
