@echo off
docker build ./ -t mcs/cel-service:V1
docker run --name cel-service -p 9543:8443 -p 9080:8080 mcs/cel-service:V1