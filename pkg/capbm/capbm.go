/*
Copyright 2019 The Kubernetes Authors.

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

package capbm

import (
	"github.com/keleustes/capi-yaml-gen/pkg/constants"
	"github.com/keleustes/capi-yaml-gen/pkg/generator"
	corev1 "k8s.io/api/core/v1"
	infrav2 "sigs.k8s.io/cluster-api-provider-baremetal/api/v1alpha3"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha2"
	bootstrapv1beta1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/v1beta1"
)

// Provider CAPA implementation of InfrastructureProvider
type Provider struct{}

// GetInfraCluster generates an BareMetal cluster
func (p Provider) GetInfraCluster(name, namespace string) generator.Object {
	baremetalCluster := &infrav2.BareMetalCluster{}
	baremetalCluster.Kind = constants.BareMetalClusterKind
	baremetalCluster.APIVersion = infrav2.GroupVersion.String()
	baremetalCluster.Name = name
	baremetalCluster.Namespace = namespace
	return baremetalCluster
}

// GetInfraMachine generates an BareMetal machine
func (p Provider) GetInfraMachine(name, namespace string) generator.Object {
	baremetalMachine := &infrav2.BareMetalMachine{}
	baremetalMachine.Kind = constants.BareMetalMachineKind
	baremetalMachine.APIVersion = infrav2.GroupVersion.String()
	baremetalMachine.Name = name
	baremetalMachine.Namespace = namespace
	baremetalMachine.Spec = infrav2.BareMetalMachineSpec{
		Image: infrav2.Image{
			URL:      "http://172.22.0.1/images/rhcos-ootpa-latest.qcow2",
			Checksum: "http://172.22.0.1/images/rhcos-ootpa-latest.qcow2.md5sum",
		},
		UserData: &corev1.SecretReference{
			Name:      "worker-user-data",
			Namespace: "otherns",
		},
	}
	return baremetalMachine
}

// GetInfraMachineTemplate generates an AWS machine template
func (p Provider) GetInfraMachineTemplate(name, namespace string) generator.Object {
	template := &infrav2.BareMetalMachineTemplate{}
	template.Name = name
	template.Namespace = namespace
	template.Kind = constants.BareMetalMachineKind + "Template"
	template.APIVersion = infrav2.GroupVersion.String()
	return template
}

// SetBootstrapConfigInfraValues fills in InfraProvider specific values into the bootstrap config
func (p Provider) SetBootstrapConfigInfraValues(c *bootstrapv1.KubeadmConfig) {
	extraArgs := map[string]string{
		"cloud-provider": "baremetal",
	}
	if c.Spec.InitConfiguration != nil {
		c.Spec.InitConfiguration.NodeRegistration = bootstrapv1beta1.NodeRegistrationOptions{
			Name:             "'{{ ds.meta_data.hostname }}'",
			KubeletExtraArgs: extraArgs,
		}
	} else if c.Spec.JoinConfiguration != nil {
		c.Spec.JoinConfiguration.NodeRegistration = bootstrapv1beta1.NodeRegistrationOptions{
			Name:             "'{{ ds.meta_data.hostname }}'",
			KubeletExtraArgs: extraArgs,
		}
	}

	if c.Spec.ClusterConfiguration != nil {
		c.Spec.ClusterConfiguration.APIServer = bootstrapv1beta1.APIServer{
			ControlPlaneComponent: bootstrapv1beta1.ControlPlaneComponent{
				ExtraArgs: extraArgs,
			},
		}

		c.Spec.ClusterConfiguration.ControllerManager = bootstrapv1beta1.ControlPlaneComponent{
			ExtraArgs: extraArgs,
		}
	}
}

// SetBootstrapConfigTemplateInfraValues fills in InfraProvider specific values into the join configuration
func (p Provider) SetBootstrapConfigTemplateInfraValues(t *bootstrapv1.KubeadmConfigTemplate) {
	extraArgs := map[string]string{
		"cloud-provider": "baremetal",
	}
	if t.Spec.Template.Spec.JoinConfiguration != nil {
		t.Spec.Template.Spec.JoinConfiguration.NodeRegistration = bootstrapv1beta1.NodeRegistrationOptions{
			Name:             "'{{ ds.meta_data.hostname }}'",
			KubeletExtraArgs: extraArgs,
		}
	}
}

func (p Provider) GetEnvironmentVariables() map[string]string {
	return map[string]string{
		"SSH_KEY_NAME":                     "default",
		"CONTROL_PLANE_INSTANCE_TYPE":      "t2.medium",
		"MACHINE_DEPLOYMENT_INSTANCE_TYPE": "t2.medium",
		"REGION":                           "us-west-2",
	}
}
