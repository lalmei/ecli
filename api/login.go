package api

// OpenSession returns a token upon successul authentication.
func OpenSession(url, login, password string) (string, error) {
	serverAddress = url
	p := map[string]interface{}{
		"login":    login,
		"password": password,
	}
	v, err := sendRequest("KeenEye.OpenSession", p)
	if err != nil {
		return "", err
	}
	return v.(map[string]interface{})["token"].(string), err
}

func CloseSession() error {
	p := make(map[string]interface{})
	_, err := sendRequest("KeenEye.CloseSession", p)
	if err != nil {
		return err
	}
	return nil
}
