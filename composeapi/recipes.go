package composeapi

import (
	"encoding/json"
	"strings"
	"time"
)

// Recipe structure
type Recipe struct {
	ID           string    `json:"id"`
	Template     string    `json:"template"`
	Status       string    `json:"status"`
	StatusDetail string    `json:"status_detail"`
	AccountID    string    `json:"account_id"`
	DeploymentID string    `json:"deployment_id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Embedded     struct {
		Recipes []Recipe `json:"recipes"`
	} `json:"_embedded"`
}

type recipeResponse struct {
	Embedded struct {
		Recipes []Recipe `json:"recipes"`
	} `json:"_embedded"`
}

//GetRecipeJSON Gets raw JSON for recipeid
func (c *Client) getRecipeJSON(deploymentid string, recipeid string) (string, []error) {
	if strings.HasPrefix(deploymentid, "bmix-") {
		deploymentid = strings.SplitN(deploymentid, "-", 4)[3]
	}

	return c.getJSON("instances/" + deploymentid + "/recipes/" + recipeid)
}

//GetRecipe gets status of Recipe
func (c *Client) GetRecipe(deploymentid string, recipeid string) (*Recipe, []error) {
	body, errs := c.getRecipeJSON(deploymentid, recipeid)

	if errs != nil {
		return nil, errs
	}

	recipe := Recipe{}
	json.Unmarshal([]byte(body), &recipe)

	return &recipe, nil
}

//GetRecipesForDeploymentJSON returns raw JSON for getRecipesforDeployment
func (c *Client) getRecipesForDeploymentJSON(deploymentid string) (string, []error) {
	instanceEndpointPrefix, errs := c.getInstanceEndpointPrefix(deploymentid)
	if errs != nil {
		return "", errs
	}

	return c.getJSON(instanceEndpointPrefix + "/recipes")
}

//GetRecipesForDeployment gets deployment recipe life
func (c *Client) GetRecipesForDeployment(deploymentid string) (*[]Recipe, []error) {
	body, errs := c.getRecipesForDeploymentJSON(deploymentid)

	if errs != nil {
		return nil, errs
	}

	recipeResponse := Recipe{}
	json.Unmarshal([]byte(body), &recipeResponse)
	recipes := recipeResponse.Embedded.Recipes

	return &recipes, nil
}
