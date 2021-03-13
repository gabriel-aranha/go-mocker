[![build](https://github.com/gabriel-aranha/go-mocker/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/gabriel-aranha/go-mocker/actions/workflows/go.yml)

# go-mocker
## Quickstart with Heroku
### Requirements
1. Create an app on Heroku
2. Download or fork the go-mocker project and add it to the app
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
    1.3 (Optional) Set the API auth key:  
    ```
    AUTH_KEY -> yourauthkey
    ```

## Quickstart with Docker Compose(go-mocker + Redis)
### Setup
1. Clone this repository:
    ```
    $ git clone git@github.com:gabriel-aranha/go-mocker.git
    ```
2. Enter the folder:
    ```
    $ cd go-mocker
    ```
3. Send the docker compose command:
    ```
    $ docker-compose up --build
    ```

## Quickstart without Docker Compose(Standalone go-mocker)
### Requirements
1. A Redis Server instance is needed

### Setup
1. Clone this repository:
    ```
    $ git clone git@github.com:gabriel-aranha/go-mocker.git
    ```
2. Enter the folder:
    ```
    $ cd go-mocker
    ```
3. Edit the Dockerfile to set the Redis url and password (If unchanged, will default to url localhost:6379 and no password):
    ```
    # Example Redis Url and Password
    ENV REDIS_URL="address:6379"
    ENV REDIS_PASS="yourpassword"
    ```

### Running
1. Command to build Docker image:
    ```
    $ docker build -t go-mocker .
    ```
2. Command to run Docker container:
    ```
    $ docker run --name go-mocker -d -p 0.0.0.0:1323:1323 --restart unless-stopped --network=host go-mocker
    ```

## Usage
1. First, check if go-mocker is running correctly by pinging its health endpoint:
    ```
    $ curl -X GET http://<address>:<port>/
    ```
2. It should return the following:
    ```
    {
        "status": "ok"
    }
    ```

3. To begin using go-mocker, you need to first send a POST request to an "/api" endpoint of your choice.  
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
5. If the request is successful, you should receive an unique id response:
    ```
    {
        "id":"1f1eee663738854c4e53bf7be7902de982f22255"
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
