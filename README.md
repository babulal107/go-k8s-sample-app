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

## ConfigMap & Secrets
  - ConfigMaps and Secrets are used to manage configuration data separately from application code.
  - ConfigMaps: Store non-sensitive configuration data (app-name, log-level, app-port, db-port etc). →  plain text
  - Secrets: Store sensitive data (passwords, tokens, keys). → Base64-encoded

  Problem: if configuration is updated then, It didn't reflect in env if used an ENV type. You need to re-create containers that is not a good way.
  Solution: using VolumeMounts: you will do the same thing, but instant of ENV used files here as we are doing mounting.

#### When to use mountPath and subPath?

1. MountPath: Mounts the entire volume
   - The entire volume will be available at this location.
   - If the volume contains multiple files, they will all be accessible inside the mountPath.
   
   Effect/Example: 
    - The entire app-config-secret secret will be mounted at /opt/config.
    - If app-config-secret contains multiple files, all of them will be in /opt/config/. like app_secret.yaml & db_secret_yaml etc.

2. SubPath: 
    - Mounts only a single file or directory.
    - If you need only one file
   
    Effect/Example: Instead of mounting the entire secret as /opt/config/, only db_secret.yaml is mounted at /opt/config/db_secret.yaml.

#### 🚀 Which One Should You Use?
  ✅ Use mountPath alone if you want all files from the volume.

  ✅ Use subPath if you only need a specific file from the volume.

### 1. ConfigMaps:
```shell
  kubectl apply -f k8s/cm.yml
```

Get all ConfigMaps: 
```shell
  kubectl get cm
```

Describe ConfigMaps get all details:
```shell
  kubectl describe cm test-config
```

Check inside your running pod:
```shell
  kubect get pods
  
  kubectl exec -it go-k8s-sample-app-668577d89d-pcznr -- ls -l /opt  
  
  OR 
  
  kubectl exec -it go-k8s-sample-app-668577d89d-pcznr -- /bin/sh 
```

### 2. Secrets:

```shell
  kubectl apply -f k8s/db_secret.yml  
  
  kubectl get secret
  
  kubectl describe secret db-secret

```
You can check your application pod logs, you will see as you logged loaded configs


#### A). Secret with mounting volume as file

##### How to Generate the Base64 Content
 If you want to encode the values yourself, run:
```shell
  echo -n "DB_NAME: test_db\nDB_USERNAME: root\nDB_PASSWORD: user@123\nDB_PORT: 3306" | base64
```
Output look like: 
```
RE5BTUU6IHRlc3RfZGIKREJfUE9SVDogMzMwNgo=
```

Check db_config.yaml file generated or not
```shell
  kubectl get pods
  
  kubectl exec -it go-k8s-sample-app-69bffb7755-bppnf -- ls -l /opt/config
```

### MountVolumes with multiple config files

db_secret.yaml
```shell
  echo -n "DB_NAME: test_db\nDB_USERNAME: root\nDB_PASSWORD: user@123\nDB_PORT: 3306" | base64
```

app_secret.yaml
```shell
  echo -n "APP_PORT: 8080\nAPP_LOG_LEVEL: root\nAPP_ENV: dev" | base64
```

```shell
   kubectl apply -f k8s/app_config_sectet.yml
   
   OR
   
   kubectl create secret generic my-app-config-secret \
    --from-file=k8s/config/app_config.yaml --from-file=k8s/config/db_secret.yaml 
```

