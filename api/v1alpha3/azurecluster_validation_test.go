/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestClusterWithPreexistingVnetValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		cluster *AzureCluster
	}

	testCase := test{
		name:    "azurecluster with pre-existing vnet - valid",
		cluster: createValidCluster(),
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := testCase.cluster.validateCluster()
		g.Expect(err).To(BeNil())
	})
}

func TestClusterWithPreexistingVnetInvalid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		cluster *AzureCluster
	}

	testCase := test{
		name:    "azurecluster with pre-existing vnet - invalid",
		cluster: createValidCluster(),
	}

	// invalid because it doesn't specify a controlplane subnet
	testCase.cluster.Spec.NetworkSpec.Subnets[0] = &SubnetSpec{
		Name: "random-subnet",
		Role: "random",
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := testCase.cluster.validateCluster()
		g.Expect(err).ToNot(BeNil())
	})
}

func TestClusterWithoutPreexistingVnetValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		cluster *AzureCluster
	}

	testCase := test{
		name:    "azurecluster without pre-existing vnet - valid",
		cluster: createValidCluster(),
	}

	// When ResourceGroup is an empty string, the cluster doesn't
	// have a pre-existing vnet.
	testCase.cluster.Spec.NetworkSpec.Vnet.ResourceGroup = ""

	t.Run(testCase.name, func(t *testing.T) {
		err := testCase.cluster.validateCluster()
		g.Expect(err).To(BeNil())
	})
}

func TestClusterSpecWithPreexistingVnetValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		cluster *AzureCluster
	}

	testCase := test{
		name:    "azurecluster spec with pre-existing vnet - valid",
		cluster: createValidCluster(),
	}

	t.Run(testCase.name, func(t *testing.T) {
		errs := testCase.cluster.validateClusterSpec()
		g.Expect(errs).To(BeNil())
	})
}

func TestClusterSpecWithPreexistingVnetInvalid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		cluster *AzureCluster
	}

	testCase := test{
		name:    "azurecluster spec with pre-existing vnet - invalid",
		cluster: createValidCluster(),
	}

	// invalid because it doesn't specify a controlplane subnet
	testCase.cluster.Spec.NetworkSpec.Subnets[0] = &SubnetSpec{
		Name: "random-subnet",
		Role: "random",
	}

	t.Run(testCase.name, func(t *testing.T) {
		errs := testCase.cluster.validateClusterSpec()
		g.Expect(len(errs)).To(BeNumerically(">", 0))
	})
}

func TestClusterSpecWithoutPreexistingVnetValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		cluster *AzureCluster
	}

	testCase := test{
		name:    "azurecluster spec without pre-existing vnet - valid",
		cluster: createValidCluster(),
	}

	// When ResourceGroup is an empty string, the cluster doesn't
	// have a pre-existing vnet.
	testCase.cluster.Spec.NetworkSpec.Vnet.ResourceGroup = ""

	t.Run(testCase.name, func(t *testing.T) {
		errs := testCase.cluster.validateClusterSpec()
		g.Expect(errs).To(BeNil())
	})
}

func TestNetworkSpecWithPreexistingVnetValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name        string
		networkSpec NetworkSpec
	}

	testCase := test{
		name:        "azurecluster networkspec with pre-existing vnet - valid",
		networkSpec: createValidNetworkSpec(),
	}

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateNetworkSpec(testCase.networkSpec, field.NewPath("spec").Child("networkSpec"))
		g.Expect(errs).To(BeNil())
	})
}

func TestNetworkSpecWithPreexistingVnetLackRequiredSubnets(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name        string
		networkSpec NetworkSpec
	}

	testCase := test{
		name:        "azurecluster networkspec with pre-existing vnet - lack required subnets",
		networkSpec: createValidNetworkSpec(),
	}

	// invalid because it doesn't specify a node subnet
	testCase.networkSpec.Subnets = testCase.networkSpec.Subnets[:1]

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateNetworkSpec(testCase.networkSpec, field.NewPath("spec").Child("networkSpec"))
		g.Expect(errs).To(HaveLen(1))
		g.Expect(errs[0].Type).To(Equal(field.ErrorTypeRequired))
		g.Expect(errs[0].Field).To(Equal("spec.networkSpec.subnets"))
		g.Expect(errs[0].Error()).To(ContainSubstring("required role node not included"))
	})
}

func TestNetworkSpecWithPreexistingVnetInvalidResourceGroup(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name        string
		networkSpec NetworkSpec
	}

	testCase := test{
		name:        "azurecluster networkspec with pre-existing vnet - invalid resource group",
		networkSpec: createValidNetworkSpec(),
	}

	testCase.networkSpec.Vnet.ResourceGroup = "invalid-name###"

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateNetworkSpec(testCase.networkSpec, field.NewPath("spec").Child("networkSpec"))
		g.Expect(errs).To(HaveLen(1))
		g.Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		g.Expect(errs[0].Field).To(Equal("spec.networkSpec.vnet.resourceGroup"))
		g.Expect(errs[0].BadValue).To(BeEquivalentTo(testCase.networkSpec.Vnet.ResourceGroup))
	})
}

func TestNetworkSpecWithoutPreexistingVnetValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name        string
		networkSpec NetworkSpec
	}

	testCase := test{
		name:        "azurecluster networkspec without pre-existing vnet - valid",
		networkSpec: createValidNetworkSpec(),
	}

	testCase.networkSpec.Vnet.ResourceGroup = ""

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateNetworkSpec(testCase.networkSpec, field.NewPath("spec").Child("networkSpec"))
		g.Expect(errs).To(BeNil())
	})
}

func TestResourceGroupValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name          string
		resourceGroup string
	}

	testCase := test{
		name:          "resourcegroup name - valid",
		resourceGroup: "custom-vnet",
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := validateResourceGroup(testCase.resourceGroup,
			field.NewPath("spec").Child("networkSpec").Child("vnet").Child("resourceGroup"))
		g.Expect(err).To(BeNil())
	})
}

func TestResourceGroupInvalid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name          string
		resourceGroup string
	}

	testCase := test{
		name:          "resourcegroup name - invalid",
		resourceGroup: "inv@lid-rg",
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := validateResourceGroup(testCase.resourceGroup,
			field.NewPath("spec").Child("networkSpec").Child("vnet").Child("resourceGroup"))
		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Type).To(Equal(field.ErrorTypeInvalid))
		g.Expect(err.Field).To(Equal("spec.networkSpec.vnet.resourceGroup"))
		g.Expect(err.BadValue).To(BeEquivalentTo(testCase.resourceGroup))
	})
}

func TestSubnetsValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		subnets Subnets
	}

	testCase := test{
		name:    "subnets - valid",
		subnets: createValidSubnets(),
	}

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateSubnets(testCase.subnets,
			field.NewPath("spec").Child("networkSpec").Child("subnets"))
		g.Expect(errs).To(BeNil())
	})
}

func TestSubnetsInvalidSubnetName(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		subnets Subnets
	}

	testCase := test{
		name:    "subnets - invalid subnet name",
		subnets: createValidSubnets(),
	}

	testCase.subnets[0].Name = "invalid-subnet-name-due-to-bracket)"

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateSubnets(testCase.subnets,
			field.NewPath("spec").Child("networkSpec").Child("subnets"))
		g.Expect(errs).To(HaveLen(1))
		g.Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		g.Expect(errs[0].Field).To(Equal("spec.networkSpec.subnets[0].name"))
		g.Expect(errs[0].BadValue).To(BeEquivalentTo("invalid-subnet-name-due-to-bracket)"))
	})
}

func TestSubnetsInvalidInternalLBIPAddress(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		subnets Subnets
	}

	testCase := test{
		name:    "subnets - invalid internal load balancer ip address",
		subnets: createValidSubnets(),
	}

	testCase.subnets[0].InternalLBIPAddress = "2550.1.1.1"

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateSubnets(testCase.subnets,
			field.NewPath("spec").Child("networkSpec").Child("subnets"))
		g.Expect(errs).To(HaveLen(1))
		g.Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		g.Expect(errs[0].Field).To(Equal("spec.networkSpec.subnets[0].internalLBIPAddress"))
		g.Expect(errs[0].BadValue).To(BeEquivalentTo("2550.1.1.1"))
	})
}

func TestSubnetsInvalidLackRequiredSubnet(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		subnets Subnets
	}

	testCase := test{
		name:    "subnets - lack required subnet",
		subnets: createValidSubnets(),
	}

	testCase.subnets[0].Role = "random-role"

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateSubnets(testCase.subnets,
			field.NewPath("spec").Child("networkSpec").Child("subnets"))
		g.Expect(errs).To(HaveLen(1))
		g.Expect(errs[0].Type).To(Equal(field.ErrorTypeRequired))
		g.Expect(errs[0].Field).To(Equal("spec.networkSpec.subnets"))
		g.Expect(errs[0].Detail).To(ContainSubstring("required role control-plane not included"))
	})
}

func TestSubnetNamesNotUnique(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name    string
		subnets Subnets
	}

	testCase := test{
		name:    "subnets - names not unique",
		subnets: createValidSubnets(),
	}

	testCase.subnets[0].Name = "subnet-name"
	testCase.subnets[1].Name = "subnet-name"

	t.Run(testCase.name, func(t *testing.T) {
		errs := validateSubnets(testCase.subnets,
			field.NewPath("spec").Child("networkSpec").Child("subnets"))
		g.Expect(errs).To(HaveLen(1))
		g.Expect(errs[0].Type).To(Equal(field.ErrorTypeDuplicate))
		g.Expect(errs[0].Field).To(Equal("spec.networkSpec.subnets"))
	})
}

func TestSubnetNameValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name       string
		subnetName string
	}

	testCase := test{
		name:       "subnet name - valid",
		subnetName: "control-plane-subnet",
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := validateSubnetName(testCase.subnetName,
			field.NewPath("spec").Child("networkSpec").Child("subnets").Index(0).Child("name"))
		g.Expect(err).To(BeNil())
	})
}

func TestSubnetNameInvalid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name       string
		subnetName string
	}

	testCase := test{
		name:       "subnet name - invalid",
		subnetName: "inv@lid-subnet-name",
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := validateSubnetName(testCase.subnetName,
			field.NewPath("spec").Child("networkSpec").Child("subnets").Index(0).Child("name"))
		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Type).To(Equal(field.ErrorTypeInvalid))
		g.Expect(err.Field).To(Equal("spec.networkSpec.subnets[0].name"))
		g.Expect(err.BadValue).To(BeEquivalentTo(testCase.subnetName))
	})
}

func TestInternalLBIPAddressValid(t *testing.T) {
	g := NewWithT(t)

	type test struct {
		name                string
		internalLBIPAddress string
	}

	testCase := test{
		name:                "subnet name - invalid",
		internalLBIPAddress: "1.1.1.1",
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := validateInternalLBIPAddress(testCase.internalLBIPAddress,
			field.NewPath("spec").Child("networkSpec").Child("subnets").Index(0).Child("internalLBIPAddress"))
		g.Expect(err).To(BeNil())
	})
}

func TestInternalLBIPAddressInvalid(t *testing.T) {
	g := NewWithT(t)

	internalLBIPAddress := "1.1.1"

	err := validateInternalLBIPAddress(internalLBIPAddress,
		field.NewPath("spec").Child("networkSpec").Child("subnets").Index(0).Child("internalLBIPAddress"))
	g.Expect(err).NotTo(BeNil())
	g.Expect(err.Type).To(Equal(field.ErrorTypeInvalid))
	g.Expect(err.Field).To(Equal("spec.networkSpec.subnets[0].internalLBIPAddress"))
	g.Expect(err.BadValue).To(BeEquivalentTo(internalLBIPAddress))
}

func TestIngressRules(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name      string
		validRule *IngressRule
		wantErr   bool
	}{
		{
			name: "ingressRule - valid priority",
			validRule: &IngressRule{
				Name:        "allow_apiserver",
				Description: "Allow K8s API Server",
				Priority:    101,
			},
			wantErr: false,
		},
		{
			name: "ingressRule - invalid low priority",
			validRule: &IngressRule{
				Name:        "allow_apiserver",
				Description: "Allow K8s API Server",
				Priority:    99,
			},
			wantErr: true,
		},
		{
			name: "ingressRule - invalid high priority",
			validRule: &IngressRule{
				Name:        "allow_apiserver",
				Description: "Allow K8s API Server",
				Priority:    5000,
			},
			wantErr: true,
		},
	}
	for _, testCase := range tests {

		t.Run(testCase.name, func(t *testing.T) {
			err := validateIngressRule(
				testCase.validRule,
				field.NewPath("spec").Child("networkSpec").Child("subnets").Index(0).Child("securityGroup").Child("ingressRules").Index(0),
			)
			if testCase.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
		})
	}
}

func createValidCluster() *AzureCluster {
	return &AzureCluster{
		Spec: AzureClusterSpec{
			NetworkSpec: createValidNetworkSpec(),
		},
	}
}

func createValidNetworkSpec() NetworkSpec {
	return NetworkSpec{
		Vnet: VnetSpec{
			ResourceGroup: "custom-vnet",
			Name:          "my-vnet",
		},
		Subnets: createValidSubnets(),
	}
}

func createValidSubnets() Subnets {
	return Subnets{
		{
			Name: "control-plane-subnet",
			Role: "control-plane",
		},
		{
			Name: "node-subnet",
			Role: "node",
		},
	}
}
