package api

func SlideInfo(id uint64) (map[string]interface{}, error) {
	p := map[string]interface{}{
		"slideId": id,
	}
	v, err := sendRequest("KeenEye.SlideInfo", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
