

# Background
Git, created by Linus Torvalds in 2005, has become one of the most widely-used version control systems in both open-source and commercial fields. It has become an essential tool for modern software development. However, while Git's remote repositories are stored in a centralized platform GitHub, users may face restrictions when accessing repositories due to external factors, such as political restrictions. As we know, in its history, GitHub has banned numerous developers' accounts and limited access to Github's open-source projects. This has sparked heated discussions about the idea of open source truly having no borders? Avoiding centralized control of code assets has become a hot topic.

As the Web3 field continues to grow, more and more people are recognizing the significance of ownership. To avoid being at the mercy of centralized platforms, many users are opting to store their data in decentralized storage networks. The Greenfield storage network, recently introduced by BNB Chain, has emerged as a powerful player in this space. It offers not only decentralized data storage but also unique permission management features. Furthermore, its cross-chain programming model with the BSC network opens up a world of possibilities for this ecosystem.

CodexField has emerged as an innovative solution that allows developers to save their code on the decentralized storage network Greenfield, the code saved on Greenfield will be distributed throughout the Storage Provider network, enabling it to resist censorship while being readily accessible, which provides developers with complete ownership of their codes.

CodexField is a decentralized platform aimed at developers, which provides a fully compatible experience with Git, allowing developers to use the toolset to develop and upload code to Greenfield. Furthurmore, codexfield proposes an innovative solution for trading codes named Code Marketplace, which is a platform where developers can sell their code saved on Greenfield at their own prices. To ensure quality, codexfield also introduces a rating mechanism, which enables users to rate the codes, creating a reputation-based trading platform for developers on the blockchain.


# Gitd 

Gitd means “Git for CodexField”, or in other words, “Git for Decentralized Storage”. The Gitd tool is fully compatible with Git's functionality and usage, enabling developers to use Gitd for version control and code submission.

The CodeSync plugin facilitates one-click migration of user-submitted code from Github to CodexField, saving it on the Greenfield network.

CodexField frontend is a web-hosted frontend page setup on Greenfield, which allows users to view and manage code stored in Greenfield through Gitd.

By default, the code uploaded through Gitd is set to private access, visible only to the author. However, users can choose to make their code public, which will be displayed on codexfield and visible to all users. Moreover, users can list their private code for sale on the Code Marketplace.

# Usage
## Set environment

```shell
export GREENFIELD_CHAIN_ID=greenfield_1017-1  // greenfield mainnet
# export GREENFIELD_CHAIN_ID=greenfield_5600-1  // greenfield testnet
# see https://docs.bnbchain.org/greenfield-docs/docs/api/endpoints/

# use metamask to generate a new account and get your private key
# use greenfield testnet faucet to get some testBNB
export GREENFIELD_PRIVATE_KEY=xxxx
```
## Build & Install

### Build from source
```
make gitd
```

### Download pre-build binary
```shell
curl -fsSL https://raw.githubusercontent.com/codexfield/gitd/develop/install.sh | bash
```

## Init A Repo Locally
```shell
mkdir <repo>
cd <repo>
gitd init
```

## Push to Greenfield Repo
```shell
gitd remote add origin gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/<repoName>
echo "Hello CodexField" >> README.md
gitd add README.md
gitd commit -m "add README.md"
gitd push origin main -f  // when push firstly, please use force push. will fix later.
```

## Clone 

```shell
cd <new_folder>
gitd clone gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/<repoName>
```

# Reference

- [go git](https://github.com/go-git/go-git): go-git is a highly extensible git implementation library written in pure Go.
- [s3 git](https://github.com/s3git/s3git): s3git is a simple CLI tool that allows you to create a distributed, decentralized and versioned repository.
- [Greenfield](https://github.com/bnb-chain/greenfield): the greenfield blockchain
- [Greenfield-go-sdk](https://github.com/bnb-chain/greenfield-go-sdk): the greenfield go sdk
- [Greenfield-cmd](https://github.com/bnb-chain/greenfield-cmd): the greenfield command line tool
- [Greenfield-Contract](https://github.com/bnb-chain/greenfield-contracts): the cross chain contract for Greenfield that deployed on BSC network.
- [Greenfield-Tendermint](https://github.com/bnb-chain/greenfield-tendermint): the consensus layer of Greenfield blockchain.
- [Greenfield-Storage-Provider](https://github.com/bnb-chain/greenfield-storage-provider): the storage service infrastructures provided by either organizations or individuals.
- [Greenfield-Relayer](https://github.com/bnb-chain/greenfield-relayer): the service that relay cross chain package to both chains.
- [Greenfield-Cmd](https://github.com/bnb-chain/greenfield-cmd): the most powerful command line to interact with Greenfield system.
- [Awesome Cosmos](https://github.com/cosmos/awesome-cosmos): Collection of Cosmos related resources which also fits Greenfield.
