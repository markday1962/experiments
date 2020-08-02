# Running tcpdump in an EKS cluster
As a pod is its own subnet to watch the network activity in a pod, we create a
container and install tcpdump, rather than create a docker file, this can be done
from the command line (shown below), once the container has been create it is pushed
into the ECR and can be added to a minifest of the pod we want to inspect.

Comand line example
```
docker build -t 111177312954.dkr.ecr.eu-west-1.amazonaws.com/aistemos/dev/tcpdump:0.2.0 - <<EOF
FROM ubuntu
RUN apt-get update && apt-get install -y tcpdump net-tools
EOF
```

Login to the AWS ECR
```
$(aws ecr get-login --region eu-west-1 --no-include-email)
```

Push the image
```
docker push -t 111177312954.dkr.ecr.eu-west-1.amazonaws.com/aistemos/dev/tcpdump:0.2.0
```

## Usage
kubectl exec -it -c tcpdump <pod> -- sh
