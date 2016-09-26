package keeneye

import (
	"ecli/config"
)

func Search(term string) (map[string]interface{}, error) {
	tok, err := config.LoadToken()
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{
		"token": tok,
		"term":  term,
	}
	v, err := sendRequest("KeenEye.Search", p)
	if err != nil {
		return nil, err
	}
	g := v.(map[string]interface{})
	return g, nil
}
