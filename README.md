## Go k8s Sample Application:
Just a simple GoLang application create using the Go-gin HTTP framework. Create a health-check endpoint.

`GET [http](http://localhost:8080/health-check)`

## Containerize Go Application by Docker

#### Build Docker Image:
```shell
  docker build -t babulal107/go-k8s-sample-app:latest .
```

#### Run Image will create a container 

```shell
  docker run --name=go_k8s_app -d -p 8080:8080 -it babulal107/go-k8s-sample-app
```

#### Stop running the container

```shell
  docker container stop a9cf0c3534a1
```

#### Verify Go app server running on localhost

Open your browser and type `http://localhost:8080/health-check`

### Using Docker Compose

#### Run Application by docker-compose like if we have multiple containers e.g Go App and Postgresql 
Run
```shell
  docker-compose up -d
  # OR force to build image always
  docker-compose up -d --build
```

#### Checking Logs
```shell
  docker-compose logs go_k8s_app
```

## Run Container in Kubernetes(k8s)
### Install Minikube:
  dock lin: https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download
### Start Minikube:
Once we do minikube start, your k8s cluster started
On Mac/Windows on => VM -> single node Kubernetes cluster with a default driver as docker
```shell
  minikube start
```

```shell
  minikube status
```

## 1. Pod:
  Pod is a specification of how to run the container.

Create Pod:
Minikube, you need to load the image into its internal Docker environment:
```shell
  minikube image load babulal107/go-k8s-sample-app:latest
```
```shell
  kubectl apply -f k8s/pod.yml
```
Get running pods
```shell
  kubectl get pods
```
OR
```shell
  kubectl get pods -o wide
```

Login to Kubernetes cluster and hit request to pod container
```shell
  minikube ssh
  
  curl 10.244.0.9:8080/health-check
```

- Check details about the pod:
```shell
  kubectl describe pod go-k8s-sample-app
```

- Check logs of pod:
```shell
  kubectl logs go-k8s-sample-app
```

- Delete pod:
```shell
  kubectl delete pods go-k8s-sample-app
```
  
## 2. Deployment:
It's just a wrapper that manages replicas of Pods (updating, scaling & rolling back of pods)
⚙️Deployment(yml wrapper file) -> ReplicaSet (k8s container) -> Pods
```shell
  kubectl apply -f deployment.yml
```
- Get all info
```shell
  kubectl get all
```

- Get Deployment & pods info:
```shell
  kubectl get deploy
  
  kubectl get pods
  
  kubectl get pods -o wide
```

- Get Pods info in vorticity level like 7 or 9(max)
```shell
  kubectl get pods -v=7
```


## 3. Service:
Service a wrapper on top up of Deployment. Exposes pods to internal/external networks. 
It enables communication between components or external access to applications.
Access application in-with org/network by NodePort mode or publicly by Load Balancer

```shell
  kubectl apply k8s/service.yml
```
  
- Get all Service info:
```shell
  kubectl get svc
  
  kubectl get svc -v=7
```


Login to Kubernetes cluster and hit request by service cluster IP-Address:
```shell
  minikube ssh
  curl http://10.110.66.39:80/health-check
```
Hit request through NodePort IP address
- Get Minikube node IP
```shell
  minikube ip
  
  curl http://192.168.49.2:30007/health-check
```
  If you're running Minikube with the Docker driver on macOS (Darwin), 
  NodePort services are not directly accessible via minikube ip due to how networking works with Docker.

Use minikube service Command:
    This will return a URL like: http://127.0.0.1:XXXXX
```shell
  minikube service go-k8s-app-service --url
```

Use this URL to test your service:
```shell
  curl http://127.0.0.1:60001/health-check
```
 
## Install KubeShark
  Doc link: https://docs.kubeshark.co/en/install
```shell
  brew install kubeshark  
```
### Run kubeShark: 
```shell
  kubeshark tap
```
It will run KubeShark and redirect to a web page with a URL like http://127.0.0.1:8899/?q=%21dns+and+%21error

We can add a filter to see specific endpoint traffic route requests:
  `http and request.path == "/health-check"` and click on apply

## Kubernetes Ingress:

#### Problems Load Balancer Type Service:
  1. Enterprise and TLS Load Balancing Capitalizes not support:
     (ratio-based, sticky session, TLS, HTTPS, host-base, path-base, whitelisting, blacklisting)
  2. Load Balancer type service: Cloud provider will charge you each every Load Balancer Service type for all static public ip-addresses.

  - User creates Ingress resource: allows you to create Ingress Resource (routing rules to services)
  - Deploy Ingress Controller (Load Balancer + API Gateway)
  - Load Balancing Company like (NGINX, F5) → They will write their own Ingress Controller and make it public like how to install.

### 🚀 How It Works:
  - Ingress Resource: Defines routing rules (e.g., paths like /api or domains like app.example.com).
  - NGINX Ingress Controller: Reads these rules and handles incoming traffic, forwarding it to the right service inside the cluster.

#### Install NGINX Ingress Controller (As Load Balancer)

```shell
  minikube addons enable ingress
```

#### Check Ingress addon running:
```shell
  minikube addons list | grep ingress
```

#### Verify Ingress controller running or not:
```shell
  kubectl get pods -n ingress-nginx
```

#### Check the Ingress Controller Service Type
```shell
  kubectl get svc -n ingress-nginx
```

#### Create Ingress as Resources based on defined rules in ingress.yml file.
```shell
  kubectl apply -f k8s/ingress.yml
  
  kubectl get ingress
```

#### Verify Ingress resource Describe Ingress resources:
```shell
  kubectl describe ingress ingress-example 
```

#### Configure Host in /etc/hosts file for local checking
```shell
  echo "$(minikube ip) example.com" | sudo tee -a /etc/hosts
```
OR 
`192.168.49.2 example.com`
