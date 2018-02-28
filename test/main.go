package main

import (
	"log"
	"os"
	"time"

	"github.com/ustream/terraform-provider-compose/compose"
	"github.com/ustream/terraform-provider-compose/composeapi"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatal("Region must be set as the first argument, DeploymentID as second")
	}

	apiToken := os.Getenv("BM_API_KEY")
	region := os.Args[1]
	deployment := os.Args[2]

	config := compose.Config{
		BluemixAPIKey: apiToken,
		Region:        region,
	}

	client, err := config.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	client.SetLogger(true, os.Stdout)

	recipe, errs := client.AddWhitelistForDeployment(deployment, composeapi.Whitelist{IP: "1.2.3.4/32", Description: "terraform teszt"})

	if errs != nil {
		log.Fatal(errs)
	}

	log.Println(recipe)

	whitelist, errs := client.GetWhitelistForDeployment(deployment)

	if errs != nil {
		log.Fatal(errs)
	}

	log.Println(whitelist.Embedded.Whitelist)

	recipe, errs = client.DeleteWhitelistForDeployment(deployment, whitelist.Embedded.Whitelist[0].ID)

	if errs != nil {
		log.Fatal(errs)
	}

	log.Println(recipe)

	time.Sleep(time.Second * 30)

	recipes, errs := client.GetRecipesForDeployment(deployment)

	if errs != nil {
		log.Fatal(errs)
	}

	log.Println(recipes)
}
