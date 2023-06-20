# hardhat node
.PHONY: hh-node
hh-node:
	cd hardhat && npx hardhat node

# hardhat deploy
.PHONY: hh-deploy
hh-deploy:
	cd hardhat && npx hardhat run scripts/rabbit-deploy.js --network localhost

# contract compile, genrate abi and go model
.PHONY: compile
compile:
	rm -rf hardhat/artifacts
	cd hardhat && npx hardhat compile
	jq '.abi' hardhat/artifacts/contracts/RabbitCollectible.sol/RabbitCollectible.json  > server/abi/RabbitCollectable.abi
	rm -rf server/contract
	mkdir -p server/contract
	abigen --abi server/abi/RabbitCollectable.abi --pkg contract --out server/contract/rabbitcollectible.go

# private blockchain genesis
.PHONY: pr-genesis
pr-genesis:
	geth --datadir=gethdata init genesis/genesis.json

# private node
.PHONY: pr-node
pr-node:
	geth --datadir=gethdata/ --networkid 1874 \
		--mine --miner.threads=1 --miner.etherbase=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 \
		--http --http.api 'web3,eth,net,debug,personal' --http.corsdomain '*' \
		--rpc.allow-unprotected-txs \
		--gcmode archive

# private node attach
.PHONY: pr-attach
pr-attach:
	geth attach gethdata/geth.ipc

# private deploy
.PHONY: pr-deploy
pr-deploy:
	cd hardhat && npx hardhat run scripts/private-deploy.js --network private

# clean private chain
.PHONY: pr-clean
pr-clean:
	rm -rf gethdata

# cluster deploy
.PHONY: clu-deploy
clu-deploy:
	cd hardhat && npx hardhat run scripts/private-deploy.js --network cluster

# start ipfs
.PHONY: start-ipfs
start-ipfs:
	cd ipfs && docker-compose up -d

# start server
.PHONY: start-server
start-server:
	cd server && go run .

# start ui
.PHONY: start-ui
start-ui:
	cd ui && yarn serve