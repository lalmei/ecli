package api

func ImageTypes() (map[string]interface{}, error) {
	v, err := sendRequest("KeenEye.ImageTypes", nil)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
