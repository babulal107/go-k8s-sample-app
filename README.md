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
default with a docker driver
```shell
  minikube start
```
OR 

Start minikube with hyperkit driver
```shell
  minikube start --memory=4098 --driver=hyperkit
```

```shell
  minikube status
```

## 1. Pod:
  Pod is a specification of how to run the container.

Create Pod:
Minikube, you need to load the image into its internal Docker environment:
Build Docker Image:
```shell
  docker build -t babulal107/go-k8s-sample-app:latest .
```
Load Docker image to minikube
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
  NodePort services are not directly accessible via minikube ip:<NodePort>.
  This is because the Minikube VM is running inside Docker, and it does not expose services to your host machine in the usual way.

Use minikube service Command:
    This will return a URL like: http://127.0.0.1:XXXXX
```shell
  minikube service go-k8s-app-service
  
  minikube service go-k8s-app-service --url
```

Use this URL to test your service:
```shell
  curl http://127.0.0.1:60001/health-check
```

### How to Use minikube tunnel for LoadBalancer Services:

Create Service as LoadBalancer
```shell
   kubectl apply -f k8s/service_with_load_balancer.yml
   
   kubectl get svc
```
You will see your service as LoadBalancer
Notice that EXTERNAL-IP is <pending>—this means the LoadBalancer isn't assigned an external IP yet.

```shell
  minikube tunnel
  
  kubectl get svc
```
Get the Assigned EXTERNAL-IP
Access App with external-ip like http://127.0.0.1:80/

 
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

#### If you expose Service as LoadBalancer then
If you use minikube tunnel, you should not need minikube ip, as the tunnel forwards traffic to 127.0.0.1. In this case, update /etc/hosts with:
```shell
    minikube tunnel
    
    echo "127.0.0.1 example.com" | sudo tee -a /etc/hosts
```
Then you can access your service with this url http://example.com



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
  echo -n "APP_NAME: go-k8s-app-service\nAPP_PORT: 8080\nAPP_LOG_LEVEL: debug\nAPP_ENV: dev" | base64
```

```shell
   kubectl apply -f k8s/app_config_sectet.yml
   
   OR
   
   kubectl create secret generic my-app-config-secret \
    --from-file=k8s/config/app_config.yaml --from-file=k8s/config/db_secret.yaml 
```

## Kubernetes Monitoring Using Prometheus & Grafana
More about Monitoring: https://github.com/iam-veeramalla/observability-zero-to-hero/tree/main/day-1

### 1. Prometheus:
 - Prometheus is an open-source systems monitoring and alerting toolkit.
 - It is known for its robust data model, powerful query language (PromQL), and ability to generate alerts based on collected time-series data.
 - And this data is stored in TSDB (time-series database)

### 2. Prometheus Web UI
 - The Prometheus Web UI allows users to explore the collected metrics data, run ad-hoc PromQL queries, and visualize the results directly within Prometheus.

### 3. Grafana:
 - Grafana is a powerful dashboard and visualization tool that integrates with Prometheus to provide rich, customizable visualizations of the metrics data.

## Install Prometheus & Grafana using helm:

### 1. Install helm
```shell
    brew install helm
    
    helm version
```

### Step 1: Add Helm Repositories:
Adding prometheus and grafana repo to helm charts
```shell
  helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
  helm repo add grafana https://grafana.github.io/helm-charts
  helm repo update
```
### Step 2: Install Prometheus using Helm
```shell
  helm install prometheus prometheus-community/prometheus
  
  OR
  helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace
```

You will see prometheus pods running (alert-manager, server, kube-state-metrics etc.) with ClusterIP
```shell
   kubectl get pods
   
   kubectl get svc
```

All prometheus services are exposes type of ClusterIP
Need to change prometheus-server expose type as NodePort to access outside k8s cluster
```shell
  kubectl expose service prometheus-server --type=NodePort --target-port=9090 --name=prometheus-server-ext
  
  kubectl get svc
```

```shell
  minikube ip
```
Access Prometheus: Now, access Prometheus using the Minikube IP and NodePort:
`http://192.168.49.2:30000`

#### Note : Minikube running with the Docker driver on macOS. Unlike other drivers (e.g., VirtualBox or HyperKit), the Docker driver does not expose NodePorts directly to the host machine.

#### Use minikube service Command (Recommended) minikube tunnel:
Minikube provides a built-in way to access services exposed as NodePort:

```shell
  minikube service prometheus-server-ext
```
This will automatically open the correct URL in your browser.


### Step 3: Install Grafana using Helm
```shell
  helm install grafana grafana/grafana
```
##### Notes: You need to read what notes showing here because it's have important commands like get admin user password, etc.

1. Get your 'admin' user password by running:
```shell
  kubectl get secret --namespace default grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```
2. Need to expose grafana service as NodePort as default it's running as ClusterIP
```shell
  kubectl get svc grafana
```
```shell
  kubectl expose service grafana --type=NodePort --target-port=3000 --name=grafana-ext
  
  kubectl get svc grafana-ext
```
#### Use minikube service Command (Recommended) minikube tunnel:
Minikube provides a built-in way to access services exposed as NodePort:
```shell
   minikube service grafana-ext
```
This will automatically open the correct URL in your browser grafana dashboard.

You can login grafana dashboard Username: admin and password you will get with this command
```shell
   kubectl get secret --namespace default grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```
1. You need to configure Prometheus as Data Source like host will be http://minikube-ip:port
2. Create Dashboard: you can used existing dashboard by importing like just type code 3662 and add

We Need to expose prometheus-kube-state-metrics service as NodePort as default it's running as ClusterIP
```shell
  kubectl expose service prometheus-kube-state-metrics --type=NodePort --target-port=8080 --name=prometheus-kube-state-metrics-ext
```

You can configure prometheus-kube-state-metrics server ip in configMap of Prometheus server
OR you can configure your custom go-lang app metrics also.
```shell
  kubectl get cm 
  
  kubectl edit cm prometheus-server
```
Add new jobs and set targets host with port like
```
scrape_configs:
    - job_name: prometheus
      static_configs:
      - targets:
        - localhost:9090
```