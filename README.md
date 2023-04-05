# SafeFetchApi
###### Discovered by: [Milkyway](https://twitter.com/milkyontheblock) and Sam
###### Written by [Prizzle](https://twitter.com/bypassedpx)

## Description
The Api uses a Google Service to preview websites over their end.

https://docs.google.com/gview?url=https://google.com

It can be used to bypass various Bot Protection Systems. Its also effective agains Ratelimits.

It has some character limits, so it might not work for some websites. I didnt test them yet.

Google doesnt seem to rate limit the requests, so you can use it as much as you want.

## Usage
JS
```shell
git clone https://github.com/Prizzledizle/SafeFetchApi.git
cd js
npm install
npm start
```

Go
```shell
git clone https://github.com/Prizzledizle/SafeFetchApi.git
cd go
go run src/main.go
```

Request
```
URL TO POST: "http:localhost:4501/api/v2/safeFetch

BODY: {
    "url": "https://example.com"
}
```

If you have any questions, feel free to contact me on Discord: Prizzle#4655 or on [Twitter](https://twitter.com/bypassedpx)
