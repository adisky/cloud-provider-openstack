package barbican

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
)

func (cfg Config) toAuthOptions() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: cfg.Global.AuthURL,
		Username:         cfg.Global.Username,
		UserID:           cfg.Global.UserID,
		Password:         cfg.Global.Password,
		TenantID:         cfg.Global.TenantID,
		TenantName:       cfg.Global.TenantName,
		DomainID:         cfg.Global.DomainID,
		DomainName:       cfg.Global.DomainName,

		// Persistent service, so we need to be able to renew tokens.
		AllowReauth: true,
	}
}

type KMSOpts struct {
	KeyID string `gcfg:"key-id"`
}

//Config to read config options
type Config struct {
	Global struct {
		AuthURL    string `gcfg:"auth-url"`
		Username   string
		UserID     string `gcfg:"user-id"`
		Password   string
		TenantID   string `gcfg:"tenant-id"`
		TenantName string `gcfg:"tenant-name"`
		DomainID   string `gcfg:"domain-id"`
		DomainName string `gcfg:"domain-name"`
		Region     string
	}
	KeyManager KMSOpts
}

// Barbican is gophercloud service client
type Barbican struct {
	Client *gophercloud.ServiceClient
}

// NewBarbicanClient creates new BarbicanClient
func NewBarbicanClient(cfg *Config) (*Barbican, error) {

	provider, err := openstack.AuthenticatedClient(cfg.toAuthOptions())

	if err != nil {
		return nil, err
	}

	client, err := openstack.NewKeyManagerV1(provider, gophercloud.EndpointOpts{
		Region: cfg.Global.Region,
	})
	if err != nil {
		return nil, err
	}

	return &Barbican{Client: client}, nil
}

// GetSecret gets unencrypted secret
func (client *Barbican) GetSecret(keyID string) ([]byte, error) {

	opts := secrets.GetPayloadOpts{
		PayloadContentType: "application/octet-stream",
	}

	key, err := secrets.GetPayload(client.Client, keyID, opts).Extract()
	if err != nil {
		return nil, err
	}

	return key, nil
}
