/*
   Copyright 2021 VMware, Inc.
   SPDX-License-Identifier: Apache-2.0
*/

package vcdclient

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	swagger "github.com/vmware/cloud-provider-for-cloud-director/pkg/vcdswaggerclient"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
)

const BusyRetries = 5

func TestCacheGatewayDetails(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	err = gm.CacheGatewayDetails(ctx)
	assert.NoError(t, err, "Unable to cache gateway details")

	require.NotNil(t, gm.GatewayRef, "Gateway reference should not be nil")
	assert.NotEmpty(t, gm.GatewayRef.Name, "Gateway Name should not be empty")
	assert.NotEmpty(t, gm.GatewayRef.Id, "Gateway Id should not be empty")

	// Missing network name should be reported
	vcdClient, err = getTestVCDClient(cloudConfig, map[string]interface{}{
		"network": "",
	})
	assert.Error(t, err, "Should get error for unknown network")
	assert.Nil(t, vcdClient, "Client should be nil when erroring out")

	return
}

func TestDNATRuleCRUDE(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	dnatRuleName := fmt.Sprintf("test-dnat-rule-%s", uuid.New().String())
	err = gm.createDNATRule(ctx, dnatRuleName, "1.2.3.4", "1.2.3.5", 80, 36123)
	assert.NoError(t, err, "Unable to create dnat rule")

	// repeated creation should not fail
	err = gm.createDNATRule(ctx, dnatRuleName, "1.2.3.4", "1.2.3.5", 80, 36123)
	assert.NoError(t, err, "Unable to create dnat rule for the second time")

	natRuleRef, err := gm.getNATRuleRef(ctx, dnatRuleName)
	assert.NoError(t, err, "Unable to get dnat rule")
	require.NotNil(t, natRuleRef, "Nat Rule reference should not be nil")
	assert.Equal(t, dnatRuleName, natRuleRef.Name, "Nat Rule name should match")
	assert.NotEmpty(t, natRuleRef.ID, "Nat Rule ID should not be empty")

	err = gm.updateDNATRule(ctx, dnatRuleName, "2.3.4.5", "2.3.4.5", 8080)
	assert.NoError(t, err, "Unable to update dnat rule")

	err = gm.updateDNATRule(ctx, dnatRuleName, "2.3.4.5", "2.3.4.5", 8080)
	assert.NoError(t, err, "repeated updates to dnat rule should not fail")

	err = gm.deleteDNATRule(ctx, dnatRuleName, true)
	assert.NoError(t, err, "Unable to delete dnat rule")

	err = gm.deleteDNATRule(ctx, dnatRuleName, true)
	assert.Error(t, err, "Should fail when deleting non-existing dnat rule")

	err = gm.deleteDNATRule(ctx, dnatRuleName, false)
	assert.NoError(t, err, "Should not fail when deleting non-existing dnat rule")

	natRuleRef, err = gm.getNATRuleRef(ctx, dnatRuleName)
	assert.NoError(t, err, "Get should not fail when nat rule is absent")
	assert.Nil(t, natRuleRef, "Deleted nat rule reference should be nil")

	return
}

func TestLBPoolCRUDE(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	lbPoolName := fmt.Sprintf("test-lb-pool-%s", uuid.New().String())
	lbPoolRef, err := gm.createLoadBalancerPool(ctx, lbPoolName, []string{"1.2.3.4", "1.2.3.5"}, 31234)
	assert.NoError(t, err, "Unable to create lb pool")
	require.NotNil(t, lbPoolRef, "LB Pool reference should not be nil")
	assert.Equal(t, lbPoolName, lbPoolRef.Name, "LB Pool name should match")

	// repeated creation should not fail
	lbPoolRef, err = gm.createLoadBalancerPool(ctx, lbPoolName, []string{"1.2.3.4", "1.2.3.5"}, 31234)
	assert.NoError(t, err, "Unable to create lb pool for the second time")
	require.NotNil(t, lbPoolRef, "LB Pool reference should not be nil")
	assert.Equal(t, lbPoolName, lbPoolRef.Name, "LB Pool name should match")

	lbPoolRefObtained, err := gm.getLoadBalancerPool(ctx, lbPoolName)
	assert.NoError(t, err, "Unable to get lbPool ref")
	require.NotNil(t, lbPoolRefObtained, "LB Pool reference should not be nil")
	assert.Equal(t, lbPoolRefObtained.Name, lbPoolRef.Name, "LB Pool name should match")
	assert.NotEmpty(t, lbPoolRefObtained.Id, "LB Pool ID should not be empty")

	updatedIps := []string{"5.5.5.5"}
	lbPoolRefUpdated, err := gm.updateLoadBalancerPool(ctx, lbPoolName, updatedIps, 55555)
	assert.NoError(t, err, "No lbPool ref for updated lbPool")
	require.NotNil(t, lbPoolRefUpdated, "LB Pool reference should not be nil")
	assert.Equal(t, lbPoolRefUpdated.Name, lbPoolRef.Name, "LB Pool name should match")
	assert.NotEmpty(t, lbPoolRefUpdated.Id, "LB Pool ID should not be empty")

	// repeated update should work
	lbPoolRefUpdated, err = gm.updateLoadBalancerPool(ctx, lbPoolName, updatedIps, 55555)
	assert.NoError(t, err, "There should be no error on repeated LB pool update")
	require.NotNil(t, lbPoolRefUpdated, "LB Pool reference should not be nil on repeated update")
	assert.Equal(t, lbPoolRefUpdated.Name, lbPoolRef.Name, "LB Pool name should match on repeated update")
	assert.NotEmpty(t, lbPoolRefUpdated.Id, "LB Pool ID should not be empty on repeated update")

	lbPoolSummaryUpdated, err := gm.getLoadBalancerPoolSummary(ctx, lbPoolName)
	assert.NoError(t, err, "No LB Pool summary for updated pool.")
	require.NotNil(t, lbPoolSummaryUpdated, "LB Pool summary reference should not be nil")
	assert.Equal(t, lbPoolSummaryUpdated.MemberCount, int32(len(updatedIps)), "LB Pool should have updated size %d", len(updatedIps))

	err = gm.deleteLoadBalancerPool(ctx, lbPoolName, true)
	assert.NoError(t, err, "Unable to delete LB Pool")

	err = gm.deleteLoadBalancerPool(ctx, lbPoolName, true)
	assert.Error(t, err, "Should fail when deleting non-existing lb pool")

	err = gm.deleteLoadBalancerPool(ctx, lbPoolName, false)
	assert.NoError(t, err, "Should not fail when deleting non-existing lb pool")

	lbPoolRef, err = gm.getLoadBalancerPool(ctx, lbPoolName)
	assert.NoError(t, err, "Get should not fail when lb pool is absent")
	assert.Nil(t, lbPoolRef, "Deleted lb pool reference should be nil")

	lbPoolRef, err = gm.updateLoadBalancerPool(ctx, lbPoolName, updatedIps, 55555)
	assert.Error(t, err, "Updating deleted lb pool should fail")
	assert.Nil(t, lbPoolRef, "Deleted lb pool reference should be nil")

	return
}

func TestGetLoadBalancerSEG(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	segRef, err := gm.getLoadBalancerSEG(ctx)
	assert.NoError(t, err, "Unable to get ServiceEngineGroup ref")
	require.NotNil(t, segRef, "ServiceEngineGroup reference should not be nil")
	assert.NotEmpty(t, segRef.Name, "ServiceEngineGroup Name should not be empty")
	assert.NotEmpty(t, segRef.Id, "ServiceEngineGroup ID should not be empty")

	return
}

func TestGetUnusedGatewayIP(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
		"subnet":       "",
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	validSubnet := cloudConfig.LB.VIPSubnet
	externalIP, err := gm.getUnusedExternalIPAddress(ctx, validSubnet)
	assert.NoError(t, err, "should not get an error for this range")
	assert.NotEmpty(t, externalIP, "should get a valid IP address in the range [%s]", validSubnet)

	invalidSubnet := "1.1.1.1/24"
	externalIP, err = gm.getUnusedExternalIPAddress(ctx, invalidSubnet)
	assert.Error(t, err, "should get an error for this range")
	assert.Empty(t, externalIP, "should not get a valid IP address in the range [%s]", invalidSubnet)

	everythingAllowedSubnet := ""
	externalIP, err = gm.getUnusedExternalIPAddress(ctx, everythingAllowedSubnet)
	assert.NoError(t, err, "should not get an error for this range")
	assert.NotEmpty(t, externalIP, "should get a valid IP address in the empty range")

	return
}

func TestVirtualServiceHttpCRUDE(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	lbPoolName := fmt.Sprintf("test-lb-pool-%s", uuid.New().String())
	lbPoolRef, err := gm.createLoadBalancerPool(ctx, lbPoolName, []string{"1.2.3.4", "1.2.3.5"}, 31234)
	assert.NoError(t, err, "Unable to create lb pool")

	segRef, err := gm.getLoadBalancerSEG(ctx)
	assert.NoError(t, err, "Unable to get ServiceEngineGroup ref")

	virtualServiceName := fmt.Sprintf("test-virtual-service-%s", uuid.New().String())
	externalIP := "10.11.12.13"
	internalIP := "2.3.4.5"
	var vsRef *swagger.EntityReference
	for i := 0; i < BusyRetries; i ++ {
		vsRef, err = gm.createVirtualService(ctx, virtualServiceName, lbPoolRef, segRef,
			internalIP, externalIP, "HTTP", 80, false, "", cloudConfig.ClusterID)
		if err != nil {
			if _, ok := err.(*VirtualServicePendingError); !ok {
				break
			}
		}
	}
	assert.NoError(t, err, "Unable to create virtual service")
	require.NotNil(t, vsRef, "VirtualServiceRef should not be nil")
	assert.Equal(t, virtualServiceName, vsRef.Name, "Virtual Service name should match")

	vsRefObtained, err := gm.getVirtualService(ctx, virtualServiceName)
	assert.NoError(t, err, "Unable to get virtual service ref")
	require.NotNil(t, vsRefObtained, "Virtual service reference should not be nil")
	assert.Equal(t, vsRefObtained.Name, vsRef.Name, "Virtual service reference name should match")
	assert.NotEmpty(t, vsRefObtained.Id, "Virtual service ID should not be empty")

	rm := NewRDEManager(vcdClient, cloudConfig.ClusterID)
	rdeVips, _, _, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Unable to get RDE vips after virtual service creation")
	assert.True(t, foundStringInSlice(externalIP, rdeVips), "external ip should be found in rde vips")
	assert.False(t, foundStringInSlice(internalIP, rdeVips), "internal ip should not be in rde vips")

	// repeated creation should not fail
	vsRef, err = gm.createVirtualService(ctx, virtualServiceName, lbPoolRef, segRef,
		internalIP, externalIP, "HTTP", 80, false, "", cloudConfig.ClusterID)
	assert.NoError(t, err, "Unable to create virtual service for the second time")
	require.NotNil(t, vsRef, "VirtualServiceRef should not be nil")
	assert.Equal(t, virtualServiceName, vsRef.Name, "Virtual Service name should match")

	err = gm.updateVirtualServicePort(ctx, virtualServiceName, 8080)
	assert.NoError(t, err, "Unable to update external port")

	// repeated update should not fail
	err = gm.updateVirtualServicePort(ctx, virtualServiceName, 8080)
	assert.NoError(t, err, "Repeated update to external port should not fail")

	err = gm.updateVirtualServicePort(ctx, virtualServiceName+"-invalid", 8080)
	assert.Error(t, err, "Update virtual service on a non-existent virtual service should fail")

	err = gm.deleteVirtualService(ctx, virtualServiceName, true, externalIP, cloudConfig.ClusterID)
	assert.NoError(t, err, "Unable to delete Virtual Service")

	rdeVips, _, _, err = rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Unable to get vips from RDE after virtual service deletion")
	assert.False(t, foundStringInSlice(externalIP, rdeVips), "external ip should not be found in RDE vips")

	err = gm.deleteVirtualService(ctx, virtualServiceName, true, externalIP, cloudConfig.ClusterID)
	assert.Error(t, err, "Should fail when deleting non-existing Virtual Service")

	err = gm.deleteVirtualService(ctx, virtualServiceName, false, externalIP, cloudConfig.ClusterID)
	assert.NoError(t, err, "Should not fail when deleting non-existing Virtual Service")

	vsRefObtained, err = gm.getVirtualService(ctx, virtualServiceName)
	assert.NoError(t, err, "Get should not fail when Virtual Service is absent")
	assert.Nil(t, vsRefObtained, "Deleted Virtual Service reference should be nil")

	err = gm.deleteLoadBalancerPool(ctx, lbPoolName, true)
	assert.NoError(t, err, "Should not fail when deleting lb pool")

	return
}

func foundStringInSlice(find string, slice []string) bool {
	for _, currElement := range slice {
		if currElement == find {
			return true
		}
	}
	return false
}

func TestVirtualServiceHttpsCRUDE(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	lbPoolName := fmt.Sprintf("test-lb-pool-%s", uuid.New().String())
	lbPoolRef, err := gm.createLoadBalancerPool(ctx, lbPoolName, []string{"1.2.3.4", "1.2.3.5"}, 31234)
	assert.NoError(t, err, "Unable to create lb pool")

	segRef, err := gm.getLoadBalancerSEG(ctx)
	assert.NoError(t, err, "Unable to get ServiceEngineGroup ref")

	externalIP := "11.12.13.14"
	internalIP := "3.4.5.6"
	virtualServiceName := fmt.Sprintf("test-virtual-service-https-%s", uuid.New().String())
	certName := cloudConfig.LB.CertificateAlias
	if certName == "" {
		certName = fmt.Sprintf("%s-cert", cloudConfig.ClusterID)
	}
	vsRef, err := gm.createVirtualService(ctx, virtualServiceName, lbPoolRef, segRef,
		internalIP, externalIP, "HTTPS", 443, true, certName, cloudConfig.ClusterID)
	assert.NoError(t, err, "Unable to create virtual service")
	require.NotNil(t, vsRef, "VirtualServiceRef should not be nil")
	assert.Equal(t, virtualServiceName, vsRef.Name, "Virtual Service name should match")

	vsRefObtained, err := gm.getVirtualService(ctx, virtualServiceName)
	assert.NoError(t, err, "Unable to get virtual service ref")
	require.NotNil(t, vsRefObtained, "Virtual service reference should not be nil")
	assert.Equal(t, vsRefObtained.Name, vsRef.Name, "Virtual service reference name should match")
	assert.NotEmpty(t, vsRefObtained.Id, "Virtual service ID should not be empty")

	rm := NewRDEManager(vcdClient, cloudConfig.ClusterID)
	rdeVips, _, _, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Unable to get RDE vips after virtual service creation")
	assert.True(t, foundStringInSlice(externalIP, rdeVips), "external ip should be found in rde vips")
	assert.False(t, foundStringInSlice(internalIP, rdeVips), "internal ip should not be found in rde vips")

	// repeated creation should not fail
	vsRef, err = gm.createVirtualService(ctx, virtualServiceName, lbPoolRef, segRef,
		internalIP, externalIP, "HTTPS", 443, true, certName, cloudConfig.ClusterID)
	assert.NoError(t, err, "Unable to create virtual service for the second time")
	require.NotNil(t, vsRef, "VirtualServiceRef should not be nil")
	assert.Equal(t, virtualServiceName, vsRef.Name, "Virtual Service name should match")

	// update and delete calls might error out if virtual services are busy. Retry if the error is caused by the busy virtual services
	err = gm.updateVirtualServicePort(ctx, virtualServiceName, 8443)
	assert.NoError(t, err, "Unable to update external port")

	// repeated update should not fail
	err = gm.updateVirtualServicePort(ctx, virtualServiceName, 8443)
	assert.NoError(t, err, "Repeated update to external port should not fail")

	// update of invalid virtual service
	err = gm.updateVirtualServicePort(ctx, virtualServiceName+"-invalid\n", 8443)
	assert.Error(t, err, "Update virtual service on a non-existent virtual service should fail")

	err = gm.deleteVirtualService(ctx, virtualServiceName, true, externalIP, cloudConfig.ClusterID)
	assert.NoError(t, err, "Unable to delete Virtual Service")

	rdeVips, _, _, err = rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Unable to get vips from RDE after virtual service deletion")
	assert.False(t, foundStringInSlice(externalIP, rdeVips), "external ip should not be found in RDE vips")

	err = gm.deleteVirtualService(ctx, virtualServiceName, true, externalIP, cloudConfig.ClusterID)
	assert.Error(t, err, "Should fail when deleting non-existing Virtual Service")

	err = gm.deleteVirtualService(ctx, virtualServiceName, false, externalIP, cloudConfig.ClusterID)
	assert.NoError(t, err, "Should not fail when deleting non-existing Virtual Service")

	vsRefObtained, err = gm.getVirtualService(ctx, virtualServiceName)
	assert.NoError(t, err, "Get should not fail when Virtual Service is absent")
	assert.Nil(t, vsRefObtained, "Deleted Virtual Service reference should be nil")

	err = gm.deleteLoadBalancerPool(ctx, lbPoolName, true)
	assert.NoError(t, err, "Should not fail when deleting lb pool")

	return
}

func TestLoadBalancerCRUDE(t *testing.T) {

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":         authDetails.Username,
		"secret":       authDetails.Password,
		"userOrg":      authDetails.UserOrg,
		"getVdcClient": true,
	})
	assert.NoError(t, err, "Unable to get VCD client")
	require.NotNil(t, vcdClient, "VCD Client should not be nil")

	ctx := context.Background()

	gm, err := NewGatewayManager(ctx, vcdClient, cloudConfig.LB.VDCNetwork, cloudConfig.LB.VIPSubnet)
	assert.NoError(t, err, "gateway manager should be created without error")

	virtualServiceNamePrefix := fmt.Sprintf("test-virtual-service-https-%s", uuid.New().String())
	lbPoolNamePrefix := fmt.Sprintf("test-lb-pool-%s", uuid.New().String())
	certName := cloudConfig.LB.CertificateAlias
	if certName == "" {
		certName = fmt.Sprintf("%s-cert", cloudConfig.ClusterID)
	}
	portDetailsList := []PortDetails{
		{
			PortSuffix:   `http`,
			ExternalPort: 80,
			InternalPort: 31234,
			Protocol:     "HTTP",
			UseSSL:       false,
		},
		{
			PortSuffix:   `https`,
			ExternalPort: 443,
			InternalPort: 31235,
			Protocol:     "HTTPS",
			UseSSL:       true,
			CertAlias:    certName,
		},
	}
	freeIP := ""
	freeIP, err = gm.CreateLoadBalancer(ctx, virtualServiceNamePrefix,
		lbPoolNamePrefix, []string{"1.2.3.4", "1.2.3.5"}, portDetailsList, cloudConfig.LB.OneArm, cloudConfig.ClusterID)
	assert.NoError(t, err, "Load Balancer should be created")
	assert.NotEmpty(t, freeIP, "There should be a non-empty IP returned")

	virtualServiceNameHttp := fmt.Sprintf("%s-http", virtualServiceNamePrefix)
	freeIPObtained, err := gm.GetLoadBalancer(ctx, virtualServiceNameHttp, cloudConfig.LB.OneArm)
	assert.NoError(t, err, "Load Balancer should be found")
	assert.Equal(t, freeIP, freeIPObtained, "The IPs should match")

	virtualServiceNameHttps := fmt.Sprintf("%s-https", virtualServiceNamePrefix)
	freeIPObtained, err = gm.GetLoadBalancer(ctx, virtualServiceNameHttps, cloudConfig.LB.OneArm)
	assert.NoError(t, err, "Load Balancer should be found")
	assert.Equal(t, freeIP, freeIPObtained, "The IPs should match")

	freeIP, err = gm.CreateLoadBalancer(ctx, virtualServiceNamePrefix,
		lbPoolNamePrefix, []string{"1.2.3.4", "1.2.3.5"}, portDetailsList, cloudConfig.LB.OneArm, cloudConfig.ClusterID)
	assert.NoError(t, err, "Load Balancer should be created even on second attempt")
	assert.NotEmpty(t, freeIP, "There should be a non-empty IP returned")

	updatedIps := []string{"5.5.5.5"}
	updatedInternalPort := int32(55555)
	// update IPs and internal port
	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-http", virtualServiceNamePrefix+"-http", updatedIps, updatedInternalPort, 80)
	assert.NoError(t, err, "HTTP Load Balancer should be updated")

	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-https", virtualServiceNamePrefix+"-https", updatedIps, updatedInternalPort, 443)
	assert.NoError(t, err, "HTTPS Load Balancer should be updated")

	// update external port only
	updatedExternalPortHttp := int32(8080)
	updatedExternalPortHttps := int32(8443)
	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-http", virtualServiceNamePrefix+"-http", updatedIps, updatedInternalPort, updatedExternalPortHttp)
	assert.NoError(t, err, "HTTP Load Balancer should be updated")

	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-https", virtualServiceNamePrefix+"-https", updatedIps, updatedInternalPort, updatedExternalPortHttps)
	assert.NoError(t, err, "HTTPS Load Balancer should be updated")

	// No error on repeated update
	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-http", virtualServiceNamePrefix+"-http", updatedIps, updatedInternalPort, updatedExternalPortHttp)
	assert.NoError(t, err, "HTTP Load Balancer should be updated")

	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-https", virtualServiceNamePrefix+"-https", updatedIps, updatedInternalPort, updatedExternalPortHttps)
	assert.NoError(t, err, "HTTPS Load Balancer should be updated")

	err = gm.DeleteLoadBalancer(ctx, virtualServiceNamePrefix, lbPoolNamePrefix, portDetailsList, cloudConfig.LB.OneArm, cloudConfig.ClusterID)
	assert.NoError(t, err, "Load Balancer should be deleted")

	freeIPObtained, err = gm.GetLoadBalancer(ctx, virtualServiceNameHttp, cloudConfig.LB.OneArm)
	assert.NoError(t, err, "Load Balancer should not be found")
	assert.Empty(t, freeIPObtained, "The VIP should not be found")

	freeIPObtained, err = gm.GetLoadBalancer(ctx, virtualServiceNameHttps, cloudConfig.LB.OneArm)
	assert.NoError(t, err, "Load Balancer should not be found")
	assert.Empty(t, freeIPObtained, "The VIP should not be found")

	err = gm.DeleteLoadBalancer(ctx, virtualServiceNamePrefix, lbPoolNamePrefix, portDetailsList, cloudConfig.LB.OneArm, cloudConfig.ClusterID)
	assert.NoError(t, err, "Repeated deletion of Load Balancer should not fail")

	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-http", virtualServiceNamePrefix+"-http", updatedIps, updatedInternalPort, 80)
	assert.Error(t, err, "updating deleted HTTP Load Balancer should be an error")
	err = gm.UpdateLoadBalancer(ctx, lbPoolNamePrefix+"-https", virtualServiceNamePrefix+"https", updatedIps, updatedInternalPort, 43)
	assert.Error(t, err, "updating deleted HTTPS Load Balancer should be an error")

	return
}

func TestUpdateRDEUsingEtag(t *testing.T) {
	// TODO: This test will currently fail unless the code below is uncommented. Refer to VCDA-3600

	cloudConfig, err := getTestConfig()
	assert.NoError(t, err, "There should be no error opening and parsing cloud config file contents.")

	authFile := filepath.Join(gitRoot, "testdata/auth_test.yaml")
	authFileContent, err := ioutil.ReadFile(authFile)
	assert.NoError(t, err, "There should be no error reading the auth file contents.")

	var authDetails authorizationDetails
	err = yaml.Unmarshal(authFileContent, &authDetails)
	assert.NoError(t, err, "There should be no error parsing auth file content.")

	vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
		"user":    authDetails.Username,
		"secret":  authDetails.Password,
		"userOrg": authDetails.UserOrg,
	})

	ctx := context.Background()

	rm := NewRDEManager(vcdClient, cloudConfig.ClusterID)

	// get rde Vips
	rdeVips1, etag1, defEnt1, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Should retrieve RDE vips on first attempt")
	rdeVips2, etag2, defEnt2, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Should retrieve RDE vips on second attempt")
	assert.Equal(t, etag1, etag2, "etags from consecutive RDE GET calls should match")
	origRdeVips := make([]string, len(rdeVips1))
	copy(origRdeVips, rdeVips1)

	// Test successfully updating using first etag
	addIp1 := "1.2.3.4"
	addIp2 := "2.3.4.5"
	updatedRdeVips1 := append(rdeVips1, addIp1)
	httpResponse1, err := rm.updateRDEVirtualIps(ctx, updatedRdeVips1, etag1, defEnt1)
	assert.NoError(t, err, "RDE should be updated")
	assert.Equal(t, http.StatusOK, httpResponse1.StatusCode, "RDE update status code should be 200 (OK)")
	rdeVips3, _, _, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Should retrieve RDE vips successfully")
	assert.True(t, foundStringInSlice(addIp1, rdeVips3), "ip [%s] should be found in rde vips", addIp1)

	// Test adding addIp2 with outdated etag
	updatedRdeVips2 := append(rdeVips2, addIp2)
	httpResponse2, err := rm.updateRDEVirtualIps(ctx, updatedRdeVips2, etag2, defEnt2)
	assert.Error(t, err, "Should have an error updating RDE with outdated etag")
	assert.Equal(t, http.StatusPreconditionFailed, httpResponse2.StatusCode, "RDE update status code should be 412 (Precondition failed)")
	rdeVips3, etag3, defEnt3, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Should retrieve RDE vips successfully")
	assert.False(t, foundStringInSlice(addIp2, rdeVips3), "ip [%s] should not be found in rde vips", addIp2)

	// Try adding addIp2 with correct etag
	updatedRdeVips3 := append(rdeVips3, addIp2)
	httpResponse3, err := rm.updateRDEVirtualIps(ctx, updatedRdeVips3, etag3, defEnt3)
	assert.NoError(t, err, "RDE should be updated")
	assert.Equal(t, http.StatusOK, httpResponse3.StatusCode, "RDE update status code should be 200 (OK)")
	rdeVips4, etag4, defEnt4, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Should retrieve RDE vips successfully")
	assert.True(t, foundStringInSlice(addIp2, rdeVips4), "ip [%s] should be found in rde vips", addIp2)

	// reset RDE vips to original state
	httpResponse4, err := rm.updateRDEVirtualIps(ctx, rdeVips1, etag4, defEnt4)
	assert.NoError(t, err, "RDE should be updated")
	assert.Equal(t, http.StatusOK, httpResponse4.StatusCode, "RDE update status code should be 200 (OK)")
	// no check to ensure ip's removed because they may have been previously present in the RDE vips
	rdeVips5, _, _, err := rm.GetRDEVirtualIps(ctx)
	assert.NoError(t, err, "Should retrieve RDE vips to check added ips are removed")
	assert.False(t, foundStringInSlice(addIp1, rdeVips5), "ip [%s] should not be found in rde vips", addIp1)
	assert.False(t, foundStringInSlice(addIp2, rdeVips5), "ip [%s] should not be found in rde vips", addIp2)
}
