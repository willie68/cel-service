docker build ./ -t mcs/cel-service:V1
docker run --rm -d --restart always --name cel-service -p 9543:8443 -p 9580:8080 mcs/cel-service:V1