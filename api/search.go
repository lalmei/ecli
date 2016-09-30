package api

func Search(term string) (map[string]interface{}, error) {
	p := map[string]interface{}{
		"term":      term,
		"rawFormat": true,
	}
	v, err := sendRequest("KeenEye.Search", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
