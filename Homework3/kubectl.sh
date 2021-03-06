sudo apt-get update && apt-get install -y apt-transport-https
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list \
   deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main \
   EOF
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl