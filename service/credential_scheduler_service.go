package service

import "leapp_daemon/core/configuration"

func RotateAllSessionsCredentials() error {
	config, err := configuration.ReadConfiguration()
	if err != nil { return err }

	for i := range config.PlainAwsSessions {
		sess := config.PlainAwsSessions[i]

		err = sess.Rotate(config, nil)

		if err != nil {
			return err
		}
	}

	/*for i := range config.FederatedAwsSessions {
		sess := config.FederatedAwsSessions[i]
		if sess.Active {
		}
	}*/

	err = configuration.UpdateConfiguration(config, false)
	if err != nil {
		return err
	}

	return nil
}
