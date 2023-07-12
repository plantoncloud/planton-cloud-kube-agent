FROM golang
COPY build/app-linux /app
RUN chmod +x /app
CMD ["/app"]
