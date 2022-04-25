# **CESS-Portal**

CESS-Portal is the client of the CESS project. By using some simple commands of CESS-Portal, you can easily realize a series of operations such as purchasing space, querying space, uploading/downloading files, and querying file information on the Linux system.

CESS-Portal  characteristics of operation function is shown below:

- provides space storage service
- low data redundancy during storage
- high data security
- Supports enterprise-level commercial storage service
- Efficient data retrieval and return ability

# **Build Code**

If you don't have git software on your machine, please install it first

```shell
yum install git -y
```

First you need to download the cess-portal project from GitHub

```shell
git clone https://github.com/CESSProject/cess-portal.git
```
Then run the build.sh(On Linux) or build.bat(On Windows) script file in the ‘install_ctl‘ folder，You can compile this project on any system，Before downloading, please install golang on the system and the version must be over 1.17.

```shell
##Compile with script
cd /cess-portal/install_ctl

##Run it on windows platform
./build.bat
##Run it on linux platform
sh build.sh
```

Finally, you can place the 'install_ctl' folder in your Linux environment,you can also operate directly in this folder

# **Install On Linux**

```shell
#If you are not in the install_ctl folder, please enter first
cd install_ctl
##Provide run permission
chmod 777 install-cessctl.sh
##Configure a one-click install script
vim install-cessctl.sh
```

Let me introduce the content of the configuration file of the one-click installation script.


## **Configuration file**

boardPath：Absolute path, the data kanban location of the result output; if not defined, output to: /etc/cess.d file.

cessRpcAddr：Chain interaction address, the address that interacts with the chain.

faucetAddress：Faucet address, the address to get coins from the faucet.

idAccountPhraseOrSeed：Account private key, which is used as the user's mnemonic when signing transactions.

accountPublicKey：The publicKey of public key conversion, used to query data on the chain, and the conversion address: https://polkadot.subscan.io/tools/ss58_transform.

walletAddress：The wallet public key address, the owner id of the file when uploading the file metadata.


```shell
##The log file of the client's operation output, your operation results will be recorded in the output.log file under this file
boardPath='/root'
##The mailing address of the CESS chain
cessRpcAddr='ws://xxx.xx.xx.xxx:9949/'
##tCESS pick-up tap address
faucetAddress='http://xx.xxx.xx.xx:9708/transfer'
##Memo Seed for Wallet
idAccountPhraseOrSeed='lazy funny invest opinion jaguar romance anger return glare flat lift clap'
##The public key address of the wallet, which is generated from the wallet address, and the generated address: https://polkadot.subscan.io/tools/ss58_transform
accountPublicKey='0x1c298066dcd205a267df5b29a2ec7104b03b27e009dd3166f7318194eb9ee77a'
##wallet address
walletAddress='5AhdZVDwjFXpvbsTjHaXv2jqNos49zFFnb5K4A1hnzVSo1iR'
##If the file upload is encrypted, the password memo of the file will be saved here, and it can be created to the next directory of the existing folder.
keyPath='/root/keypath'
##The path address of the file download, the downloaded files will appear here, support to create the next level directory of the existing folder
installPath='/root/cessDownload'
```

Please edit the configuration of the above file, press the ESC key on the keyboard and enter': wq', then press the Enter key on keyboard for save it.Next you can run the script to install.

```shell
./install-cessctl.sh
```

# **Getting Started**

## Prerequisites

**Software requirement:**
* :one: Centos
* :two: Go 1.17 and above

**Hardware requirement:**
* :one: Memory 8GB and above
* :two: internet


## **Global command**

-h,--help：Get the specific operation method of the command line

-c,--config：Absolute path, the address of the configuration file; used when not defined:/etc/cess.d/cess_client.yaml



## **Command group**

| command group name | subcommand name | features                                                     |
| ------------------ | --------------- | ------------------------------------------------------------ |
| find               | price           | Query the current storage price per MB                       |
| find               | space           | Query real-time storage space information                    |
| find               | file            | Query the uploaded files information                         |
| file               | upload          | Upload the specific files                                    |
| file               | download        | Download the specific file                                   |
| trade              | exp             | Buy CESS storage space                                       |
| trade              | obtain          | Top up tokens from the faucet                                |
| file               | delete          | Delete the specifc file                                      |










## **Operate example**

### (A)Query storage unit price

* instruction：

  ​		Query the current storage price per MB(Unit: TCESS / MB)

* usage：


  ```shell
  ​		cessctl find price
  ```
      
* example：

  [root@iZbp18tsw8ozfwv5y1z6avZ ~]# cessctl find price
  
  

  2022/04/25 14:06:23 Connecting to *** . ** . ** . *** : ****...


  [Success]The current storage price is:8.342023 per (MB)



### (B)Buy CESS storage space 

* instruction：

  ​		buy the storage space
  
* usage：

  ​		cessctl trade exp <quantity><duration>

* example：

  cessctl trade exp 1 20
  
  2022/04/21 16:05:14 Connecting to ws://*** .** .** .***:****/...
  
  [Success]Buy space on chain success!



### (C)Query real-time storage space information 

* instruction：

  ​		Query storage space information include storage space, purchased storage space, and remaining storage space

* usage：

  ​		cessctl find space


* example：
  
[root@iZbp18tsw8ozfwv5y1z6avZ admin]# cessctl find space
  
  
2022/04/21 16:09:55 Connecting to ws://*** .** .** .*** :****/...

  ——You Purchased Space——
  
  PurchasedSpace:3145728(KB)
  
  UsedSpace:24582(KB)
  
  RemainingSpace:3121146(KB)
  




### (D)Upload files

* instruction：

  ​		Send local source files to scheduling nodes

* usage：

  ​		cessctl upload file <filepath> <downloadfee>

  ​		filepath：The absolute path of the file, not a folder

  ​		downloadfee：The cost of downloading the file, in cess

* example：

  ​		cessctl upload file /root/cess_client.yaml 10 -c /root/cess_client.yaml



### (E)Download file

* instruction：

  ​		Download file based on fileid

* usage：

  ​		cessctl download <fileid>

  ​		fileid：The unique id of the file

* example：

  ​		cessctl download file 1483720947931287552 -c /root/cess_client.yaml



### (F)Buy space

* instruction：

  ​		Send on-chain transactions, buy space

* usage：

  ​		download exp <spacequantity> <expected price>

  ​		spacequantity：The number of expansion capacity, unit: 1/512MB

  ​		expected price：The maximum acceptable price for buying space, in cess; if it is empty, all prices are accepted

* example：

  ​		expected price 20cess：download exp 1 20 -c /root/cess_client.yaml

  ​		All price accepted：download exp 1 -c /root/cess_client.yaml



### (Y)Tap to get tokens

* instruction：

  ​		Get a certain amount of tokens through the faucet service

* usage：

  ​		cessctl obtain <address>

  ​		address：publickey of the account

* example：

  ​		cessctl obtain 0x2ed4a2c67291bf3eaa4de538ab120ba21b3de1b5704551864226d2fae8f87937 -c /root/cess_client.yaml



### (T)File delete

* instruction：

  ​		Delete file meta information.

* usage:

  ​		cessctl file delete <fileid>

  ​		fileid：file unique id

* example：

  ​		cessctl file delete 1506154108548026368 -c /root/cess_client.yaml
