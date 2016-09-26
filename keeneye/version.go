package keeneye

func Version() (string, error) {
	v, err := sendRequest("KeenEye.Version", nil)
	if err != nil {
		return "", err
	}
	return v.(map[string]interface{})["version"].(string), err
}
