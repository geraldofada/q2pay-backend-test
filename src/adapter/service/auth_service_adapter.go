package service

import (
	"encoding/json"
	"net/http"

	"github.com/geraldofada/q2pay-backend-test/src/core"
)

type Service struct{}

func (s Service) AuthorizeTransfer() (bool, error) {
	// NOTE: Essa URI poderia ser extraída pro .env, mas como é apenas uma chamada simples
	// vou deixar direto aqui
	resp, err := http.Get("https://run.mocky.io/v3/d02168c6-d88d-4ff2-aac6-9e9eb3425e31")
	if err != nil {
		return false, core.AccountTransferNotAuthorized{}
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return false, core.AccountTransferNotAuthorized{}
	}

	if body["authorization"] == nil || !body["authorization"].(bool) {
		return false, core.AccountTransferNotAuthorized{}
	}

	return true, nil
}
