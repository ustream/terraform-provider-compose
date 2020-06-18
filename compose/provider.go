package compose

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ustream/terraform-provider-compose/composeapi"
)

type Config struct {
	BluemixAPIKey string
	Region        string
}

func (c *Config) NewClient() (*composeapi.Client, error) {
	var apiBase string
	switch c.Region {
	case "us-east":
		apiBase = composeapi.BxUsEastApiBase
	case "eu-de":
		apiBase = composeapi.BxEuDeApiBase
	case "eu-gb":
		apiBase = composeapi.BxEuGbApiBase
	default:
		apiBase = composeapi.BxUsSouthApiBase
	}
	client, err := composeapi.NewClient(c.BluemixAPIKey, apiBase)
	return client, err
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bluemix_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Bluemix API Key",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_API_KEY", "BLUEMIX_API_KEY", "IBMCLOUD_API_KEY", "IC_API_KEY"}, ""),
				Sensitive:   true,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Bluemix Region (for example 'us-south').",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_REGION", "BLUEMIX_REGION", "IBMCLOUD_REGION", "IC_REGION"}, "us-south"),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"compose_whitelist": resourceComposeWhitelist(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	bluemixAPIKey := d.Get("bluemix_api_key").(string)
	region := d.Get("region").(string)

	config := Config{
		BluemixAPIKey: bluemixAPIKey,
		Region:        region,
	}

	return config.NewClient()
}
