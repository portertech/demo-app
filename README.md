## What is this?

A simple Golang web app that demonstrates the HTTP `GET /healthz` endpoint
pattern for Kubernetes `readinessProbe` and `livenessProbe` checks.

## How do I use it?

1. Clone this repo

   ```
   $ git clone git@github.com:portertech/dummy.git dummy
   ```

2. Compile a statically linked binary of the dummy app

   ```
   $ cd dummy
   $ CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .
   ```

   _NOTE: Installing golang is left as an exercise for the reader._

3. Build a Docker image

   ```
   $ docker build -t dummy:latest .
   ```

   _NOTE: Installing docker and general help around building Docker containers
   is left as an exercise for the reader._

4. Run the container

   ```
   $ docker run -p 8080:8080 -t dummy:latest
   ```

5. Query the web app using curl to observe it running:

   ```
   $ curl -XGET -I http://127.0.0.1:8080/healthz
   HTTP/1.1 200 OK
   Date: Wed, 18 Oct 2017 15:57:03 GMT
   Content-Length: 7
   Content-Type: text/plain; charset=utf-8

   ```

   _NOTE: determining the IP address of your running container is left as an
   exercise for the reader; hint: run `docker ps` for clues._

6. Toggle the status of the running app by posting to the `/healthz` endpoint:

   ```
   $ curl -XPOST -I http://127.0.0.1:8080/healthz
   HTTP/1.1 200 OK
   Date: Wed, 18 Oct 2017 15:58:21 GMT
   Content-Length: 0
   Content-Type: text/plain; charset=utf-8

   ```

7. Query the web app again using curl to observe the "unhealthy" status:

   ```
   $ curl -XGET -I http://127.0.0.1:8080/healthz
   HTTP/1.1 500 Internal Server Error
   Content-Type: text/plain; charset=utf-8
   X-Content-Type-Options: nosniff
   Date: Wed, 18 Oct 2017 15:58:24 GMT
   Content-Length: 10

   ```
