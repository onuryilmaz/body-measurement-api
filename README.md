# Body Measurement API

* This repository provides an API to the outside to which health apps can connect and upload measured data about the users. 
* A body measurement consists of a few parameters like type and value. Examples are a one-time measurement of blood pressure or body temperature.
* HTTP REST API is provided with the following methods and parameters:

| Verb 	| Endpoint                               	| Comment                                                                                                    	|
|------	|----------------------------------------	|------------------------------------------------------------------------------------------------------------	|
| GET  	| `/api/filter/:user/:type/:from/:to` 	| Filter measurement data for user `user`, measurement type `type` and time frame between `from` and `to` 	|
| GET  	| `/api/last/:user/:type`             	| Get last measurement data for user `user` and measurement type `type`                                   	|
| GET  	| `/api/record/:user/:type/:value`    	| Record measurement for user `user` and measurement type `type` with the value of  `value`               	|
| POST 	| `/api/save`                            	| Save measurement data sent as a body of request:                                                           	|

```
{  
   "ID": 1,
   "Type": "testType",
   "Value": 1.1,
   "UserID": "testUser",
   "Timestamp": "2017-08-31T11:56:47.286854389+03:00"
}
```

* Build and packaging steps are done inside Docker and as a result Docker image is delivered.

## Assumptions
* Authentication, encryption and user-management are assumed to be handled by other components.

## Requirements
* Docker (> version 17.05)
* GNU make

## Test
```
make test
```

## Build
```
make build
```
## Run
```
make run
```

## Push
```
make push DOCKER_REGISTRY=$REGISTRY
```

## Example Flow
```
$ make build
$ make run

# on another shell
$ curl -i localhost:8080/api/record/onuryilmaz/weight/80
HTTP/1.1 200 OK
..
$ curl -i localhost:8080/api/record/onuryilmaz/weight/81
HTTP/1.1 200 OK
..

$ curl localhost:8080/api/last/onuryilmaz/weight
{"ID":2,"Type":"weight","Value":81,"UserID":"onuryilmaz","Timestamp":"2017-09-04T07:48:30.758669975Z"}

$ curl localhost:8080/api/filter/onuryilmaz/weight/1501545600/1533081600
 [ 
  {"ID":1,"Type":"weight","Value":80,"UserID":"onuryilmaz","Timestamp":"2017-09-04T07:47:51.999662761Z"},
  {"ID":2,"Type":"weight","Value":81,"UserID":"onuryilmaz","Timestamp":"2017-09-04T07:48:30.758669975Z"}
 ]
```

## Dependency Management
* [govendor](https://github.com/kardianos/govendor) is used for dependency management.
* Fixed versions can be checked from [vendor.json](vendor/vendor.json)

