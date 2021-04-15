package credentials_ini_file

import (
  "gopkg.in/ini.v1"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
  "os"
)

func CreateNamedProfileSection(credentialsFile *ini.File, profileName string, accessKeyId string,
	secretAccessKey string, sessionToken string, region string) (*ini.Section, error) {

	section, err := credentialsFile.NewSection(profileName)
	if err != nil {
		return nil, http_error2.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_access_key_id", accessKeyId)
	if err != nil {
		return nil, http_error2.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_secret_access_key", secretAccessKey)
	if err != nil {
		return nil, http_error2.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_session_token", sessionToken)
	if err != nil {
		return nil, http_error2.NewInternalServerError(err)
	}

	if region != "" {
		_, err = section.NewKey("region", region)
		if err != nil {
			return nil, http_error2.NewInternalServerError(err)
		}
	}

	return section, nil
}

func AppendToFile(file *ini.File, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return http_error2.NewNotFoundError(err)
	}

	_, err = file.WriteTo(f)
	if err != nil {
		return http_error2.NewUnprocessableEntityError(err)
	}

	return nil
}

func OverwriteFile(file *ini.File, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return http_error2.NewNotFoundError(err)
	}

	_, err = file.WriteTo(f)
	if err != nil {
		return http_error2.NewUnprocessableEntityError(err)
	}

	return nil
}
