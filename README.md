### Dockerization

1. Build Image: `docker build -t recipe-api .`

2. Run Container: `docker run --env-file .env -p 8080:8080 --network host recipe-api` 