package gcp

import (
  "context"
  "fmt"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/iam/v1"
  "io/ioutil"
  "log"
)

func getService() (*iam.Service, error) {
  ctx := context.Background()
  service, err := iam.NewService(ctx)
  if err != nil {
    return nil, fmt.Errorf("iam.NewService: %v", err)
  }
  return service, nil
}

func createKey(serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
  service, err := getService()
  if err != nil {
    return nil, fmt.Errorf("cannot access GCP service: %v", err)
  }
  resource := "projects/-/serviceAccounts/" + serviceAccountEmail
  request := &iam.CreateServiceAccountKeyRequest{
    KeyAlgorithm:    "",
    PrivateKeyType:  "",
    ForceSendFields: nil,
    NullFields:      nil,
  }
  key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
  if err != nil {
    return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
  }
  return key, nil
}

func listKeys(serviceAccountEmail string) ([]*iam.ServiceAccountKey, error) {
  service, err := getService()
  if err != nil {
    return nil, fmt.Errorf("cannot access GCP service: %v", err)
  }
  resource := "projects/-/serviceAccounts/" + serviceAccountEmail
  response, err := service.Projects.ServiceAccounts.Keys.List(resource).Do()
  if err != nil {
    return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.List: %v", err)
  }
  return response.Keys, nil
}

// deleteKey deletes a service account key.
func deleteKey(fullKeyName string) error {
  service, err := getService()
  if err != nil {
    return fmt.Errorf("cannot access GCP service: %v", err)
  }

  _, err = service.Projects.ServiceAccounts.Keys.Delete(fullKeyName).Do()
  if err != nil {
    return fmt.Errorf("Projects.ServiceAccounts.Keys.Delete: %v", err)
  }
  return nil
}

func oauthGetConfig() *oauth2.Config {
  b, err := ioutil.ReadFile("~/leapp/google-cloud-sdk-config.json")
  if err != nil {
    log.Fatalf("Unable to read client secret file: %v", err)
  }

  config, err := google.ConfigFromJSON(b, "openid",
    "https://www.googleapis.com/auth/userinfo.email",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/appengine.admin",
    "https://www.googleapis.com/auth/compute",
    "https://www.googleapis.com/auth/accounts.reauth")
  if err != nil {
    log.Fatalf("Unable to read configuration .json: %v", err)
  }
  return config
}

func oauthGetAuthorizationUrl(config *oauth2.Config) string {
  return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func oauthGetTokenFromWeb(config *oauth2.Config, authCode string) *oauth2.Token {
  ctx := context.Background()
  tok, err := config.Exchange(ctx, authCode)
  if err != nil {
    log.Fatalf("Unable to retrieve token from web: %v", err)
  }
  return tok
}
