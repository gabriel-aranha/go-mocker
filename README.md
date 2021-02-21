# go-mocker

## Docker project to run a Go Mocker Server for use by other systems on the same network

## Quickstart

### Requirements
A Redis Server instance is needed

### Setup
Clone this repository:
```
$ git clone git@github.com:gabriel-aranha/go-mocker.git
```
Enter the folder:
```
$ cd go-mocker
```
Edit the Dockerfile to set the Redis host and password (If unchanged, will default to localhost:6379 and no password):
```
$ nano Dockerfile
```

### Running
Command to build Docker image:
```
$ docker build -t go-mocker .
```
Command to run Docker container:
```
$ docker run --name go-mocker -d -p 0.0.0.0:1323:1323 --restart unless-stopped --network=host go-mocker
```
