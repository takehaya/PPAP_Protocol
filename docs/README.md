# docs memo

## gen sequence image
https://www.websequencediagrams.com/
```
title PPAP Protocol

Alice->Bob: PPAP Request
Bob->Alice: Authentication Response. (I have 1)
Bob->Alice: Authentication Response. (I have 2)
Alice->Bob: Ack. (Ah!)
Alice->Bob: Key Gen from the response. (<have 2>-<have 1>!)

Bob->Carol: Response Proxy. (<have 2>-<have 1>!)
Carol->Bob: Authentication Response. (I have 3)
Carol->Bob: Authentication Response. (I have 4)
Bob->Carol: Ack. (Ah!)
Bob->Carol: Key Gen from the response. (<have 3>-<have 4>-<have 2>-<have 1>!)
Carol->Bob: Close. (Pico)
Bob->Alice: Close. (Pico)
```
