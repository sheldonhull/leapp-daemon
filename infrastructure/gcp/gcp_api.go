package gcp

import (
  "context"
  "fmt"
  "google.golang.org/api/iam/v1"
  "google.golang.org/api/option"
)

// For reference see:
// https://cloud.google.com/iam/docs/creating-managing-service-account-keys#iam-service-account-keys-create-go

func GetServiceAccountKey(apiKey string, serviceAccountName string, projectId string) (*iam.ServiceAccountKey, error) {
  ctx := context.Background()
  service, err := iam.NewService(ctx, option.WithAPIKey(apiKey))
  if err != nil {
    return nil, fmt.Errorf("iam.NewService: %v", err)
  }
  resource := "projects/-/serviceAccounts/" + serviceAccountName + "@" + projectId + ".iam.gserviceaccount.com"
  request := &iam.CreateServiceAccountKeyRequest{}
  key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
  if err != nil {
    return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
  }
  return key, nil
}
