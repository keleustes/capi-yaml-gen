# capi-yaml

Developer tool to generate yaml for Cluster-API and necessary provider resources.

Read more about this [here](https://docs.google.com/document/d/1Tzx6IXOoQnUxaVSYA2I8IdcNrFE4zgOnzkk55KHOU20/edit)

## Usage Examples

### Docker Infrastructure and Kubeadm Bootstrap Providerrs

- InfraProvider: [Docker](https://github.com/kubernetes-sigs/cluster-api-provider-docker)
- BootstrapProvider: [Kubeadm](https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm)
- KubernetesVersion: v1.14.2
- ControlPlaneMachineCount: 3
- WorkerMachineCountt: 1

```(bash)
$ go run main.go generate --controlplane-count 3
# yaml written to stdout
```

### AWS Infrastructure and Kubeadm Bootstrap Providers

- InfraProvider: [AWS](https://github.com/kubernetes-sigs/cluster-api-provider-aws)
- BootstrapProvider: [Kubeadm](https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm)
- KubernetesVersion: v1.14.2
- ControlPlaneMachineCount: 3
- WorkerMachineCountt: 1

```(bash)
$ go run main.go generate --control-plane-count 3 --infrastructure-provider aws
# yaml written to stdout
```

### MachineDeployments

By default workers will be managed by a MachineDeployment. If you do not want this behavior set the
`--generate-machined-deployment` flag to false like this:

```bash
go run main.go generate --generate-machine-deployment=false
```

### Customizations

Some providers require fields custom to each user. We provider some default values that can be overridden by setting
appropriate environment variables. For instance, if you do not have any environment variables set and you select the
`aws` provider, you will be greeted with a message that looks something like this:

```text
Consider setting these default values and rerunning.
If you do not want to interpolate the values, rerun with the --allow-empty-env-vars flag.

export REGION=us-west-2
export SSH_KEY_NAME=default
export CONTROL_PLANE_INSTANCE_TYPE=t2.medium
```

You can customize the environment variables for your use case.  If you're just looking for some YAML to modify by hand
and don't want to supply any values you can use the `--allow-empty-env-vars` flag to skip verification. This will
interpolate all environment variables you have set and ignore ones you do not have set.

## Output Examples

Generation of a baremetal cluster named capigen-demo with 3 control-plane nodes and 9 worker nodes.

```bash
GO111MODULE=on go install
capi-yaml-gen generate -i baremetal -c capigen-demo -m 3 -k 1.16.1 -n capigen-namespace -w 9 > capigen-demo.yaml
```

```yaml
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: BareMetalCluster
metadata:
  creationTimestamp: null
  name: capigen-demo
  namespace: capigen-namespace
spec: {}
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Cluster
metadata:
  creationTimestamp: null
  name: capigen-demo
  namespace: capigen-namespace
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
      - 192.168.0.0/16
    services:
      cidrBlocks: []
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: BareMetalCluster
    name: capigen-demo
    namespace: capigen-namespace
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: BareMetalMachine
metadata:
  creationTimestamp: null
  name: controlplane-0
  namespace: capigen-namespace
spec:
  hostSelector: {}
  image:
    checksum: http://172.22.0.1/images/rhcos-ootpa-latest.qcow2.md5sum
    url: http://172.22.0.1/images/rhcos-ootpa-latest.qcow2
  userData:
    name: worker-user-data
    namespace: otherns
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Machine
metadata:
  creationTimestamp: null
  labels:
    cluster.x-k8s.io/cluster-name: capigen-demo
    cluster.x-k8s.io/control-plane: "true"
  name: controlplane-0
  namespace: capigen-namespace
spec:
  bootstrap:
    configRef:
      apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
      kind: KubeadmConfig
      name: controlplane-0-config
      namespace: capigen-namespace
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: BareMetalMachine
    name: controlplane-0
    namespace: capigen-namespace
  metadata: {}
---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: KubeadmConfig
metadata:
  creationTimestamp: null
  name: controlplane-0-config
  namespace: capigen-namespace
spec:
  clusterConfiguration:
    apiServer:
      extraArgs:
        cloud-provider: baremetal
    certificatesDir: ""
    controlPlaneEndpoint: ""
    controllerManager:
      extraArgs:
        cloud-provider: baremetal
    dns:
      type: ""
    etcd: {}
    imageRepository: ""
    kubernetesVersion: ""
    networking:
      dnsDomain: ""
      podSubnet: ""
      serviceSubnet: ""
    scheduler: {}
  initConfiguration:
    localAPIEndpoint:
      advertiseAddress: ""
      bindPort: 0
    nodeRegistration:
      kubeletExtraArgs:
        cloud-provider: baremetal
      name: '''{{ ds.meta_data.hostname }}'''
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: BareMetalMachine
metadata:
  creationTimestamp: null
  name: controlplane-1
  namespace: capigen-namespace
spec:
  hostSelector: {}
  image:
    checksum: http://172.22.0.1/images/rhcos-ootpa-latest.qcow2.md5sum
    url: http://172.22.0.1/images/rhcos-ootpa-latest.qcow2
  userData:
    name: worker-user-data
    namespace: otherns
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Machine
metadata:
  creationTimestamp: null
  labels:
    cluster.x-k8s.io/cluster-name: capigen-demo
    cluster.x-k8s.io/control-plane: "true"
  name: controlplane-1
  namespace: capigen-namespace
spec:
  bootstrap:
    configRef:
      apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
      kind: KubeadmConfig
      name: controlplane-1-config
      namespace: capigen-namespace
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: BareMetalMachine
    name: controlplane-1
    namespace: capigen-namespace
  metadata: {}
---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: KubeadmConfig
metadata:
  creationTimestamp: null
  name: controlplane-1-config
  namespace: capigen-namespace
spec:
  joinConfiguration:
    controlPlane:
      localAPIEndpoint:
        advertiseAddress: ""
        bindPort: 6443
    discovery: {}
    nodeRegistration:
      kubeletExtraArgs:
        cloud-provider: baremetal
      name: '''{{ ds.meta_data.hostname }}'''
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: BareMetalMachine
metadata:
  creationTimestamp: null
  name: controlplane-2
  namespace: capigen-namespace
spec:
  hostSelector: {}
  image:
    checksum: http://172.22.0.1/images/rhcos-ootpa-latest.qcow2.md5sum
    url: http://172.22.0.1/images/rhcos-ootpa-latest.qcow2
  userData:
    name: worker-user-data
    namespace: otherns
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Machine
metadata:
  creationTimestamp: null
  labels:
    cluster.x-k8s.io/cluster-name: capigen-demo
    cluster.x-k8s.io/control-plane: "true"
  name: controlplane-2
  namespace: capigen-namespace
spec:
  bootstrap:
    configRef:
      apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
      kind: KubeadmConfig
      name: controlplane-2-config
      namespace: capigen-namespace
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: BareMetalMachine
    name: controlplane-2
    namespace: capigen-namespace
  metadata: {}
---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: KubeadmConfig
metadata:
  creationTimestamp: null
  name: controlplane-2-config
  namespace: capigen-namespace
spec:
  joinConfiguration:
    controlPlane:
      localAPIEndpoint:
        advertiseAddress: ""
        bindPort: 6443
    discovery: {}
    nodeRegistration:
      kubeletExtraArgs:
        cloud-provider: baremetal
      name: '''{{ ds.meta_data.hostname }}'''
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: BareMetalMachineTemplate
metadata:
  creationTimestamp: null
  name: worker-md
  namespace: capigen-namespace
spec:
  template:
    spec:
      hostSelector: {}
      image:
        checksum: ""
        url: ""

---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: KubeadmConfigTemplate
metadata:
  creationTimestamp: null
  name: worker-md
  namespace: capigen-namespace
spec:
  template:
    spec:
      joinConfiguration:
        discovery: {}
        nodeRegistration:
          kubeletExtraArgs:
            cloud-provider: baremetal
          name: '''{{ ds.meta_data.hostname }}'''

---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: MachineDeployment
metadata:
  creationTimestamp: null
  name: worker-md
  namespace: capigen-namespace
spec:
  replicas: 9
  selector:
    matchLabels:
      cluster.x-k8s.io/cluster-name: capigen-demo
  template:
    metadata:
      labels:
        cluster.x-k8s.io/cluster-name: capigen-demo
    spec:
      bootstrap:
        configRef:
          apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
          kind: KubeadmConfigTemplate
          name: worker-md
          namespace: capigen-namespace
      infrastructureRef:
        apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
        kind: BareMetalMachineTemplate
        name: worker-md
        namespace: capigen-namespace
      metadata: {}
      version: 1.16.1

```
