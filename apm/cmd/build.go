/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"example/apm/domain"
	"example/apm/interactive"
	"log"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build up a APM system that can be used by a single tenant.",
	Long: `By following a series of commands, you can build a APM system that can be used by a single tenant.
For example: apm tenant build
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Create meta information to use Grandview.")
		log.Println(" ")
		log.Println("===========================================")
		log.Println("1. Create an Organization. (1/9) ")

		var organizationName, orgID string
		k, _ := interactive.PromptAddOrganization()
		if k == "y" {
			po, err := interactive.PromptOrganization()
			if err != nil {
				log.Panic(err)
			}

			organizationName, orgID, _ = domain.CreateOrganization(po)

			log.Println("Result: " + organizationName)
		} else {
			orgID, organizationName, _ = interactive.SelectOrg()
			log.Println("Result: " + organizationName)
		}

		u, _ := interactive.PromptAddUser()
		if u == "y" {
			for {
				userID, userName, _ := interactive.SelectUser()
				domain.LinkOrgMember(orgID, userID, userName)
				domain.AddOrgMember(orgID, userID, userName)

				log.Println(" ")
				u, _ = interactive.PromptAddUser()
				if u == "n" {
					break
				}
			}
		} else {
			log.Println("You skipped this Step.")
		}

		log.Println(" ")
		log.Println("===========================================")
		log.Println("2. Create a Location. (2/9) ")

		s, _ := interactive.PromptAddLocation()
		if s == "y" {
			for {
				child, parents, depth, _ := interactive.PromptLocation()
				domain.CreateLocation(child, parents, depth)

				log.Println(" ")
				s, _ = interactive.PromptAddMoreLocation()
				if s == "n" {
					break
				}
			}
		} else {
			log.Println("You skipped this Step.")
		}

		log.Println(" ")
		log.Println("===========================================")
		log.Println("3. Create a Tenant. (3/9) ")

		rootCatalogID, rootCatalogName, _ := interactive.SelectRootCatalog()
		log.Println("Selected ID: " + rootCatalogID + ", Name: " + rootCatalogName)

		domain.CreateHierarchy(rootCatalogID)

		var tenantID, tenantName string

		k, _ = interactive.PromptAddTenant()
		if k == "y" {
			tenantID, tenantName, _, _ = domain.CreateTenant(organizationName, rootCatalogID)
			log.Println("Tenant: " + tenantName + ", ID: " + tenantID)

			domain.CreateMenu(tenantID)
			log.Println("The menu was created.")

			users, _ := domain.GetUser(tenantID)
			for _, u := range users {
				domain.CreateRoleOfTenant(organizationName, rootCatalogID, u.Id)
				log.Println("added user: " + u.Id)
			}
			log.Println("user was added to directory.")

		} else {
			tenantID, tenantName, _ = interactive.SelectTenant()
		}

		domain.CreateTypeOfTenant(tenantName)

		log.Println(" ")
		log.Println("===========================================")
		log.Println("4. Create a Type. (ASSET, PARAM) (4/9) ")

		t, _ := interactive.PromptAddType("ASSET")
		if t == "y" {
			for {
				typeName, _ := interactive.PromptType(tenantID, "ASSET")
				typeID, _ := domain.CreateType(tenantID, typeName, "ASSET")
				log.Println("Role:ASSET, Type Name: " + typeName + ", ID: " + typeID)

				log.Println(" ")
				r, _ := interactive.PromptAddMoreType()
				if r == "n" {
					break
				}
			}
		} else {
			log.Println("You skipped this Step. (ASSET) ")
		}

		t, _ = interactive.PromptAddType("PARAM")
		if t == "y" {
			for {
				typeName, _ := interactive.PromptType(tenantID, "PARAM")
				typeID, _ := domain.CreateType(tenantID, typeName, "PARAM")
				log.Println("Role:PARAM, Type Name: " + typeName + ", ID: " + typeID)

				log.Println(" ")
				r, _ := interactive.PromptAddMoreType()
				if r == "n" {
					break
				}
			}
		} else {
			log.Println("You skipped this Step. (PARAM) ")
		}

		log.Println(" ")
		log.Println("===========================================")
		log.Println("5. Create a asset template (5/9) ")

		c, _ := interactive.PromptAddTemplate()
		if c == "y" {
			typeID, _, _ := interactive.SelectType(tenantID, "ASSET")
			templateName, _ := interactive.PromptTemplate()
			templateID, _ := domain.CreateTemplate(tenantID, templateName, typeID)
			log.Println("template ID: " + templateID)
		} else {
			log.Println("You skipped this Step. ")
		}

		log.Println(" ")
		log.Println("===========================================")
		log.Println("6. Create a parameter group. (6/9) ")

		c, _ = interactive.PromptAddParamGroup()
		if c == "y" {
			typeID, _, _ := interactive.SelectType(tenantID, "PARAM")
			paramGroupName, _ := interactive.PromptParamGroup()
			paramGroupID, _ := domain.CreateParamGroup(tenantID, typeID, paramGroupName)
			log.Println("Parameter Group ID: " + paramGroupID)
		} else {
			log.Println("You skipped this Step. ")
		}

		log.Println(" ")
		log.Println("===========================================")
		log.Println("7. Create a parameter with spec. (7/9) ")

		c, _ = interactive.PromptAddParameter()
		if c == "y" {
			paramGroupID, _, _ := interactive.SelectParamGroup(tenantID)

			for {
				parameterName, upper, target, lower, _ := interactive.PromptParameter(paramGroupID)
				parameterID, _ := domain.CreateParameter(paramGroupID, parameterName)
				log.Println("Parameter ID: " + parameterID)
				domain.CreateParamSpecWithModel(parameterID, upper, target, lower)

				log.Println(" ")
				r, _ := interactive.PromptAddMoreParameter()
				if r == "n" {
					break
				}
			}
		}

		log.Println(" ")
		log.Println("===========================================")
		log.Println("8. Create a device.(8/9) ")

		c, _ = interactive.PromptAddDevice()
		if c == "y" {
			for {

				deviceName, _ := interactive.PromptDeviceName()
				ipAddress, _ := interactive.PromptDeviceAddress()
				macAddress, _ := interactive.PromptDeviceMacAddress()
				modelNum, _ := interactive.PromptDeviceModelNumber()
				serialNum, _ := interactive.PromptDeviceSerialNumber()

				domain.CreateDevice(tenantID, ipAddress, macAddress, modelNum, serialNum, deviceName)

				log.Println(" ")
				r, _ := interactive.PromptAddMoreDevice()
				if r == "n" {
					break
				}
			}
		}

		log.Println(" ")
		log.Println("===========================================")
		log.Println("9. Create a asset with sensor.(9/9) ")

		for {
			nestID, _, _ := interactive.SelectNest(tenantID)
			templateID, _, _ := interactive.SelectTemplate(tenantID)
			paramGroupID, _, _ := interactive.SelectParamGroup(tenantID)
			catalogID, _, _ := interactive.SelectCatalogs(rootCatalogID)
			assetName, _ := interactive.PromptAsset()

			phyAsset, assetID, _ := domain.CreateAsset(tenantID, templateID, catalogID, nestID, assetName)
			log.Println("Physical Asset Name: " + phyAsset)

			domain.CreateParamSpecWithAsset(assetID, tenantID, paramGroupID)

			deviceID, _, _ := interactive.SelectDevice(tenantID)
			log.Println("device is selected: " + deviceID)
			sensorName, _ := interactive.PromptSensor(assetID)
			sensorDuration, _ := interactive.PromptSensorDuration()
			domain.CreateSensor(assetID, deviceID, paramGroupID, sensorDuration, sensorName)
			log.Println("Physical Sensor Name: " + sensorName)
			log.Println(" ")

			r, _ := interactive.PromptAddMoreAsset()
			if r == "n" {
				break
			}
		}

		log.Println(" ")
		log.Println("===========================================")
	},
}

func init() {
	tenantCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
