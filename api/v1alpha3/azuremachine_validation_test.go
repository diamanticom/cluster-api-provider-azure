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
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"testing"

	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestAzureMachine_ValidateSSHKey(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name    string
		sshKey  string
		wantErr bool
	}{
		{
			name:    "valid ssh key",
			sshKey:  generateSSHPublicKey(),
			wantErr: false,
		},
		{
			name:    "invalid ssh key",
			sshKey:  "invalid ssh key",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSSHKey(tc.sshKey, field.NewPath("sshPublicKey"))
			if tc.wantErr {
				g.Expect(err).ToNot(HaveLen(0))
			} else {
				g.Expect(err).To(HaveLen(0))
			}
		})
	}
}

func generateSSHPublicKey() string {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicRsaKey, _ := ssh.NewPublicKey(&privateKey.PublicKey)
	return base64.StdEncoding.EncodeToString(ssh.MarshalAuthorizedKey(publicRsaKey))
}

type osDiskTestInput struct {
	name    string
	wantErr bool
	osDisk  OSDisk
}

func TestAzureMachine_ValidateOSDisk(t *testing.T) {
	g := NewWithT(t)

	testcases := []osDiskTestInput{
		{
			name:    "valid os disk spec",
			wantErr: false,
			osDisk:  generateValidOSDisk(),
		},
	}
	testcases = append(testcases, generateNegativeTestCases()...)

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateOSDisk(test.osDisk, field.NewPath("osDisk"))
			if test.wantErr {
				g.Expect(err).NotTo(HaveLen(0))
			} else {
				g.Expect(err).To(HaveLen(0))
			}
		})
	}
}

func generateNegativeTestCases() []osDiskTestInput {
	inputs := []osDiskTestInput{}
	testCaseName := "invalid os disk spec"

	invalidDiskSpecs := []OSDisk{
		{},
		{
			DiskSizeGB: 0,
			OSType:     "blah",
		},
		{
			DiskSizeGB: -10,
			OSType:     "blah",
		},
		{
			DiskSizeGB: 2050,
			OSType:     "blah",
		},
		{
			DiskSizeGB: 20,
			OSType:     "",
		},
		{
			DiskSizeGB:  30,
			OSType:      "blah",
			ManagedDisk: ManagedDisk{},
		},
		{
			DiskSizeGB: 30,
			OSType:     "blah",
			ManagedDisk: ManagedDisk{
				StorageAccountType: "",
			},
		},
		{
			DiskSizeGB: 30,
			OSType:     "blah",
			ManagedDisk: ManagedDisk{
				StorageAccountType: "invalid_type",
			},
		},
	}

	for _, input := range invalidDiskSpecs {
		inputs = append(inputs, osDiskTestInput{
			name:    testCaseName,
			wantErr: true,
			osDisk:  input,
		})
	}

	return inputs
}

func generateValidOSDisk() OSDisk {
	return OSDisk{
		DiskSizeGB: 30,
		OSType:     "Linux",
		ManagedDisk: ManagedDisk{
			StorageAccountType: "Premium_LRS",
		},
	}
}
