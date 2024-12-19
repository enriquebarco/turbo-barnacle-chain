# turbo-barnacle-chain
A p2p network that allows users to chat. This app has a blockchain integration that allows nodes to send transactions!

![turbo-barnacle-chain](./turbo-barnacle-chain.webp)


## Installation
Ensure you have Go installed on your machine. If not, make follow the instructions of the [official documentation](https://go.dev/doc/install)

If you are connecting to nodes outside of your local network and do not want to set up port forwarding, you will need to create a tunnel with Ngrok. Make sure you install / create an account / authenticate. All the necessary steps can be found [here](https://dashboard.ngrok.com/signup)

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

This is a basic p2p app, so connecting to peers is a bit rudimentary. You need to specify a couple of flags to make this happen and include: 
> - -port <the port your node will live and the application will run>
> - -name <the name your node will identify as>
> - -connect<ip:port> (the IP address of the node you wish to connect to, along with the specified port the process is running)

### Connecting to external nodes
A tcp connection works by specifying a private IP address and a port. However, when nodes are on different networks, it is not possible to connect to a private IP as only the public IP is exposed for security reasons. Since a public IP is actually just a router's IP, in order for this application to work, we would need to expose a port on the router to external connections and then port forward connections to the node's private IP. 

This is how Bitcoin Core software works (a router's port 8333 is setup to port forward connections to the local nodes). 

However, for the scope of this project, it is too complicated to ask users to tamper with their router settings. Therefore, we are going to use **Ngrok** - a secure unified ingress platform - to create a secure tunnel between the nodes.

Please follow these steps to create a tunnel, swap forwarding urls with your friends, and run the process to connect with them 

1. Create a tunnel 
```bash
# 3000 in this example is the port we are going to create the tunnel on
ngrok tcp 3000
```

```bash
# expected response
Full request capture now available in your browser: https://ngrok.com/r/ti

Session Status                online
Account                       Enrique Barco (Plan: Free)
Version                       3.9.0
Region                        United States (us)
Web Interface                 http://127.0.0.1:4040
Forwarding                    tcp://4.tcp.ngrok.io:15925 -> localhost:3000

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```

> Ngrok tunnel can be stopped at any time by killing the process, this should always be done at the end of a session for security reasons.

2. Share the forwarding url (which will act as an IP address) to the node who will connect with you

> Following the example above, this URL would be `4.tcp.ngrok.io:15925`. **NOTE: you NEED to remove 'tcp://' from the forwarding url** 

3. Once both nodes have created tunnels, run the application specifying the correct flags

For example:

> FirstNode:
> ```bash
> ./turbo-barnacle-chain -port 3000 -name firstNode -connect 4.tcp.ngrok.io:15925
> ```

> SecondNode:
> ```bash
> ./turbo-barnacle-chain -port 3000 -name secondNode -connect <FORWARDING_URL>
> ```

> If all went well, you will see the following print statement:
> ```bash
>  _____             _                   ______                                   _                _____  _             _
> |_   _|           | |                  | ___ \                                 | |              /  __ \| |           (_)
>   | | _   _  _ __ | |__    ___  ______ | |_/ /  __ _  _ __  _ __    __ _   ___ | |  ___  ______ | /  \/| |__    __ _  _  _ __> 
>   | || | | || '__|| '_ \  / _ \|______|| ___ \ / _` || '__|| '_ \  / _` | / __|| | / _ \|______|| |    | '_ \  / _` || || '_ \
>   | || |_| || |   | |_) || (_) |       | |_/ /| (_| || |   | | | || (_| || (__ | ||  __/        | \__/\| | | || (_| || || | | |
>   \_/ \__,_||_|   |_.__/  \___/        \____/  \__,_||_|   |_| |_| \__,_| \___||_| \___|         \____/|_| |_| \__,_||_||_| |_|
>
> 2024/05/16 11:26:53 Listening for P2P connections on 3000...


### Connecting with local nodes (on the same network)

1. You need to copy and share your private IP address with the node who wants to connect with you. This can easily be done by running the following command in the terminal (for Mac)
```bash
ip=$(ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}'); echo $ip; echo $ip | pbcopy
```

2. Once the node has shared their private IP run
```bash
./turbo-barnacle-chain -port 3000 -name firstNode -connect <YourFriendsIPAddress>:3000
```

### Testing on a local machine

If you want to test out how two nodes would connect locally, open up a new terminal and run the following bash script in the root directory of the project:
```bash
./turbo-barnacle-chain -port 3000 -name firstNode -connect localhost:3001
```
Then, open up another terminal at the root directory of the project and run:
```bash
./turbo-barnacle-chain -port 3001 -name secondNode -connect localhost:3000
```


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
