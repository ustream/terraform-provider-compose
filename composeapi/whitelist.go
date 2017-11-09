package composeapi

import (
	"encoding/json"
	"fmt"
)

type Instance struct {
	ID           string `json:"id"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

// Whitelist structure
type Whitelist struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	IP          string `json:"ip"`
	Embedded    struct {
		Whitelist []Whitelist `json:"whitelist"`
	} `json:"_embedded"`
}

type whitelistDeployment struct {
	Whitelist Whitelist `json:"whitelist"`
}

type whitelistParams struct {
	Deployment whitelistDeployment `json:"deployment"`
}

//GetRecipeJSON Gets raw JSON for recipeid
func (c *Client) getWhitelistForDeploymentJSON(deploymentid string) (string, []error) {
	instanceEndpointPrefix, errs := c.getInstanceEndpointPrefix(deploymentid)
	if errs != nil {
		return "", errs
	}
	return c.getJSON(instanceEndpointPrefix + "/whitelist")
}

//GetRecipe gets status of Recipe
func (c *Client) GetWhitelistForDeployment(deploymentid string) (*Whitelist, []error) {
	body, errs := c.getWhitelistForDeploymentJSON(deploymentid)

	if errs != nil {
		return nil, errs
	}

	whitelist := Whitelist{}
	json.Unmarshal([]byte(body), &whitelist)

	return &whitelist, nil
}

//GetRecipe gets status of Recipe
func (c *Client) AddWhitelistForDeployment(deploymentid string, whitelist Whitelist) (*Recipe, []error) {
	instanceEndpointPrefix, errs := c.getInstanceEndpointPrefix(deploymentid)
	if errs != nil {
		return nil, errs
	}

	response, body, errs := c.newRequest("POST", c.apiBase+instanceEndpointPrefix+"/whitelist").
		Send(whitelistParams{Deployment: whitelistDeployment{Whitelist: whitelist}}).
		End()

	if errs != nil {
		return nil, errs
	}

	recipe := Recipe{}

	if response.StatusCode != 202 { // Expect Accepted on success - assume error on anything else
		myerrors := Errors{}
		err := json.Unmarshal([]byte(body), &myerrors)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to parse error - status code %d - body %s", response.StatusCode, response.Body))
		} else {
			errs = append(errs, fmt.Errorf("%v", myerrors.Error))
		}
	} else {
		err := json.Unmarshal([]byte(body), &recipe)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to parse recipe - status code %d - body %s", response.StatusCode, response.Body))
		}
	}

	return &recipe, errs
}

//GetRecipe gets status of Recipe
func (c *Client) DeleteWhitelistForDeployment(deploymentid string, whitelistID string) (*Recipe, []error) {
	instanceEndpointPrefix, errs := c.getInstanceEndpointPrefix(deploymentid)
	if errs != nil {
		return nil, errs
	}

	response, body, errs := c.newRequest("DELETE", c.apiBase+instanceEndpointPrefix+"/whitelist/"+whitelistID).End()

	if errs != nil {
		return nil, errs
	}

	recipe := Recipe{}

	if response.StatusCode != 202 { // Expect Accepted on success - assume error on anything else
		myerrors := Errors{}
		err := json.Unmarshal([]byte(body), &myerrors)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to parse error - status code %d - body %s", response.StatusCode, response.Body))
		} else {
			errs = append(errs, fmt.Errorf("%v", myerrors.Error))
		}
	} else {
		err := json.Unmarshal([]byte(body), &recipe)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to parse recipe - status code %d - body %s", response.StatusCode, response.Body))
		}
	}

	return &recipe, errs
}
