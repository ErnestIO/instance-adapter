/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"encoding/json"
)

type disk struct {
	ID   int `json:"id"`
	Size int `json:"size"`
}

type builderEvent struct {
	Uuid                  string   `json:"_uuid"`
	BatchID               string   `json:"_batch_id"`
	Type                  string   `json:"type"`
	Service               string   `json:"service"`
	Name                  string   `json:"name"`
	CPU                   int      `json:"cpus"`
	RAM                   int      `json:"ram"`
	IP                    string   `json:"ip"`
	PublicIP              string   `json:"public_ip"`
	Catalog               string   `json:"reference_catalog"`
	Image                 string   `json:"reference_image"`
	Disks                 []disk   `json:"disks"`
	AssignElasticIP       bool     `json:"assign_elastic_ip"`
	InstanceAWSID         string   `json:"instance_aws_id"`
	RouterName            string   `json:"router_name"`
	RouterType            string   `json:"router_type"`
	RouterIP              string   `json:"router_ip"`
	ClientName            string   `json:"client_name"`
	DatacenterName        string   `json:"datacenter_name"`
	DatacenterPassword    string   `json:"datacenter_password"`
	DatacenterRegion      string   `json:"datacenter_region"`
	DatacenterType        string   `json:"datacenter_type"`
	DatacenterUsername    string   `json:"datacenter_username"`
	DatacenterAccessToken string   `json:"datacenter_token"`
	DatacenterAccessKey   string   `json:"datacenter_secret"`
	NetworkName           string   `json:"network_name"`
	NetworkAWSID          string   `json:"network_aws_id"`
	KeyPair               string   `json:"key_pair"`
	SecurityGroupAWSIDs   []string `json:"security_group_aws_ids"`
	VCloudURL             string   `json:"vcloud_url"`
	Status                string   `json:"status"`
	ErrorCode             string   `json:"error_code"`
	ErrorMessage          string   `json:"error_message"`
}

type instanceResource struct {
	CPU     int    `json:"cpus"`
	RAM     int    `json:"ram"`
	IP      string `json:"ip"`
	Catalog string `json:"reference_catalog"`
	Image   string `json:"reference_image"`
	Disks   []disk `json:"disks"`
}

type vcloudEvent struct {
	Uuid               string           `json:"_uuid"`
	BatchID            string           `json:"_batch_id"`
	Type               string           `json:"_type"`
	Service            string           `json:"service_id"`
	InstanceName       string           `json:"instance_name"`
	InstanceType       string           `json:"instance_type"`
	Resource           instanceResource `json:"instance_resource"`
	RouterName         string           `json:"router_name"`
	RouterType         string           `json:"router_type"`
	RouterIP           string           `json:"router_ip"`
	ClientName         string           `json:"client_name"`
	DatacenterName     string           `json:"datacenter_name"`
	DatacenterPassword string           `json:"datacenter_password"`
	DatacenterRegion   string           `json:"datacenter_region"`
	DatacenterType     string           `json:"datacenter_type"`
	DatacenterUsername string           `json:"datacenter_username"`
	NetworkName        string           `json:"network_name"`
	VCloudURL          string           `json:"vcloud_url"`
	Status             string           `json:"status"`
	ErrorCode          string           `json:"error_code"`
	ErrorMessage       string           `json:"error_message"`
}

type awsEvent struct {
	Uuid                    string   `json:"_uuid"`
	BatchID                 string   `json:"_batch_id"`
	Type                    string   `json:"_type"`
	DatacenterRegion        string   `json:"datacenter_region,omitempty"`
	DatacenterAccessToken   string   `json:"datacenter_access_token"`
	DatacenterAccessKey     string   `json:"datacenter_access_key"`
	DatacenterVpcID         string   `json:"datacenter_vpc_id,omitempty"`
	NetworkAWSID            string   `json:"network_aws_id"`
	SecurityGroupAWSIDs     []string `json:"security_group_aws_ids"`
	InstanceName            string   `json:"instance_name"`
	InstanceImage           string   `json:"instance_image"`
	InstanceType            string   `json:"instance_type"`
	InstanceIP              string   `json:"instance_ip"`
	InstanceElasticIP       string   `json:"instance_elastic_ip"`
	InstanceAWSID           string   `json:"instance_aws_id"`
	InstanceKeyPair         string   `json:"instance_key_pair"`
	InstanceAssignElasticIP bool     `json:"instance_assign_elastic_ip"`
	ErrorMessage            string   `json:"error"`
}

type Translator struct{}

func (t Translator) BuilderToConnector(j []byte) []byte {
	var input builderEvent
	var output []byte
	json.Unmarshal(j, &input)

	switch input.DatacenterType {
	case "vcloud", "vcloud-fake", "fake":
		output = t.builderToVCloudConnector(input)
	case "aws", "aws-fake":
		output = t.builderToAwsConnector(input)
	}

	return output
}

func (t Translator) builderToVCloudConnector(input builderEvent) []byte {
	var output vcloudEvent

	resource := instanceResource{
		CPU:     input.CPU,
		RAM:     input.RAM,
		IP:      input.IP,
		Catalog: input.Catalog,
		Image:   input.Image,
		Disks:   input.Disks,
	}

	output.Uuid = input.Uuid
	output.BatchID = input.BatchID
	output.Type = input.DatacenterType
	output.Service = input.Service
	output.InstanceName = input.Name
	output.InstanceType = input.DatacenterType
	output.Resource = resource
	output.RouterIP = input.RouterIP
	output.RouterName = input.RouterName
	output.RouterType = input.RouterType
	output.NetworkName = input.NetworkName
	output.ClientName = input.ClientName
	output.DatacenterName = input.DatacenterName
	output.DatacenterRegion = input.DatacenterRegion
	output.DatacenterUsername = input.DatacenterUsername
	output.DatacenterPassword = input.DatacenterPassword
	output.DatacenterType = input.DatacenterType
	output.VCloudURL = input.VCloudURL
	output.Status = input.Status
	output.ErrorCode = input.ErrorCode
	output.ErrorMessage = input.ErrorMessage
	body, _ := json.Marshal(output)

	return body
}

func (t Translator) builderToAwsConnector(input builderEvent) []byte {
	var output awsEvent

	output.Uuid = input.Uuid
	output.BatchID = input.BatchID
	output.Type = input.DatacenterType
	output.DatacenterRegion = input.DatacenterRegion
	output.DatacenterAccessToken = input.DatacenterAccessToken
	output.DatacenterAccessKey = input.DatacenterAccessKey
	output.DatacenterVpcID = input.DatacenterName
	output.NetworkAWSID = input.NetworkAWSID
	output.SecurityGroupAWSIDs = input.SecurityGroupAWSIDs
	output.InstanceName = input.Name
	output.InstanceImage = input.Image
	output.InstanceType = input.Type
	output.InstanceIP = input.IP
	output.InstanceKeyPair = input.KeyPair
	output.InstanceAssignElasticIP = input.AssignElasticIP
	output.InstanceAWSID = input.InstanceAWSID

	body, _ := json.Marshal(output)

	return body
}

func (t Translator) ConnectorToBuilder(j []byte) []byte {
	var output []byte
	var input map[string]interface{}

	dec := json.NewDecoder(bytes.NewReader(j))
	dec.Decode(&input)

	switch input["_type"] {
	case "vcloud", "vcloud-fake", "fake":
		output = t.vcloudConnectorToBuilder(j)
	case "aws", "aws-fake":
		output = t.awsConnectorToBuilder(j)
	}

	return output
}

func (t Translator) vcloudConnectorToBuilder(j []byte) []byte {
	var input vcloudEvent
	var output builderEvent
	json.Unmarshal(j, &input)

	output.Uuid = input.Uuid
	output.BatchID = input.BatchID
	output.Type = input.DatacenterType
	output.Service = input.Service
	output.Name = input.InstanceName
	output.CPU = input.Resource.CPU
	output.RAM = input.Resource.RAM
	output.IP = input.Resource.IP
	output.Catalog = input.Resource.Catalog
	output.Image = input.Resource.Image
	output.Disks = input.Resource.Disks
	output.RouterIP = input.RouterIP
	output.RouterName = input.RouterName
	output.RouterType = input.RouterType
	output.ClientName = input.ClientName
	output.NetworkName = input.NetworkName
	output.DatacenterName = input.DatacenterName
	output.DatacenterUsername = input.DatacenterUsername
	output.DatacenterPassword = input.DatacenterPassword
	output.DatacenterType = input.DatacenterType
	output.VCloudURL = input.VCloudURL
	output.Status = input.Status
	output.ErrorCode = input.ErrorCode
	output.ErrorMessage = input.ErrorMessage

	body, _ := json.Marshal(output)

	return body
}

func (t Translator) awsConnectorToBuilder(j []byte) []byte {
	var input awsEvent
	var output builderEvent
	json.Unmarshal(j, &input)

	output.Uuid = input.Uuid
	output.BatchID = input.BatchID
	output.Type = input.Type
	output.DatacenterRegion = input.DatacenterRegion
	output.DatacenterAccessToken = input.DatacenterAccessToken
	output.DatacenterAccessKey = input.DatacenterAccessKey
	output.DatacenterName = input.DatacenterVpcID
	output.SecurityGroupAWSIDs = input.SecurityGroupAWSIDs
	output.Name = input.InstanceName
	output.Image = input.InstanceImage
	output.Type = input.InstanceType
	output.IP = input.InstanceIP
	output.AssignElasticIP = input.InstanceAssignElasticIP
	output.PublicIP = input.InstanceElasticIP
	output.InstanceAWSID = input.InstanceAWSID

	if input.ErrorMessage != "" {
		output.Status = "errored"
		output.ErrorCode = "0"
		output.ErrorMessage = input.ErrorMessage
	}

	body, _ := json.Marshal(output)

	return body
}
