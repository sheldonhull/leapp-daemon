package repository

import (
	"leapp_daemon/test/mock"
	"os"
	"path"
	"testing"
)

var (
	tempDirPath               string
	tempAwsDirPath            string
	tempCredentialsFilePath   string
	awsFsMock                 mock.FileSystemMock
	awsEnvMock                mock.EnvironmentMock
	awsRepo                   AwsConfigurationRepository
	awsCredentials            AwsTempCredentials
	awsCredentialsFileContent string
)

func awsConfigurationRepositorySetup() {
	tempDirPath = path.Join(os.TempDir(), "temp-aws")
	_ = os.RemoveAll(tempDirPath)
	tempAwsDirPath = path.Join(tempDirPath, ".aws")
	_ = os.MkdirAll(tempAwsDirPath, 0600)
	tempCredentialsFilePath = path.Join(tempAwsDirPath, "credentials")

	awsFsMock = mock.FileSystemMock{}
	awsFsMock.ExpHomeDir = tempDirPath

	awsEnvMock = mock.EnvironmentMock{}

	awsRepo = AwsConfigurationRepository{
		FileSystem:  &awsFsMock,
		Environment: &awsEnvMock,
	}

	awsCredentials = AwsTempCredentials{
		ProfileName:  "profile-name",
		AccessKeyId:  "access-key-id",
		SecretKey:    "secret-key",
		SessionToken: "session-token",
		Region:       "region",
	}

	awsCredentialsFileContent = "; Credentials managed by Leapp, do not manually edit this file\n" +
		"[profile-name]\n" +
		"aws_access_key_id     = access-key-id\n" +
		"aws_secret_access_key = secret-key\n" +
		"aws_session_token     = session-token\n" +
		"region                = region\n" +
		"\n"
}

func TestWriteCredentials_missingFile(t *testing.T) {
	awsConfigurationRepositorySetup()

	err := awsRepo.WriteCredentials([]AwsTempCredentials{awsCredentials})
	if err != nil {
		return
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("missing credentials file")
	}

	actualFileContent := string(fileData)
	if actualFileContent != awsCredentialsFileContent {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}
}

func TestWriteCredentials_missingFile_noRegion(t *testing.T) {
	awsConfigurationRepositorySetup()

	awsCredentials.Region = ""
	err := awsRepo.WriteCredentials([]AwsTempCredentials{awsCredentials})
	if err != nil {
		return
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("missing credentials file")
	}

	expectedFileContent := "; Credentials managed by Leapp, do not manually edit this file\n" +
		"[profile-name]\n" +
		"aws_access_key_id     = access-key-id\n" +
		"aws_secret_access_key = secret-key\n" +
		"aws_session_token     = session-token\n" +
		"\n"

	actualFileContent := string(fileData)
	if actualFileContent != expectedFileContent {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}
}

func TestWriteCredentials_missingFile_multipleCredentials(t *testing.T) {
	awsConfigurationRepositorySetup()

	awsCredentials2 := AwsTempCredentials{
		ProfileName:  "profile-name-2",
		AccessKeyId:  "access-key-id-2",
		SecretKey:    "secret-key-2",
		SessionToken: "session-token-2",
		Region:       "region-2",
	}

	err := awsRepo.WriteCredentials([]AwsTempCredentials{awsCredentials, awsCredentials2})
	if err != nil {
		return
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("missing credentials file")
	}

	expectedFileContent := awsCredentialsFileContent +
		"; Credentials managed by Leapp, do not manually edit this file\n" +
		"[profile-name-2]\n" +
		"aws_access_key_id     = access-key-id-2\n" +
		"aws_secret_access_key = secret-key-2\n" +
		"aws_session_token     = session-token-2\n" +
		"region                = region-2\n" +
		"\n"

	actualFileContent := string(fileData)
	if actualFileContent != expectedFileContent {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}
}

func TestWriteCredentials_existingFile(t *testing.T) {
	t.Fatalf("implement other tests")
}
