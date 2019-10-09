package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/keleustes/capi-yaml-gen/pkg/generate"
)

var update = flag.Bool("update", false, "update golden files")

func TestGoldenFiles(t *testing.T) {
	testcases := []struct {
		name       string
		goldenfile string
		options    generate.GenerateOptions
	}{
		{
			"./capi-yaml-gen generate --generate-machine-deployment=false",
			"default-capd-no-machine-deployment",
			generate.GenerateOptions{
				InfraProvider:            defaultInfrastructureProvider,
				ClusterName:              defaultClusterName,
				ClusterNamespace:         defaultNamespace,
				BsProvider:               defaultBootstrapProvider,
				K8sVersion:               defaultVersion,
				ControlplaneMachineCount: defaultControlPlaneCount,
				WorkerMachineCount:       defaultWorkerCount,
				MachineDeployment:        false,
			},
		},
		{
			"./capi-yaml-gen generate --infrastructure-provider aws",
			"default-capa-no-machine-deployment",
			generate.GenerateOptions{
				InfraProvider:            "aws",
				ClusterName:              defaultClusterName,
				ClusterNamespace:         defaultNamespace,
				BsProvider:               defaultBootstrapProvider,
				K8sVersion:               defaultVersion,
				ControlplaneMachineCount: defaultControlPlaneCount,
				WorkerMachineCount:       defaultWorkerCount,
				MachineDeployment:        false,
			},
		},
		{
			"./capi-yaml-gen generate",
			"default-capd",
			generate.GenerateOptions{
				InfraProvider:            defaultInfrastructureProvider,
				ClusterName:              defaultClusterName,
				ClusterNamespace:         defaultNamespace,
				BsProvider:               defaultBootstrapProvider,
				K8sVersion:               defaultVersion,
				ControlplaneMachineCount: defaultControlPlaneCount,
				WorkerMachineCount:       defaultWorkerCount,
				MachineDeployment:        true,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			if err := generate.RunGenerateCommand(tc.options, &stdout); err != nil {
				t.Fatal(err)
			}

			if *update {
				if err := ioutil.WriteFile(goldenFileName(tc.goldenfile), stdout.Bytes(), 0644); err != nil {
					t.Fatal(err)
				}
				return
			}

			golden, err := ioutil.ReadFile(goldenFileName(tc.goldenfile))
			if err != nil {
				t.Fatal(err)
			}
			diff := cmp.Diff(string(golden), stdout.String())
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func goldenFileName(name string) string {
	return fmt.Sprintf("testdata/%s.golden", name)
}
