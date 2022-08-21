# kubectl-nodetop

This plugin displays resource (CPU/memory) usage of pods grouped by nodes

## Usage

```bash

1. Node (k3s-civo-d0a2a2a2-node-pool-c777-yu5p9)
==============================================

NAME                                     CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%   
k3s-civo-d0a2a2a2-node-pool-c777-yu5p9   180m         9%     1904Mi          50%       

Pods
====

NAMESPACE       NAME                                          CPU(cores)   MEMORY(bytes)   
database        db-0                                          27m          346Mi           
sample          sampleapi-68585cd94b-8fdtg                    1m           93Mi            
ingress-nginx   ingress-nginx-controller-59c69f6c7-v28jt      4m           83Mi            
sample          sampleapi-68585cd94b-f7kgh                    1m           46Mi            
sample          sampleapi-68585cd94b-5vjsd                    2m           41Mi            
sample          sampleapi-68585cd94b-687wr                    3m           37Mi            
kube-system     civo-csi-controller-0                         2m           36Mi            
cert-manager    cert-manager-cainjector-54f4cc6b5-vqkk6       5m           33Mi            
cert-manager    cert-manager-848f547974-5rk8c                 2m           22Mi            
kube-system     metrics-server-7bb44b587b-4jvtf               13m          18Mi            
kube-system     coredns-7796b77cd4-wllqb                      4m           14Mi            
kube-system     coredns-7796b77cd4-hvhzf                      4m           14Mi            
sample          sampleapi-68585cd94b-v2gck                    0m           14Mi            
kube-system     civo-csi-node-6wxlj                           1m           11Mi            
cert-manager    cert-manager-webhook-58fb868868-prb2b         5m           10Mi            
sample          sampleapi-68585cd94b-qh2hh                    1m           7Mi             
sample          sampleapi-68585cd94b-w65f7                    1m           7Mi             
kube-system     local-path-provisioner-84bb864455-xpqbz       1m           6Mi             
ingress-nginx   svclb-ingress-nginx-controller-bxgq4          0m           1Mi             
                                                              ________     ________        
                                                              67m          849Mi           


2. Node (k3s-civo-d0a2a2a2-node-pool-b26e-rnvk2)
==============================================

NAME                                     CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
k3s-civo-d0a2a2a2-node-pool-b26e-rnvk2   31m          3%     439Mi           47%

Pods
====

NAMESPACE       NAME                                   CPU(cores)   MEMORY(bytes)
kube-system     civo-csi-node-mz5xl                    1m           9Mi
ingress-nginx   svclb-ingress-nginx-controller-nft5n   0m           1Mi
                                                       ________     ________
                                                       1m           10Mi

3. Node (k3s-civo-d0a2a2a2-node-pool-b26e-euiq4)
==============================================

NAME                                     CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
k3s-civo-d0a2a2a2-node-pool-b26e-euiq4   36m          3%     420Mi           45%

Pods
====

NAMESPACE       NAME                                   CPU(cores)   MEMORY(bytes)
kube-system     civo-csi-node-flsgg                    1m           9Mi
ingress-nginx   svclb-ingress-nginx-controller-69wgv   0m           1Mi
                                                       ________     ________
                                                       1m           11Mi
4. Node (k3s-civo-d0a2a2a2-node-pool-b26e-h417i)
==============================================

NAME                                     CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
k3s-civo-d0a2a2a2-node-pool-b26e-h417i   34m          3%     398Mi           42%

Pods
====

NAMESPACE       NAME                                   CPU(cores)   MEMORY(bytes)
kube-system     civo-csi-node-zqc9d                    1m           9Mi
ingress-nginx   svclb-ingress-nginx-controller-mkq8m   0m           1Mi
                                                       ________     ________
                                                       1m           11Mi

```


## TODO

[] Show summary of pods instead of full details
[] 