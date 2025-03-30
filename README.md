[![](https://img.shields.io/discord/828676951023550495?color=5865F2&logo=discord&logoColor=white)](https://lunish.nl/support)
![](https://ghcr-badge.egpl.dev/luna-devv/transgirl/latest_tag)
![](https://ghcr-badge.egpl.dev/luna-devv/transgirl/size)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I3I6AFVAP)

## About
transgirl allows you to get random files from any AWS S3 compatible storage bucket using a lightweight, scalable go http api.

This project is a part of the [Wamellow](https://wamellow.com) project for serving random images of Blåhajs for `/blahaj`.
An example of this project in deployment can be found at [transgirl.wamellow.com](https://transgirl.wamellow.com).

If you need help deploying with this api, join **[our Discord Server](https://discord.com/invite/yYd6YKHQZH)**.

## Deploy
To deploy this project, create the following `docker-compose.yml`:
```yml
services:
  app:
    image: ghcr.io/luna-devv/transgirl:latest
    container_name: transgirl
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

---
### Statistics

To see how many files are currently in the cache (any file from the bucket that starts with `FILE_PREFIX`), make a [GET request to the `/stats`](https://transgirl.wamellow.com/stats) endpoint.

```sh
curl http://localhost:8080/stats
{"file_count": 242}
```

---
### Refreshing

To refresh the cache, make a [POST request to the `/refresh`](https://transgirl.wamellow.com/refresh) endpoint.
As authorization header, you need to provide the `Bearer` token that you set in the `AWS_SECRET_ACCESS_KEY` environment variable.

```sh
curl -X POST http://localhost:8080/refresh \
    -H "Authorization: Bearer <your-secret-access-key>"
```

## Examples
Here are some Copy-Paste examples for using this api in your projects.

---
### TypeScript
A simple example of how you can use this api in a TypeScript project.

`blahaj.ts`
```ts
export interface BlahajResponse {
    url: string;
}

export async function getBlahaj() {
    const blahaj = await fetch(this.url)
        .then((r) => r.json())
        .catch(() => null) as BlahajResponse | null;

    return blahaj;
}
```
`index.ts`
```ts
import { getBlahaj } from './blahaj.ts';

const blahaj = await getBlahaj();
if (!blahaj) process.exit();

console.log(blahaj.url);
```

---
### JavaScript Discord Bot
To use this snippet, you need to do the following things first:
1. Install the `discord.js` package using `npm install discord.js`.
2. Replace `your-token` with your bot token. (https://discord.com/developers/applications)
3. Run `node .`
4. Invite the bot to your server and run `!blahaj`.

```js
const { Client, GatewayIntent } = require('discord.js');

const client = new Client({
    intents: [
        GatewayIntent.GuildMessages,
        GatewayIntent.MessageContent,
    ]
});

client.on('messageCreate', async (message) => {
    if (message.content === '!blahaj') {
        const blahaj = await fetch('http://localhost:8080')
            .then((r) => r.json())
            .catch(() => null);

        if (!blahaj) {
            message.channel.send({ content: 'Failed to get a random blahaj' });
            return;
        }

        message.channel.send({ content: blahaj.url });
    }
});

client.login('your-token');
```

## Used in
This API is used for:
- [Wamellow Daily Posts](https://wamellow.com/docs/dailyposts) - Discord Bot to get daily Blåhajs.
- [Wamellow](https://wamellow.com/invite) - Discord Bot to get random Blåhajs by using `/blahaj`.
- [blahaj.4k.pics](https://blahaj.4k.pics) - Web view for this API.
- [@blahaj.4k.pics](https://bsky.app/profile/blahaj.4k.pics) - Daily curated Blåhajs on Bluesky.