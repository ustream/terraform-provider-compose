package compose

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ustream/terraform-provider-compose/composeapi"
)

func resourceComposeWhitelist() *schema.Resource {
	return &schema.Resource{
		Create:   resourceComposeWhitelistCreate,
		Read:     resourceComposeWhitelistRead,
		Delete:   resourceComposeWhitelistDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, field string) (warnings []string, errors []error) {
					value := val.(string)
					if _, _, err := net.ParseCIDR(value); err != nil {
						errors = append(
							errors,
							fmt.Errorf(
								"Provided value '(%s)' is not a valid IPv4 network: %s",
								value,
								err,
							),
						)
					}
					return
				},
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"deployment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

//  Creates DNS Domain Resource Record
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/createObject
func resourceComposeWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*composeapi.Client)

	deploymentID := d.Get("deployment_id").(string)
	ip := d.Get("ip").(string)
	description := d.Get("description").(string)

	whitelist := composeapi.Whitelist{IP: ip, Description: description}

	recipe, errs := client.AddWhitelistForDeployment(deploymentID, whitelist)

	if errs != nil {
		return fmt.Errorf("Error adding whitelist entry: %s", errs)
	}

	stateChangeConf := &resource.StateChangeConf{
		Pending:    []string{"waiting", "running"},
		Target:     []string{"complete"},
		Refresh:    jobCompletedRefreshFunc(client, deploymentID, recipe.ID),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateChangeConf.WaitForState()
	if err != nil {
		return err
	}

	newWhitelist, errs := client.GetWhitelistForDeployment(deploymentID)

	if errs != nil {
		return fmt.Errorf("Error querying whitelist entry: %s", errs)
	}

	for _, whitelistEntry := range newWhitelist.Embedded.Whitelist {
		if whitelistEntry.Description == description {
			d.SetId(whitelistEntry.ID)
			return resourceComposeWhitelistRead(d, meta)
		}
	}

	return errors.New("Failed to find newly created whitelist entry")
}

func resourceComposeWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*composeapi.Client)

	deploymentID := d.Get("deployment_id").(string)

	whitelist, errs := client.GetWhitelistForDeployment(deploymentID)

	if errs != nil {
		return fmt.Errorf("Error querying whitelist entry: %s", errs)
	}

	for _, whitelistEntry := range whitelist.Embedded.Whitelist {
		if whitelistEntry.ID == d.Id() {
			d.Set("description", whitelistEntry.Description)
			d.Set("ip", whitelistEntry.IP)
			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceComposeWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*composeapi.Client)

	deploymentID := d.Get("deployment_id").(string)

	recipe, errs := client.DeleteWhitelistForDeployment(deploymentID, d.Id())

	if errs != nil {
		return fmt.Errorf("Error deleting whitelist entry: %s", errs)
	}

	stateChangeConf := &resource.StateChangeConf{
		Pending:    []string{"waiting", "running"},
		Target:     []string{"complete"},
		Refresh:    jobCompletedRefreshFunc(client, deploymentID, recipe.ID),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateChangeConf.WaitForState()
	if err != nil {
		return err
	}

	return nil
}

func jobCompletedRefreshFunc(client *composeapi.Client, deploymentid string, recipeid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		recipe, err := client.GetRecipe(deploymentid, recipeid)
		if err != nil {
			return nil, "", err[0]
		}
		return recipe.ID, recipe.Status, nil
	}
}
