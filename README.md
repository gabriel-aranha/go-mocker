# go-mocker

## Quickstart with Docker Compose(go-mocker + Redis)
Clone this repository:
```
$ git clone git@github.com:gabriel-aranha/go-mocker.git
```
Enter the folder:
```
$ cd go-mocker
```
Send the docker compose command:
```
$ docker-compose up --build
```

## Quickstart without Docker Compose(Standalone go-mocker)
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
Edit the Dockerfile to set the Redis host and password (If unchanged, will default to  host localhost:6379 and no password):
```
# Example Redis Host and Password
ENV REDIS_HOST="192.168.0.1:6379"
ENV REDIS_PASS="pass123"
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
The go-mocker server can then be accessed from other systems using the following endpoint:
```
<ip-on-local-network>:1323
```

Finally, check if go-mocker is running correctly by pinging its health endpoint:
```
$ curl -X GET http://<ip-on-local-network>:1323/health
```
It should return the following:
```
{
    "status": "ok"
}
```

## Usage
To begin using go-mocker, you need to first send a POST request to an "/api" endpoint of your choice.  
As an example, let's create the following mock endpoint:
```
http://<ip-on-local-network>:1323/api/get-countries
```
To accomplish this, you send a POST request to the above endpoint and include in the request body your JSON response data formatted like this:
```
{
    "response": {
        "countries": [
            "Brazil",
            "United States",
            "France"
        ]
    }
}
```
If the request is successful, you should receive an unique id response:
```
{
    "id":"1f1eee663738854c4e53bf7be7902de982f22255"
}
```
Now you can send a GET request to the same endpoint and receive the following response:
```
{
    "countries": [
        "Brazil",
        "United States",
        "France"
    ]
}
```

This above example was for GET requests without JSON body. However, you can also send a body to the POST request and receive an unique GET response.
```
{
    "body": {
        "continent": "South America"
    },
    "response": {
        "countries": [
            "Brazil"
        ]
    }
}
```

If the above POST request was successful, you can then GET the data by adding the following JSON body in the request:
```
{
    "continent": "South America"
}
```
And you will receive:
```
{
    "countries": [
        "Brazil"
    ]
}
```