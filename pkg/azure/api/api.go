// Copyright 2020 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/cilium/cilium/pkg/api/helpers"
	"github.com/cilium/cilium/pkg/azure/types"
	"github.com/cilium/cilium/pkg/ipam"
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/cilium/cilium/pkg/logging"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/spanstat"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	userAgent = "cilium"
)

var log = logging.DefaultLogger.WithField(logfields.LogSubsys, "azure-api")

// Client represents an EC2 API client
type Client struct {
	resourceGroup   string
	interfaces      network.InterfacesClient
	virtualnetworks network.VirtualNetworksClient
	vmscalesets     compute.VirtualMachineScaleSetsClient
	limiter         *helpers.ApiLimiter
	metricsAPI      MetricsAPI
}

type MetricsAPI interface {
	ObserveAPICall(call, status string, duration float64)
	ObserveRateLimit(operation string, duration time.Duration)
}

// NewClient returns a new EC2 client
func NewClient(subscriptionID, resourceGroup string, metrics MetricsAPI, rateLimit float64, burst int) (*Client, error) {
	c := &Client{
		resourceGroup:   resourceGroup,
		interfaces:      network.NewInterfacesClient(subscriptionID),
		virtualnetworks: network.NewVirtualNetworksClient(subscriptionID),
		vmscalesets:     compute.NewVirtualMachineScaleSetsClient(subscriptionID),
		metricsAPI:      metrics,
		limiter:         helpers.NewApiLimiter(metrics, rateLimit, burst),
	}

	// Authorizer based on environment variables
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}

	c.interfaces.Authorizer = authorizer
	c.interfaces.AddToUserAgent(userAgent)
	c.virtualnetworks.Authorizer = authorizer
	c.virtualnetworks.AddToUserAgent(userAgent)
	c.vmscalesets.Authorizer = authorizer
	c.vmscalesets.AddToUserAgent(userAgent)

	return c, nil
}

// deriveStatus returns a status string
func deriveStatus(err error) string {
	if err != nil {
		return "Failed"
	}

	return "OK"
}

// describeNetworkInterfaces lists all Azure Interfaces
func (c *Client) describeNetworkInterfaces(ctx context.Context) ([]network.Interface, error) {
	var networkInterfaces []network.Interface

	c.limiter.Limit(ctx, "VirtualMachineScaleSets.ListAll")
	sinceStart := spanstat.Start()
	result, err := c.vmscalesets.ListComplete(ctx, c.resourceGroup)
	c.metricsAPI.ObserveAPICall("VirtualMachineScaleSets.ListAll", deriveStatus(err), sinceStart.Seconds())
	if err != nil {
		return nil, err
	}

	for result.NotDone() {
		if err != nil {
			return nil, err
		}

		scaleset := result.Value()
		err = result.Next()

		log.Infof("scaleset %+v", scaleset)

		if scaleset.Name == nil {
			continue
		}

		c.limiter.Limit(ctx, "Interfaces.ListAll")
		sinceStart := spanstat.Start()
		result2, err2 := c.interfaces.ListVirtualMachineScaleSetNetworkInterfacesComplete(ctx, c.resourceGroup, *scaleset.Name)
		c.metricsAPI.ObserveAPICall("Interfaces.ListVirtualMachineScaleSetNetworkInterfacesComplete", deriveStatus(err2), sinceStart.Seconds())
		if err2 != nil {
			return nil, err2
		}

		for result2.NotDone() {
			if err2 != nil {
				return nil, err2
			}

			log.Infof("interface %+v", result2.Value())

			networkInterfaces = append(networkInterfaces, result2.Value())
			err2 = result2.Next()
		}
	}

	return networkInterfaces, nil
}

// parseInterfaces parses a network.Interface as returned by the Azure API
// converts it into a v2.AzureInterface
func parseInterface(iface *network.Interface) (instanceID string, i *v2.AzureInterface) {
	i = &v2.AzureInterface{}

	if iface.VirtualMachine != nil && iface.VirtualMachine.ID != nil {
		instanceID = strings.ToLower(*iface.VirtualMachine.ID)
	}

	if iface.MacAddress != nil {
		// Azure API reports MAC addresses as AA-BB-CC-DD-EE-FF
		i.MAC = strings.ReplaceAll(*iface.MacAddress, "-", ":")
	}

	if iface.ID != nil {
		i.ID = *iface.ID
	}

	if iface.NetworkSecurityGroup != nil {
		if iface.NetworkSecurityGroup.ID != nil {
			i.SecurityGroup = *iface.NetworkSecurityGroup.ID
		}
	}

	if iface.IPConfigurations != nil {
		for _, ip := range *iface.IPConfigurations {
			if ip.PrivateIPAddress != nil {
				addr := v2.AzureAddress{
					IP:    *ip.PrivateIPAddress,
					State: strings.ToLower(string(ip.ProvisioningState)),
				}

				if ip.Subnet != nil {
					addr.Subnet = *ip.Subnet.ID
				}

				i.Addresses = append(i.Addresses, addr)
			}
		}
	}

	return
}

// GetInstances returns the list of all instances including their ENIs as
// instanceMap
func (c *Client) GetInstances(ctx context.Context) (types.InstanceMap, error) {
	instances := types.InstanceMap{}

	networkInterfaces, err := c.describeNetworkInterfaces(ctx)
	if err != nil {
		return nil, err
	}

	log.Infof("Successfully called getInstances %+v", networkInterfaces)

	for _, iface := range networkInterfaces {
		if id, azureInterface := parseInterface(&iface); id != "" {
			log.Infof("Interface: ID=%s, %+v", id, azureInterface)
			instances.Update(id, azureInterface)
		}
	}

	return instances, nil
}

// describeVpcs lists all VPCs
func (c *Client) describeVpcs(ctx context.Context) ([]network.VirtualNetwork, error) {
	var vpcs []network.VirtualNetwork

	c.limiter.Limit(ctx, "VirtualNetworks.List")

	sinceStart := spanstat.Start()
	result, err := c.virtualnetworks.ListComplete(ctx, c.resourceGroup)
	c.metricsAPI.ObserveAPICall("Interfaces.ListAll", deriveStatus(err), sinceStart.Seconds())
	if err != nil {
		return nil, err
	}

	for result.NotDone() {
		if err != nil {
			return nil, err
		}

		vpcs = append(vpcs, result.Value())
		err = result.Next()
	}

	return vpcs, nil
}

func parseSubnet(subnet *network.Subnet) (s *ipam.Subnet) {
	s = &ipam.Subnet{ID: *subnet.ID}
	if subnet.Name != nil {
		s.Name = *subnet.Name
	}

	if subnet.AddressPrefix != nil {
		s.CIDR = *subnet.AddressPrefix
	}

	return
}

// GetVpcsAndSubnets retrieves and returns all Vpcs
func (c *Client) GetVpcsAndSubnets(ctx context.Context) (ipam.VirtualNetworkMap, ipam.SubnetMap, error) {
	vpcs := ipam.VirtualNetworkMap{}
	subnets := ipam.SubnetMap{}

	vpcList, err := c.describeVpcs(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, v := range vpcList {
		if v.ID == nil {
			continue
		}

		vpc := &ipam.VirtualNetwork{ID: *v.ID}
		log.Infof("VPC: %+v", vpc)
		vpcs[vpc.ID] = vpc

		if v.Subnets != nil {
			for _, subnet := range *v.Subnets {
				if subnet.ID == nil {
					continue
				}

				log.Infof("Subnet: %+v", subnet)
				subnets[*subnet.ID] = parseSubnet(&subnet)

			}
		}
	}

	return vpcs, subnets, nil
}

func (c *Client) AssignPrivateIpAddresses(ctx context.Context, subnetID, interfaceID string, ips []net.IP) error {
	var ipConfigurations []network.InterfaceIPConfiguration

	for _, ip := range ips {
		ipString := ip.String()
		ipConfigurations = append(ipConfigurations,
			network.InterfaceIPConfiguration{
				InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
					PrivateIPAddress: &ipString,
					Subnet:           &network.Subnet{ID: &subnetID},
				},
			},
		)

	}

	ifaceParams := network.Interface{
		ID: &interfaceID,
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &ipConfigurations,
		},
	}

	result, err := c.interfaces.CreateOrUpdate(ctx, c.resourceGroup, interfaceID, ifaceParams)
	if err == nil {
		log.Infof("Updated interface: %+v", result)
	}

	return err
}
