package keeneye

// OpenSession returns a token upon successul authentication.
func OpenSession(login, password string) (string, error) {
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

func CloseSession(token string) error {
	p := map[string]interface{}{
		"token": token,
	}
	_, err := sendRequest("KeenEye.CloseSession", p)
	if err != nil {
		return err
	}
	return nil
}
