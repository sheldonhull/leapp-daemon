package use_case

import (
	"encoding/json"
	"io/ioutil"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/infrastructure/logging"
)

// temporaneamente forzo ad AK
const mode string = "AK"

type CredentialsFile struct {
	Current   string                `json:"current"`
	Profiles  []NamedProfileSection `json:"profiles"`
	Meta_path string                `json:"meta_path"`
}

type NamedProfileSection struct {
	Name              string `json:"name"`
	Mode              string `json:"mode"`
	Access_key_id     string `json:"access_key_id"`
	Access_key_secret string `json:"access_key_secret"`
	Sts_token         string `json:"sts_token"`
	Ram_role_name     string `json:"ram_role_name"`
	Ram_role_arn      string `json:"ram_role_arn"`
	Ram_session_name  string `json:"ram_session_name"`
	Private_key       string `json:"private_key"`
	Key_pair_name     string `json:"key_pair_name"`
	Expired_seconds   int    `json:"expired_seconds"`
	Verified          string `json:"verified"`
	Region_id         string `json:"region_id"`
	Output_format     string `json:"output_format"`
	Language          string `json:"language"`
	Site              string `json:"site"`
	Retry_timeout     int    `json:"retry_timeout"`
	Connect_timeout   int    `json:"connect_timeout"`
	Retry_count       int    `json:"retry_count"`
	Process_command   string `json:"process_command"`
}

type AlibabaCredentialsApplier struct {
	FileSystem FileSystem
	Keychain   Keychain
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) UpdatePlainAlibabaSessions(oldPlainAlibabaSessions []session.PlainAlibabaSession, newPlainAlibabaSessions []session.PlainAlibabaSession) error {
	for i, oldSess := range oldPlainAlibabaSessions {
		if i < len(newPlainAlibabaSessions) {
			newSess := newPlainAlibabaSessions[i]

			if oldSess.Status == session.NotActive && newSess.Status == session.Pending {

				homeDir, err := alibabaCredentialsApplier.FileSystem.GetHomeDir()
				if err != nil {
					return err
				}

				credentialsFilePath := homeDir + "/" + constant.AlibabaCredentialsFilePath
				profileName := named_profile.GetNamedProfilesFacade().GetNamedProfileById(newSess.Account.NamedProfileId).Name
				region := newSess.Account.Region

				accessKeyId, secretAccessKey, err := alibabaCredentialsApplier.getAccessKeys(newSess.Id)
				if err != nil {
					return err
				}

				namedProfileSection := NamedProfileSection{Name: profileName, Mode: mode, Access_key_id: accessKeyId, Access_key_secret: secretAccessKey, Region_id: region, Output_format: "json", Language: "en"}
				profiles := []NamedProfileSection{namedProfileSection}
				config := CredentialsFile{Current: namedProfileSection.Name, Profiles: profiles}
				out, _ := json.MarshalIndent(config, "", "  ")
				alibabaCredentialsApplier.overwriteFile(out, credentialsFilePath)
			}
		}
	}

	return nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) getAccessKeys(sessionId string) (accessKeyId string, secretAccessKey string, error error) {
	accessKeyId = ""
	secretAccessKey = ""

	accessKeyIdSecretName := sessionId + "-plain-alibaba-session-access-key-id"
	accessKeyId, err := alibabaCredentialsApplier.Keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := sessionId + "-plain-alibaba-session-secret-access-key"
	secretAccessKey, err = alibabaCredentialsApplier.Keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) overwriteFile(content []byte, path string) error {
	logging.Entry().Error("flag 3")

	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	return nil
}
