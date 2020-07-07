FROM alpine
ADD leads-service /leads-service
ENTRYPOINT [ "/leads-service" ]
