package use_case

type FileSystem interface {
  DoesFileExist(path string) bool
  GetHomeDir() (string, error)
  ReadFile(path string) ([]byte, error)
  RemoveFile(path string) error
  WriteToFile(path string, data []byte) error
}

type Keychain interface {
  DoesSecretExist(label string) (bool, error)
  GetSecret(label string) (string, error)
  SetSecret(secret string, label string) error
}
