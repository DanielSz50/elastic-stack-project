# elastic-stack-project
Assignment for the NoSQL databases course.

# Run elastic stack and gin server
Run the following command from the project root directory:
```sh
docker-compose up
```

# Add index templates
Wait for services to become healthy and run the script:
```sh
./script/put_templates.sh
```

# Generate logs
Logstash is configured with pipeline to gather data from filebeat. Filebeat tracks logs from HTTP server running in the *app* service.
Run the following command to do some requests, the logs will be pushed to elastic:
```sh
./script/do_app_requests.sh
```
