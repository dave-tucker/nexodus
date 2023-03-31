# Deploying the Nexodus Agent

## Deploying on a Node

The following sections contain general instructions to deploy the Nexodus agent on any node. The minimum requirement is that the node runs Linux, Darwin, or Windows-based Operating systems.

### Installing the Agent

#### Install Script

The Nexodus agent (`nexd`) is run on any node that will join a Nexodus Zone to communicate with other peers in that zone. This agent communicates with the Nexodus Controller and manages local wireguard configuration.

The `hack/nexodus_installer.sh` script will download the latest build of `nexd` and install it for you. It will also ensure that `wireguard-tools` has been installed. This installer supports MacOS and Linux. You may also install `wireguard-tools` yourself and build `nexd` from source.

```sh
hack/nexodus_installer.sh
```

#### RPM from a Copr Repository

A Fedora [Copr repository](https://copr.fedorainfracloud.org/coprs/russellb/nexodus/) is updated with new rpms after each new commit to the `main` branch that passes CI. The rpm will include `nexctl`, `nexd`, and integration with systemd.

You can add this repository to your Fedora host with the following command:

```sh
sudo dnf copr enable russellb/nexodus
```

Then you should be able to install Nexodus with:

```sh
sudo dnf install nexodus
```

#### Custom RPM Build

You can also build a custom rpm from the git repository. You must have `mock` installed to build the package.

```sh
make rpm
```

After running this command, the resulting rpm can be found in `./dist/rpm/mock/`.

To install the rpm, you may use `dnf`.

```sh
sudo dnf install ./dist/rpm/mock/nexodus-0-0.1.20230216git068fedd.fc37.src.rpm
```

#### Systemd

If you did not install `nexd` via the rpm, you can still use the systemd integration if you would like. The following commands will put the files in the right place.

```sh
sudo cp contrib/rpm/nexodus.service /usr/lib/systemd/service/nexodus.service
sudo cp contrib/rpm/nexodus.sysconfig /etc/sysconfig/nexodus
sudo systemctl daemon-reload
```

#### Starting the Agent

> **Note**
> In a self-signed dev environment, each agent machine needs to have the [imported cert](../deployment/nexodus-service.md#https) and the [host entry](../deployment/nexodus-service.md#add-required-dns-entries) detailed above.

You may start `nexd` directly. You must include the URL to the Nexodus service as an argument.

```sh
sudo nexd-linux-amd64 https://try.nexodus.local
```

Alternatively, you can start `nexd` as a systemd service. First, edit `/etc/sysconfig/nexodus` to reflect the URL of the Nexodus service. Then, start the agent with the following command:

```sh
sudo systemctl start nexodus
```

If you would like `nexd` to run automatically on boot, run this command as well:

```sh
sudo systemctl enable nexodus
```

### Interactive Enrollment

If the agent is able to successfully reach the controller API, it will provide a one-time code to provide to the controller web UI to complete enrollment of this node into a Nexodus Zone. If you ran `nexd` manually, you will see a message like the following in your terminal:

```sh
Your device must be registered with Nexodus.
Your one-time code is: LTCV-OFFS
Please open the following URL in your browser to sign in:
https://auth.try.nexodus.local/realms/nexodus/device?user_code=LTCV-OFFS
```

If the agent was started using systemd, you will find the same thing in the service logs. You can also retrieve this status information using `nexctl`.

```sh
$ sudo nexctl nexd status
Status: WaitingForAuth
Your device must be registered with Nexodus.
Your one-time code is: LTCV-OFFS
Please open the following URL in your browser to sign in:
https://auth.try.nexodus.local/realms/nexodus/device?user_code=LTCV-OFFS
```

Once enrollment is completed in the web UI, the agent will show progress.

```text
Authentication succeeded.
...
INFO[0570] Peer setup complete
```

### Verifying Agent Setup

Once the Agent has been started successfully, you should see a wireguard interface with an address assigned. For example, on Linux:

```sh
$ ip address show wg0
161: wg0: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1420 qdisc noqueue state UNKNOWN group default qlen 1000
    link/none
    inet 100.100.0.1/32 scope global wg0
       valid_lft forever preferred_lft forever
```

### Verifying Organization Connectivity

Once more than one node has enrolled in the same Nexodus organization, you will see additional routes populated for reaching other nodes' endpoints in the same organization. For example, we have just added a second node to this zone. The new node's address in the Nexodus organization is 100.100.0.2. On Linux, we can check the routing table and see:

```sh
$ ip route
...
100.100.0.2 dev wg0 scope link
```

You should now be able to reach that node over the wireguard tunnel.

```sh
$ ping 100.100.0.2
PING 100.100.0.2 (100.100.0.2) 56(84) bytes of data.
64 bytes from 100.100.0.2: icmp_seq=1 ttl=64 time=7.63 ms
```

You can explore the web UI by visiting the URL of the host you added in your `/etc/hosts` file. For example, `https://try.nexodus.local/`.

### Cleanup Agent From Node

If you want to remove the node from the network, and want to clean up all the configuration done on the node. Fire away following commands:

Ctrl + c (cmd+c) the agent process. and remove the wireguard interface and relevant configuration files.
*Linux:*

```shell
sudo ip link del wg0
```

*OSX/Windows:*

Since the wireguard agents are userspace in both Windows and Darwin, the tunnel interface is removed when the agent process exits.

## Deploying on Kubernetes-managed Nodes

Instructions mentioned in [Deploying on a Node](#deploying-on-a-node) can be used here to deploy the Nexodus agent on Kubernetes-managed nodes. However, deploying the agent across all the nodes in a sizable Kubernetes cluster can be a challenging task. The following section provides instructions to deploy Kubernetes style manifest to automate the deployment process.

### Setup the configuration

Agent deployment in Kubernetes requires a few initial configuration details to successfully deploy the agent and onboard the node. This configuration is provided through `kustomization.yaml`. Make a copy of the sample [`kustomization.yaml.sample`](https://github.com/nexodus-io/nexodus/blob/main/deploy/nexodus-client/overlays/dev/kustomization.yaml.sample) and rename it to `kustomization.yaml`.

Fetch the CA cert from the Kubernetes cluster that is running the Nexodus controller, and set the return cert string to `cert` literal of `secretGenerator` in `kustomization.yaml`.

```sh
kubectl get secret -n nexodus nexodus-ca-key-pair -o json | jq -r '.data."ca.crt"'
```

Update the `<nexodus_CONTROLLER_IP>` with the IP address where the Nexodus controller is reachable. If you are running the Nexodus stack in Kind cluster on your local machine, set it to the public IP of the machine.

**Note**
Current username and password are default configuration for the development environment, which is very likely to change in near future.

You can refer to [kustomization.yaml.sample](https://github.com/nexodus-io/nexodus/blob/main/deploy/nexodus-client/overlays/dev/kustomization.yaml.sample) for an example manifest.

If you have set up your Nexodus stack with a non-default configuration, please copy the [sample](./../deploy/nexodus-client/overlays/sample/) directory and update the sample file accordingly to create a new overlay for your setup and deploy it.

### Deploying the Agent in the Kind Dev Environment

If you're using the kind-based development environment, you can deploy the agent to each node in that cluster using this command:

```sh
make deploy-nexodus-agent
```

Then you may skip down to the [Verify the deployment](#verify-the-deployment) section.

Otherwise, if you are working with another cluster, continue to the next section.

### Deploying the Nexodus Agent Manifest

Once the configuration is set up, you can deploy Nexodus's manifest files.

```sh
kubectl apply -k ./deploy/nexodus-client/overlays/dev
```

It will deploy a DaemonSet in the newly created `Nexodus` namespace. DaemonSet deploys a Pod that runs a privileged container that does all the required configuration on the local node. It also starts the agent and onboard the device automatically using the non-interactive access token-based onboarding.

**Note**
If your Kubernetes cluster enforces security context to deny privileged container deployment, you need to make sure that the security policy is added to the service account `nexodus` (created for the agent deployment) to allow the deployment.

### Controlling the Agent Deployment

By default, Nexodus agent is deployed using DaemonSet, so Kubernetes will deploy Nexodus agent pod on each worker node. This might not be the ideal strategy for onboarding Kubernetes worker nodes for many reasons. You can control the Nexodus agent deployment by configuring the `nodeAffinity` in the [node_selector.yaml](https://github.com/nexodus-io/nexodus/blob/main/deploy/nexodus-client/overlays/dev/node_selector.yaml).

The default behavior is set to deploy Nexodus pod on any node that is running Linux Operating System and is tagged with `app.kubernetes.io/nexodus=`. With this deployment strategy, once you apply the Nexodus manifest, Kubernetes won't deploy Nexodus pod on any worker node. To deploy the Nexodus pod on any worker node, tag that node with `app.kubernetes.io/nexodus=` label.

```sh
kubectl label nodes <NODE_NAME> app.kubernetes.io/nexodus=
```

If you want to remove the deployment from that node, just remove the label.

```sh
kubectl label nodes <NODE_NAME> app.kubernetes.io/nexodus-
```

If you want to change the deployment strategy for Nexodus pod, please copy the [sample](./../deploy/nexodus-client/overlays/sample/) directory to create a new overlay, and configure the  [node_selector.yaml.sample](https://github.com/nexodus-io/nexodus/blob/main/deploy/nexodus-client/overlays/dev/node_selector.yaml) file as per your requirements. After making the required changes rename the file to `node_selector.yaml` and deploy it.

Currently, the sample file provides two strategies to control the deployment, but feel free to change it based on your requirements.

1 ) Deploy Nexodus pod on any node that is tagged with `app.kubernetes.io/nexodus=`

```yaml
  - key: app.kubernetes.io/nexodus
    operator: Exists
```

2 ) Deploy Nexodus pod on specific node/s in the Kubernetes cluster. Uncomment the following lines in `node_selector.yaml.sample` and add the list of the nodes.

```yaml
# Deploy Nexodus client on  specific nodes
#              - key: kubernetes.io/hostname
#                operator: In
#                values:
#                - <NODE_1>
#                - <NODE_2>
```

### Verify the deployment

Check that nexodus pod is running on the nodes

```sh
kubectl get pods -n nexodus -o wide
```

Login to one of the nodes, and you should see a wireguard interface with an address assigned. For example, on Linux:

```sh
$ ip address show wg0
161: wg0: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1420 qdisc noqueue state UNKNOWN group default qlen 1000
    link/none
    inet 100.100.0.1/32 scope global wg0
       valid_lft forever preferred_lft forever
```

### Cleanup Agent Pod From Node

Removing the Nexodus manifest will remove the Nexodus pod and clean up WireGuard configuration from the node as well.

```sh
kubectl delete -k ./deploy/nexodus-client/overlays/dev
```