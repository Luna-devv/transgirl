[![](https://img.shields.io/discord/828676951023550495?color=5865F2&logo=discord&logoColor=white)](https://lunish.nl/support)
![](https://ghcr-badge.egpl.dev/Luna-devv/mellow-transgirl/latest_tag)
![](https://ghcr-badge.egpl.dev/Luna-devv/mellow-transgirl/size)

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