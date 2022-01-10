# cel-service
HTTP JSON and gRPC service for evaluating an expression using google cel (https://opensource.google/projects/cel)

This service is a small, fast go microservice based on go-micro (https://github.com/willie68/go-micro) a template for a simple go microservice.

Features:

- gRPC and HTTP JSON APIs
- simple API to evaluate an expression
- go 1.17

## Installation

As this service is not part of the docker hub, you have to build up your docker image yourself.

For this simply start then docker build 

`docker build ./ -t mcs/cel-service:V1`

The docker file is a 2 phase docker build. First building the binaries, than building the runnable image.

You can start the image with

`docker run --name cel-service -p 8443 -p 8080 mcs/cel-service:V1` 

 

## Example http json

in the api folder you will find a postman collection for this service using the http json interface, as the all the health and metrics endpoints.

A simple curl example for the service:

```sh
curl --location --request POST 'https://127.0.0.1:9543/api/v1/evaluate' \
--header 'apikey: 8723a34c54a53c70071cf86dfb1d8744' \
--header 'Content-Type: application/json' \
--data-raw '{"context": {"data": {"index": 1}},"expression": "data.index == 1"}'
```

The result should be something like this:

```json
{
  "error": null,
  "message": "result ok: true",
  "result": true
}
```



In case of an error in your evaluation you will get an special error response: 

for an eval of:

```json
{
  "context": {
      "data": {
          "index": 1
      }
  },
  "expression": "data.index == 1.0"
} 
```



```json
{
  "error": "no such overload",
  "message": "program evaluation error: no such overload",
  "result": false
}
```

because you can't compare an float (in the expression) with an int literal. (1 in cel is an int literal)

see the cel project for further information. (https://opensource.google/projects/cel)

## Expression Cache

The service has implemented an expression cache. Most time consuming operations are the parameter analyzing and the expression program compiling. The result of this two steps can be cached, so that you can reuse the same expression program with different contexts. The context definition should be equal, the values of course can be changed. To cache an expression simply add an identifier to the request:

```json
{
  "context": {
      "data": {
          "index": 1
      }
  },
  "expression": "data.index == 1.0",
  "identifier": "12345"
} 
```

Every subsequent call with the same **id** and **expression** will be accessing the cached expression program. For updating simply change the expression, that the cache will automatically updated with the newly compiled expression program. 

## Example gPRC

The service also expose a grpc server (with the default port 50051 with TSL). The definition of the service and the models you can find in the api folder (cel-service.proto)

Be aware, because of https://github.com/google/cel-go/issues/203 in gRPC numeric parameters in the context are always interpreted as float. In the example above  the expression should that be 

int(data.index) == 1

The problem here is that you can't use the same expression for both HTTP JSON and gRPC. 
