package websocket

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var clientAlias *ClientAlias
var onceAlias sync.Once

type ClientAlias struct {
	Alias map[string][]string `json:"alias"`
}

func (ca *ClientAlias) Add(clientIdent string, alias string) error {
	if _, ok := ca.Alias[alias]; ok {
		ca.Alias[alias] = append(ca.Alias[alias], clientIdent)
	} else {
		ca.Alias[alias] = []string{clientIdent}
	}
	return ca.Dumps()
}

func (ca *ClientAlias) Remove(clientIdent string, alias string) error {
	values := ca.Alias[alias]
	for i, v := range values {
		if v != clientIdent {
			continue
		}
		values[1] = values[i]
		ca.Alias[alias] = values[1:]
		return ca.Dumps()
	}
	return nil
}

func (ca *ClientAlias) Has(alias string) bool {
	if _, ok := ca.Alias[alias]; true {
		return ok
	}
	return false
}

func (ca *ClientAlias) Loads() error {
	if data, err := ioutil.ReadFile("client_alias"); err == nil {
		if len(data) != 0 {
			json.Unmarshal(data, ca)
		}
		return nil
	} else {
		return err
	}
}

func (ca *ClientAlias) Dumps() error {
	if data, err := json.Marshal(&ca); err == nil {
		ioutil.WriteFile("client_alias", data, 0644)
		return nil
	} else {
		return err
	}
}

func NewClientAlias() *ClientAlias {
	onceAlias.Do(func() {
		clientAlias = &ClientAlias{}
		clientAlias.Alias = make(map[string][]string)
		clientAlias.Loads()
	})
	return clientAlias
}
