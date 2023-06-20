## Rabbit NFT
An NFT market.

![demo](demo.gif)

## prerequisite env
* docker
* golang
* node
* yarn

## development
* install hardhat
  ```
  cd hardhat
  yarn install
  ```
* install ui dependency
  ```
  cd ui
  yarn install
  ```
* start hardhat and deploy contract
  ```
  make hh-node
  make hh-deploy
  ``` 
* start ipfs
  ```
  make start-ipfs
  ```
* start server
  ```
  make start-server
  ```
* start ui
  ```
  make start-ui
  ```
* visit: http://localhost:8080
