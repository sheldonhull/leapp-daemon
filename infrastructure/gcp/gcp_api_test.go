package gcp

import (
  "os"
  "testing"
)

func TestWorkingDefaultAuth(t *testing.T) {
  if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
    t.Fatalf("GOOGLE_APPLICATION_CREDENTIALS should not be set")
  }

  keys, err := listKeys("leapp-test@forlunch-laboratory.iam.gserviceaccount.com")
  if err != nil {
    t.Fatalf("Error: %v", err)
  }
  print(keys)
}

func TestOAuthGetAuthUrl(t *testing.T) {
  config := oauthGetConfig()
  url := oauthGetAuthorizationUrl(config)
  print(url)
}

func TestOAuthGetToken(t *testing.T) {
  config := oauthGetConfig()
  authCode := ""

  token := oauthGetTokenFromWeb(config, authCode)
  print("refresh token: " + token.RefreshToken + "\n")
  print("token: " + token.AccessToken + "\n")
}
