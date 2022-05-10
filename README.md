# **CESS-Portal**

CESS-Portal is the storage client of the CESS project. By using some simple commands of CESS-Portal, you can easily realize a series of operations such as purchasing space, querying space, uploading/downloading files, and querying file information on the Linux system.

CESS-Portal characteristics of operation function is shown below:

- Supports enterprise-level commercial storage service
- Efficient data retrieval and return ability
- low data redundancy during storage
- provides space storage service
- high data security

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


## **Global command**

```shell
-h,--help：Get the specific operation method of the command line
```

```shell
-c,--config：Absolute path, the address of the configuration file
```


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
| file               | delete          | Delete the specific file                                      |



## Top up tokens from the faucet

* instruction：

  ​		Get a certain amount of tokens through the faucet service

* usage：

  ```shell
  cessctl trade obtain <public key>
  ```
  
* example：

  ```shell
  cessctl trade obtain 0x2ed4a2c67291bf3eaa4de538ab120ba21b3de1b5704551864226d2fae8f87937 
  ```




## **Operate example on Linux**

### (A)Query the current storage price per MB

* instruction：

  ​		Query the current storage price per MB(Unit: TCESS / MB)

* usage：


  ```shell
  cessctl find price
  ```
      
* example：

  ```shell
  # cessctl find price
  ```
  

  2022/04/25 14:06:23 Connecting to ws://


  [Success]The current storage price is:8.342023 per (MB)



### (B)Buy CESS storage space 

* instruction：

  ​		buy the storage space
  
* usage：

  ```shell
  cessctl trade exp <quantity><duration>
  ```
  
* example：

  ```shell
  cessctl trade exp 1 20
  
  2022/04/21 16:05:14 Connecting to ws://
  
  [Success]Buy space on chain success!
  ```

  Tips: 
  
  [quantity] - The amount of space you want to buy, the unit is GB   
  [duration] - The Time of storage space you want to use, the unit is month



### (C)Query real-time storage space information 

* instruction：

  ​		Query storage space information include storage space, purchased storage space, and remaining storage space

* usage：

  ```shell
  cessctl find space
  ```

* example：
  
  ```shell
  # cessctl find space
  
  
  2022/04/21 16:09:55 Connecting to ws://

  ——You Purchased Space——
  
  PurchasedSpace:3145728(KB)
  
  UsedSpace:24582(KB)
  
  RemainingSpace:3121146(KB)
  ```




### (D)Upload the specific files

* instruction：

  ​		Send local source files to scheduling nodes

* usage：

  ```shell
  cessctl upload file <filepath><backups><private key>
  ```
  
* example：

  ```shell
  # cessctl file upload /root/test.txt 3 1234567887654321
  
  [Warming] Do you want to upload your file without private key (it's means your file status is public)?
  
  You can type the 'private key' or enter with nothing to jump it:
  
  2022/04/21 16:42:38 Connecting to ws://

  File meta info upload:success! ,fileid is:1517061233797238784
  
  [██████████████████████████████████████████████████]100%         1/1
  
  [Success]:upload file:/root/test.txt successful!#
  ```

### (E)Download the specific file

* instruction：

  ​		Download specific file based on fileid

* usage：

  ```shell
  cessctl file download <fileid>
  ```
  
* example：

  ```shell
  # cessctl file download 1517061233797238784
  
  2022/04/21 16:58:39 Connecting to ws://
  
  [██████████████████████████████████████████████████]100%         1/1
  
  [OK]:File 'test.txt' has been downloaded to the directory :/root/installpath/test.txt
  ```


### (F)Delete the specific file

* instruction：

  ​		Delete the specific file meta information

* usage：

  ```shell
  cessctl file delete <fileid>
  ```
  
* example：

  ```shell
  # cessctl file delete 1517061233797238784
  
  2022/04/21 17:02:57 Connecting to ws://
  
  [OK]Delete fileid:1517061233797238784 success!
  ```

