# Pod Janitor

Pod Janitor is an in-cluster process to remove succeeded pods from the Kubernetes cluster. Typically run via a Cron Job deployed in the cluster.

### Built With
- Go
- Docker

## Getting Started
To get a local instance up and running follow these steps:

### Prerequisites
- Kubernetes cluster running locally

### Installation
1. Clone the repo
```
git clone https://github.com/filetrust/pod-janitor.git
```

2. cd to root
```
cd .\pod-janitor\
```

3. Build Docker image
```
docker build -t <repository-name>/pod-janitor .
```

4. Push Docker image to repository
```
docker push <repository-name>/pod-janitor
```

5. Deploy cron job which runs image to cluster