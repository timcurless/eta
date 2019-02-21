package vault

import (
	"strconv"

	"github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *api.Client
}

func NewVaultClient(vaultURL, vaultToken string) (*VaultClient, error) {
	config := api.Config{Address: vaultURL}

	c, err := api.NewClient(&config)
	if err != nil {
		return nil, err
	}

	c.SetToken(vaultToken)
	return &VaultClient{client: c}, nil
}

func (v *VaultClient) GetVaultHealth() (interface{}, error) {
	health, err := v.client.Sys().Health()
	if err != nil {
		return nil, err
	}
	return health, nil
}

func (v *VaultClient) GetDatabaseCreds() (string, string, string, error) {
	s, err := v.client.Logical().Read("database/creds/eta")
	if err != nil {
		return "", "", "", err
	}
	return s.Data["username"].(string),
		s.Data["password"].(string),
		strconv.Itoa(s.LeaseDuration),
		nil
}
