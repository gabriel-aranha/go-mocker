[![Build](https://github.com/gabriel-aranha/minimock/actions/workflows/build.yml/badge.svg)](https://github.com/gabriel-aranha/minimock/actions/workflows/build.yml)

# Minimock
Minimock is a simple Mock HTTP Server that is deployable completely for free on Heroku. Its use case is if you want some free customizable API for GET requests in early stage projects. Basically, you send a POST request to an endpoint with the data to be received, and the GET endpoint for the same endpoint will return the data. It also supports Bearer Token authentication on the POST method, so you can be sure no one will modify your data.
## Quickstart with Heroku
### Requirements
1. Create an app on Heroku
2. Download or fork the Minimock project and add it to the app
3. Add a Redis instance to the app such as Heroku Redis

### Setup
1. Go to the settings pane of your app and set the following config vars:    
    1.1 Set the Go Version config var:  
    ```
    GOVERSION -> 1.16
    ```
    1.2 Set the Redis URL and Password config vars:  
    ```
    REDIS_PASS -> yourpassword
    REDIS_URL  -> address:port
    ```
    1.3 (Optional) Set the POST API auth key. If set, the POST method needs the auth key as a bearer token:  
    ```
    AUTH_KEY -> yourauthkey
    ```

## Quickstart with Docker Compose(Minimock + Redis)
### Setup
1. Clone this repository:
    ```
    $ git clone git@github.com:gabriel-aranha/minimock.git
    ```
2. Enter the folder:
    ```
    $ cd minimock
    ```
3. Send the docker compose command:
    ```
    $ docker-compose up --build
    ```

## Quickstart without Docker Compose(Standalone Minimock)
### Requirements
1. A Redis Server instance is needed

### Setup
1. Clone this repository:
    ```
    $ git clone git@github.com:gabriel-aranha/minimock.git
    ```
2. Enter the folder:
    ```
    $ cd minimock
    ```
3. Set the Redis environment variables if needed. If unset, will default to url localhost:6379 and no password):
    ```
    # Example Redis Url and Password
    REDIS_URL="address:6379"
    REDIS_PASS="yourpassword"
    ```

### Running
1. Command to build Docker image:
    ```
    $ docker build -t minimock .
    ```
2. Command to run Docker container:
    ```
    $ docker run --name minimock -d -p 0.0.0.0:1323:1323 --restart unless-stopped --network=host minimock
    ```

## Usage
1. First, check if Minimock is running correctly by pinging its health endpoint:
    ```
    $ curl -X GET http://<address>:<port>/
    ```
2. It should return the following:
    ```
    {
        "message": "OK"
    }
    ```

3. To begin using Minimock, you need to first send a POST request to an "/api" endpoint of your choice.  
As an example, let's create the following mock endpoint:
    ```
    http://<address>:<port>/api/get-countries
    ```
4. To accomplish this, you send a POST request to the above endpoint and include in the request body your JSON response data formatted like this:
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
5. If the request is successful, you should receive OK response:
    ```
    {
        "message": "OK"
    }
    ```
6. Now you can send a GET request to the same endpoint and receive the following response:
    ```
    {
        "countries": [
            "Brazil",
            "United States",
            "France"
        ]
    }
    ```

7. This above example was for GET requests without JSON body. However, you can also send a body to the POST request and receive an unique GET response.
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

8. If the above POST request was successful, you can then GET the data by adding the following JSON body in the request:
    ```
    {
        "continent": "South America"
    }
    ```
9. And you will receive:
    ```
    {
        "countries": [
            "Brazil"
        ]
    }
    ```
