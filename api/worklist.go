package api

func WorkList(parentId string) (map[string]interface{}, error) {
	p := make(map[string]interface{}, 0)
	if parentId != "" {
		p["parentId"] = parentId
	}
	v, err := sendRequest("KeenEye.WorkList", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
