# turbo-barnacle-chain
A p2p network that allows users to chat. In the process of integrating a simple blockchain

![turbo-barnacle-chain](./turbo-barnacle-chain.webp)


## Installation
Ensure you have Go installed on your machine. If not, make follow the instructions of the [official documentation](https://go.dev/doc/install)

1. Clone the repo either through ssh or https
```bash
# the following example is for ssh
git clone:git@github.com:enriquebarco/turbo-barnacle-chain.git
```
3. Go to the root directory and create a build
```bash
cd turbo-barnacle-chain
go build -o turbo-barnacle-chain ./cmd/blockchain-node
```

## Running the process

1. This is a basic p2p app, so connecting to peers is a bit rudimentary. You need to specify a couple of flags to make this happen and include: 
- -port <the port your node will live and the application will run>
- -name <the name your node will identify as>
- -connect<ip:port> (the IP address of the node you wish to connect to, along with the specified port the process is running)

2. You need to copy and share your IP address with the node who wants to connect with you. This can easily be done by running the following command in the terminal (for Mac)
```bash
 ip=$(ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}'); echo $ip; echo $ip | pbcopy
```

3. After cloning the repo and running the build
```bash
./turbo-barnacle-chain -port 3000 -name firstNode -connect <YourFriendsIPAddress>:3000
```

3. (Optional) If you want to test out how two nodes would connect locally, open up a new terminal and run the following bash script in the root directory of the project:
```bash
./turbo-barnacle-chain -port 3000 -name firstNode -connect localhost:3001
```
Then, open up another terminal at the root directory of the project and run:
```bash
./turbo-barnacle-chain -port 3001 -connect localhost:3000 -name secondNode
```


> If all went well, you will see the following print statement:
> ```bash
> ./turbo-barnacle-chain -port 3000 -name firstNode -connect localhost:3001
>  _____             _                   ______                                   _                _____  _             _
> |_   _|           | |                  | ___ \                                 | |              /  __ \| |           (_)
>   | | _   _  _ __ | |__    ___  ______ | |_/ /  __ _  _ __  _ __    __ _   ___ | |  ___  ______ | /  \/| |__    __ _  _  _ __> 
>   | || | | || '__|| '_ \  / _ \|______|| ___ \ / _` || '__|| '_ \  / _` | / __|| | / _ \|______|| |    | '_ \  / _` || || '_ \
>   | || |_| || |   | |_) || (_) |       | |_/ /| (_| || |   | | | || (_| || (__ | ||  __/        | \__/\| | | || (_| || || | | |
>   \_/ \__,_||_|   |_.__/  \___/        \____/  \__,_||_|   |_| |_| \__,_| \___||_| \___|         \____/|_| |_| \__,_||_||_| |_|
>
> 2024/05/16 11:26:53 Listening for P2P connections on 3000...


## p2p

The P2P network implementation involves nodes connecting via TCP connections. Each message is sent over a new TCP connection, which is closed after the message is transmitted. The app allows messaging of any kind, and only taps into the blockchain with a specified command that is discussed in more detail below.

> The application supports basic chats between the nodes. For example:
> ```bash
> secondNode [MESSAGE]: hey!
> How are you doing today?
> secondNode [MESSAGE]: Good! thanks
> ```


## Blockchain 

This simple blockchain uses Proof of Work (PoW) as its consensus mechanism. During mining, the nonce is incremented until the block's hash meets the difficulty requirement, ensuring computational effort is required to add a new block.

### Requesting and updating to the latest blockchain

When a node goes online, it will automatically request the latest blockchain in the network. When it gets the blockchain back, it will validate it and add replace the local blockchain if necessary

> Here is an example of a node updating its blockchain:
> ```bash
> 2024/05/16 11:35:02 Listening for P2P connections on 3001...
> # marshaled chain has been truncated for the example
> firstNode [RECEIVE_CHAIN]: [{"Data":null,"Hash":"0 Hello Mel","PreviousHash":"","Timestamp":"2024-05-16T15:26:53.56527Z","Nonce":0},{"Data":{"amount":10,"from":"mel","to":"kike"}...
> Received blockchain from remote node
> Blockchain replaced with the received chain
> Current Blockchain:
> Transaction: map[amount:10 from:mel to:kike], Nonce: 60
> Transaction: map[amount:120 from:jamil to:juan], Nonce: 489
> Transaction: map[amount:236 from:mel to:kike], Nonce: 197
> ```

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
># once the block is added, it will give you a confirmation such as:
>New block added to the blockchain
>Current Blockchain:
>Transaction: map[amount:10 from:kike to:mel], Nonce: 7
>Broadcasting new block to the network...
>```
>
> Once the block has been added to the local blockchain, it will broadcast the new block to the
> connected node

### Updating the blockchain based on broadcasted blocks

When a node mines a new block, it broadcasts it to the network. The node that recieves the block validates it by
- Ensuring the previous block's hash matches the new block's `previousBlockHash`: `previousBlock.Hash == newBlock.PreviousHash`
- ensuring the new block's hash is valid by recalculating the hash given the new block data `bash newBlock.Hash == newBlock.calculateHash()`

If the block is valid, it is added to the blockchain

An example of a succesful block being recieved, validated, and added would be:
>
>```bash
> firstNode: BLOCK:{"Data":{"amount":10,"from":"kike","to":"mel"},"Hash":"0014f8f216f57293931ff2d90f4748e9b08e5cbfe8414594601b85263b96dfb6","PreviousHash":"0 Hello Mel","Timestamp":"2024-05-15T03:38:59.374901Z","Nonce":7}
> New block added to the blockchain
> From: kike, To: mel, Amount: 10.000000
>Current Blockchain:
>Transaction: map[amount:10 from:kike to:mel], Nonce: 7
>```
