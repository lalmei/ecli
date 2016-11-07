package api

func ImageTypes() (map[string]interface{}, error) {
	// No explicit parameter, but will add the token.
	p := map[string]interface{}{}
	v, err := sendRequest("KeenEye.ImageTypes", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
