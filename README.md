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
## Start Minikube:
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
  
  If you are unable to access the application, it may be if you run on a VM with a specific networking configuration.

Try Minikube Tunnel (If Using Minikube):
  This may be required if you’re running Minikube on a VM with specific networking configurations.
  Get URL with port:
```shell
  minikube service go-k8s-app-service --url
```

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
It will run Kubeshark and redirect to a web page with a URL like http://127.0.0.1:8899/?q=%21dns+and+%21error

We can add a filter to see specific endpoint traffic route requests:
  `http and request.path == "/health-check"` and click on apply

