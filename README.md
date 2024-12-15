<div align="center">
<img
    width=100%
    src="banner.png"
    alt="massolit banner"
/>
</div>

# Massolit - Managebac Extension

> **This app is entirely separate from ManageBac and is not an official product.**

This little project was made to scratch my own itch, your milage may vary.

**While the app is being used in the real world, it's internal tests are, lets say, lacking.**

# Features

## Book Tracker

Massolit is able to keep track of which books are borrowed by which student.

Currently, it only works if with QR codes which must be generated with a specially formatted code.

```
MASSOLIT|1|12345678|IG-PSYCH-20
MASSOLIT|<version>|<isbn>|<unique book id>
```

Upon scanning, Massolit tries to find the book and the book instance. If it does not exist, it will add it to the database.

Best effort is made to fetch the title and cover image of the provided ISBN.

## Detention Tracker

Massolit also keeps track of ManageBac behaviour notes and sends daily email reports with students who have detention that day.

Behaviour note types must include the word `Detention` to be detected.

Emails are sent at 13:00 every work day. Currently, this is hardcoded and not configurable.

# Production

## Build

As this project is written in Go with the webapp embedded, it builds to a single executable. This can be run anywhere.

### Docker

No docker image is provided, but it would be trivial to do this <- Recommended

## Deploy

This one is up to you. Interal servers or a service like Hetzner are both good choices.

# Developing and Contributing

## Run Locally

This will run both a Nuxt3 dev server with hot reload and a pocketbase instance with a authentication and an api ready to go.

Clone the project

```bash
  git clone https://github.com/veritymedia/massolit
```

Go to the project directory

```bash
  cd massolit
```

Install dependencies

```bash
  pnpm install
```

Start the server

```bash
  pnpm dev // starts full app
  # pnpm dev:nuxt // only starts frontend
  # pnpm dev:pocketbase // only starts backend
```

| URL                 | Function                       |
| ------------------- | ------------------------------ |
| localhost:8090/\_/  | pocketbase admin setup and log |
| localhost:8090/api/ | pocketbase api                 |
| localhost:3000/     | nuxt3 dev server               |

From there build your nuxt app as normal. Follow the [pocketbase docs](https://pocketbase.io/docs/) for more info on how to use pocketbase.

## Deployment

This will create a single binary containing Nuxt and PocketBase for deployment.

```bash
  yarn build:prod
```

```bash
  ./pocketnuxt serve --http "yourdomain.com:80" --https "yourdomain.com:443"
```

https://pocketbase.io/docs/going-to-production/ for more examples.

## Acknowledgements

- [PocketBase](https://github.com/pocketbase/pocketbase)
- [PocketNuxt](https://github.com/j-wil/pocket-nuxt)
