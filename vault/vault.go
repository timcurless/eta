package vault

import (
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
