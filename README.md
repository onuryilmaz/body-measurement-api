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

## Components
* REST API Server
Methods, parameters table

* Data Layer

## Requirements
* Docker (> version 17.05)
* GNU make

## Build & Push
```
make build
make push
```

## Run
```
make run
```

## Test
```
make test
```