package keeneye

import (
	"ecli/config"
)

func SlideInfo(id uint64) (map[string]interface{}, error) {
	tok, err := config.LoadToken()
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{
		"token":   tok,
		"slideId": id,
	}
	v, err := sendRequest("KeenEye.SlideInfo", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
