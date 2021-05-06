## Notes
https://aws.amazon.com/blogs/opensource/network-load-balancer-nginx-ingress-controller-eks/
https://docs.aws.amazon.com/eks/latest/userguide/load-balancing.html

## Helm
https://github.com/kubernetes/ingress-nginx/tree/master/charts/ingress-nginx

Download manifest
```
wget https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-0.32.0/deploy/static/provider/aws/deploy.yaml
```

Download resources
```
wget https://raw.githubusercontent.com/cornellanthony/nlb-nginxIngress-eks/master/apple.yaml
wget https://raw.githubusercontent.com/cornellanthony/nlb-nginxIngress-eks/master/banana.yaml
```

TSL Cert
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout \
tls.key -out tls.crt -subj "/CN=aistemos.com/O=aistemos.com"
```

Create Secret
```
kubectl create secret tls tls-secret --key tls.key --cert tls.crt
```

Download Ingress
```
wget https://raw.githubusercontent.com/cornellanthony/nlb-nginxIngress-eks/master/example-ingress.yaml
```
