# CQL is a **Byzantine Fault Tolerant** distributed relational database

Features:
- **ServerLess**: Free, High Availabile, Auto Sync Database Service for Serverless App
- **SQL**: Most SQL-92 support.
- **Decentralize**: Running on Open Internet without Central Coordination.
- **Privacy**: Access with Granted Permission and Encryption Pass.
- **Immutable**: Query History in CQL is Immutable and Trackable.
- **Permission**: Column Level ACL and SQL Pattern Whitelist.

## What is CQL?

- Open source alternative of Amazon QLDB
- Low cost DBaaS
- Just like filecoin + IPFS is the decentralized file system, CQL is the decentralized database

CQL is a derived work of [CovenantSQL](https://github.com/CovenantSQL).

## Quick Start

CQL client supports:

- macOS X 10.9+
- Linux 2.6.23+

## How CQL works

### 3 Layers Arch

- Layer 1: **Global Consensus Layer** (the main chain, the middle ring in the architecture diagram):
    - There will only be one main chain throughout the network.
    - Mainly responsible for database Miner and the user’s contract matching, transaction settlement, anti-cheating, shard chain lock hash and other global consensus matters.
- Layer 2: **SQL Consensus Layer** (shard chain, rings on both sides):
    - Each database will have its own separate shard chain.
    - Mainly responsible for: the signature, delivery and consistency of the various Transactions of the database. The data history of the permanent traceability is mainly implemented here, and the hash lock is performed in the main chain.
- Layer 3: **Datastore Layer** (database engine with SQL-92 support):
    - Each Database has its own independent distributed engine.
    - Mainly responsible for: database storage & encryption, query processing & signature, efficient indexing.

### Consensus Algorithm

CQL supports 2 kinds of consensus algorithm:

1. DPoS (Delegated Proof-of-Stake) is applied in `Eventually consistency mode` database and also `Layer 1 (Global Consensus Layer)` in BlockProducer. CQL miners pack all SQL queries and its signatures by the client into blocks thus form a blockchain.
2. BFT-Raft (Byzantine Fault-Toleranted Raft)is applied in `Strong consistency mode` database. The CQL miner leader does a `Two-Phase Commit` with `Kayak` to support `Transaction`.

CQL database consistency mode and node count can be selected in database creation with command  `cql create '{"UseEventualConsistency": true, "Node": 3}'`

## Comparison

|                              | Ethereum          | Hyperledger Fabric     | Amazon QLDB | CQL                                                          |
| ---------------------------- | ----------------- | ---------------------- | ----------- | ------------------------------------------------------------ |
| **Dev language**             | Solidity, ewasm   | Chaincode (Go, NodeJS) | ?           | Python, Golang, Java, PHP, NodeJS, MatLab                    |
| **Dev Pattern**              | Smart   Contract  | Chaincode              | SQL         | SQL                                                          |
| **Open Source**              | Y                 | Y                      | N           | Y                                                            |
| **Nodes for HA**             | 3                 | 15                     | ?           | 3                                                            |
| **Column Level ACL**         | N                 | Y                      | ?           | Y                                                            |
| **Data Format**              | File              | Key-value              | Document    | Key-value, Structured                                        |
| **Storage Encryption**       | N                 | API                    | Y           | Y                                                            |
| **Data Desensitization**     | N                 | N                      | N           | Y                                                            |
| **Multi-tenant**             | DIY               | DIY                    | N           | Y                                                            |
| **Consistency Delay**        | 2~6 min           | < 1 s                  | ?           | < 10 ms                                                      |
| **Secure for Open Internet** | Y                 | N                      | Only in AWS | Y                                                            |
| **Consensus**                | PoW + PoS(Casper) | CFT                    | ?           | DPoS (Eventually Consistency)<br/>BFT-Raft (Strong Consistency) |

## Use cases

### Traditional App

#### Privacy data

If you are a developper of password management tools just like [1Password](https://1password.com/) or [LastPass](https://www.lastpass.com/). You can use CQL as the database to take benefits:

1. Serverless: no need to deploy a server to store your user's password for sync which is the hot potato.
2. Security: CQL handles all the encryption work. Decentralized data storage gives more confidence to your users.
3. Regulation: CQL naturally comply with [GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation).

#### IoT storage

CQL miners are deployed globally, IoT node can write to nearest CQL miner directly.

1. Cheaper: Without passing all the traffic through a gateway, you can save a large bandwidth fee. And, CQL is a shared economic database which makes storage cheaper.
2. Faster: CQL consensus protocol is designed for Internet where network latency is unavoidable.

#### Open data service

For example, you are the most detailed Bitcoin OHLC data maintainer. You can directly expose an online SQL interface to your customers to meet a wide range of query needs.

1. CQL can limit specific SQL query statements to meet the needs while also balancing data security;
2. CQL can record SQL query records on the blockchain, which is very convenient for customers to check their bills for long-tail customers and billing.
3. For customers with high performance requirements, Slave nodes can be deployed at the customer to meet the needs of customers with low latency queries while enabling almost real-time data updates.

#### Secure storage

Thanks to the CQL data history is immutable, CQL can be used as a storage for sensitive operational logs to prevent hacking and erasure access logs.

### ĐApp

Storing data on Bitcoin or Ethereum is quite expensiveProgramming is very complicated due to the lack of support for structured data storage. CQL gives you a low-cost structured SQL database and also provides more room for ĐApp to exchange data with real-world.
