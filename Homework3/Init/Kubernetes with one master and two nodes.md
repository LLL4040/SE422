# Kubernetes with one master and two nodes

## Master initialization

```shell
$ sudo kubeadm init --pod-network-cidr=192.168.0.0/16
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
$ kubectl apply -f https://docs.projectcalico.org/v3.10/manifests/calico.yaml
$ watch kubectl get pods --all-namespaces
ctrl+c退出
$ kubectl taint nodes --all node-role.kubernetes.io/master-
$ kubectl get nodes -o wide
```

## Node join

```shell
$ kubeadm join ip:port --token [TOKEN] \
     --discovery-token-ca-cert-hash xxx
```

## Set the role of node

```shell
$ kubectl label nodes [NAME] node-role.kubernetes.io/master=
$ kubectl label nodes [NAME] node-role.kubernetes.io/node=
```

