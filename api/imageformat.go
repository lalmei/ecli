package api

func ImageFormats() (map[string]interface{}, error) {
	// No explicit parameter, but will add the token.
	p := map[string]interface{}{}
	v, err := sendRequest("KeenEye.ImageFormats", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
