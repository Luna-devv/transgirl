[![](https://img.shields.io/discord/828676951023550495?color=5865F2&logo=discord&logoColor=white)](https://lunish.nl/support)
![](https://ghcr-badge.egpl.dev/luna-devv/mellow-transgirl/latest_tag)
![](https://ghcr-badge.egpl.dev/luna-devv/mellow-transgirl/size)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I3I6AFVAP)

## About
mellow-transgirl allows you to get random files from any AWS S3 compatible storage bucket using a lightweight, scalable go http api. This project is a part of the Wamellow project.

An example of this project in deployment can be found at [transgirl.wamellow.com](https://transgirl.wamellow.com).

If you need help deploying this api, join **[our Discord Server](https://discord.com/invite/yYd6YKHQZH)**.

## Deploy
To deploy this project, create the following `docker-compose.yml`:
```yml
services:
  app:
    image: ghcr.io/luna-devv/mellow-transgirl:latest
    container_name: mw-transgirl
    ports:
      - "8080:8080"
    environment:
      AWS_REGION: auto
      AWS_BUCKET: wamellow
      AWS_ACCESS_KEY_ID: <your-access-key-id>
      AWS_SECRET_ACCESS_KEY: <your-secret-access-key>
      AWS_ENDPOINT: https://xxx.r2.cloudflarestorage.com
      AWS_PUBLIC_URL: https://r2.wamellow.com
      FILE_PREFIX: blahaj/
    restart: unless-stopped
```

To deploy the project, run:
```sh
docker compose up -d
```

## Usage
To get a random file from the bucket, make a [GET request to the `/`](https://transgirl.wamellow.com) endpoint.
You can also specify a prefix by setting the `FILE_PREFIX` environment variable, when this is set, the api will only return files that start with the specified prefix (i.e.: [`blahaj/GqpEBx.webp`](https://r2.wamellow.com/blahaj/wAuI4n.webp)).

```sh
curl http://localhost:8080
{"url":"https://r2.wamellow.com/blahaj/wAuI4n.webp"}
```

To see how many files are currently in the cache (any file from the bucket that starts with `FILE_PREFIX`), make a [GET request to the `/stats`](https://transgirl.wamellow.com/stats) endpoint.

```sh
curl http://localhost:8080/stats
{"file_count": 242}
```

To refresh the cache, make a [POST request to the `/refresh`](https://transgirl.wamellow.com/refresh) endpoint.
As authorization header, you need to provide the `Bearer` token that you set in the `AWS_SECRET_ACCESS_KEY` environment variable.

```sh
curl -X POST http://localhost:8080/refresh \
    -H "Authorization: Bearer <your-secret-access-key>"
```