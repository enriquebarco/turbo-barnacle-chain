# turbo-barnacle-chain
A p2p network that allows users to chat. In the process of integrating a simple blockchain

![turbo-barnacle-chain](./turbo-barnacle-chain.webp)


## Installation
Ensure you have Go installed on your machine. If not, make follow the instructions of the [official documentation](https://go.dev/doc/install)

## Running the process

1. This is a basic p2p app, so connecting to peers is a bit rudimentary. You need to specify a couple of flags to make this happen and include: 
- -port <the port your node will live and the application will run>
- -name <the name your node will identify as>
- -connect<ip:port> (the IP address of the node you wish to connect to, along with the specified port the process is running)

2. You need to copy and share your IP address with the node who wants to connect with you. This can easily be done by running the following command in the terminal (for Mac)
```bash
 ip=$(ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}'); echo $ip; echo $ip | pbcopy
```

3. After cloning the repo, run the main function
```bash
cd cmd/blockchain-node/
go run . -port 3000 -name firstNode -connect <YourFriendsIPAddress>:3000
```

3. (Optional) If you want to test out how two nodes would connect locally, open up a new terminal and run the following bash script in the root directory of the project:
```bash
cd cmd/blockchain-node/
go run . -port 3000 -name firstNode -connect localhost:3001
```
Then, open up another terminal at the root directory of the project and run:
```bash
cd cmd/blockchain-node/
go run . -port 3001 -connect localhost:3000 -name secondNode
```

## Blockchain 

This simple blockchain uses Proof of Work (PoW) as its consensus mechanism. During mining, the nonce is incremented until the block's hash meets the difficulty requirement, ensuring computational effort is required to add a new block, thus enhancing security.

### Creating and sending blocks

To create a transaction, which generates a new block, a user must specify the following syntax:

> send
>  - personSendingTransaction,
>  - personReceivingTransaction,
>  - amount
>
> for example, the following would be a valid transaction:
>
> 
>```bash
>send kike,mel,10
>
>// once the block is added, it will give you a confirmation such as:
>New block added to the blockchain
>From: kike, To: mel, Amount: 10.000000
>Current Blockchain:
>Data: map[amount:10 from:kike to:mel]
>```
