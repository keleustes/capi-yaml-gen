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

package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/keleustes/capi-yaml-gen/pkg/generate"
)

const (
	defaultClusterName            = "my-cluster"
	defaultNamespace              = "default"
	defaultInfrastructureProvider = "docker"
	defaultBootstrapProvider      = "kubeadm"
	defaultVersion                = "v1.15.3"
	defaultControlPlaneCount      = 1
	defaultWorkerCount            = 1
)

// RootCmd returns the root command for capi-yaml-gen tool
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "capi-yaml",
		Short: "Devtool to help with YAML for CAPI and CAPI providers",
		Long:  "Devtool to help with YAML for CAPI and CAPI providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.AddCommand(getGenerateCommand())

	return rootCmd
}

func getGenerateCommand() *cobra.Command {
	opts := generate.GenerateOptions{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate yaml for CAPI and its providers",
		Long:  "generate yaml for CAPI and its providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generate.RunGenerateCommand(opts, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&opts.ClusterName, "cluster-name", "c", defaultClusterName, "Name for the cluster")
	cmd.Flags().StringVarP(&opts.ClusterNamespace, "namespace", "n", defaultNamespace, "Namespace where the cluster will be created")
	cmd.Flags().StringVarP(&opts.InfraProvider, "infrastructure-provider", "i", defaultInfrastructureProvider, "Infrastructure provider for the cluster")
	cmd.Flags().StringVarP(&opts.BsProvider, "boostrap-provider", "b", defaultBootstrapProvider, "Bootstrap provider for the cluster")
	cmd.Flags().StringVarP(&opts.K8sVersion, "k8s-version", "k", defaultVersion, "Version of kubernetes for the cluster")
	cmd.Flags().BoolVarP(&opts.MachineDeployment, "generate-machine-deployment", "d", true, "Generate a machine deployment instead of individual machines")

	cmd.Flags().IntVarP(&opts.ControlplaneMachineCount, "control-plane-count", "m", defaultControlPlaneCount, "Number of control plane machines in the cluster")
	cmd.Flags().IntVarP(&opts.WorkerMachineCount, "worker-count", "w", defaultWorkerCount, "Number of worker machines in the cluster")

	return cmd
}

// Execute starts the process
func Execute() {
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}
