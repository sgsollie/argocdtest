FROM alpine/k8s:1.23.14
WORKDIR /app
COPY pijob /app/pijob
#RUN chmod +x /app/pijob.sh
CMD ["/app/pijob"]
