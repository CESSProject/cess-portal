# **cess-portal**

cess-portal is the client of the cess project. By using some simple commands of cess-portal, you can easily realize a series of operations such as purchasing space, querying space, uploading/downloading files, and querying file information on the Linux system.

# **Build Code**

First you need to download the cess-portal project from GitHub

```shell
git clone https://github.com/CESSProject/cess-portal.git
```
Then run the build.bat script file in the ‘install_ctl‘ folder，You can compile this project on any system

```shell
##Compile with script
cd tar
./build.bat
```

Finally, please place the 'install_ctl' folder in your Linux environment

# **Install On Linux**

```shell
cd install_ctl
##Provide run permission
chmod 777 install-cessctl.sh
##Configure a one-click install script
vim install-cessctl.sh
```

Let me introduce the content of the configuration file of the one-click installation script.

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

* :one: Centos 8 and above
* :two: Dedicate IP

## **Command group**

| command group name | subcommand name | features                                                     |
| ------------------ | --------------- | ------------------------------------------------------------ |
| find               | price           | Query the current storage price per MB                       |
| find               | space           | Query currently used space, purchased space, remaining space |
| find               | file            | Query file or file list                                      |
| file               | upload          | upload files                                                 |
| file               | download        | download file                                                |
| trade              | exp             | buy storage                                                  |
| trade              | obtain          | Get coins from the faucet                                    |
| file               | delete          | delete file                                                  |



## **Global command**

-h,--help：Get the specific operation method of the command line

-c,--config：Absolute path, the address of the configuration file; used when not defined:/etc/cess.d/cess_client.yaml



## **Configuration file**

boardPath：Absolute path, the data kanban location of the result output; if not defined, output to: /etc/cess.d file.

cessRpcAddr：Chain interaction address, the address that interacts with the chain.

faucetAddress：Faucet address, the address to get coins from the faucet.

idAccountPhraseOrSeed：Account private key, which is used as the user's mnemonic when signing transactions.

accountPublicKey：The publicKey of public key conversion, used to query data on the chain, and the conversion address: https://polkadot.subscan.io/tools/ss58_transform.

walletAddress：The wallet public key address, the owner id of the file when uploading the file metadata.



## **Operate example**

### (A)Query storage unit price

* instruction：

  ​		Chain query and displays the current lease storage space Price (Unit: Cess / MB)

* usage：

  ​		cessctl find price

* example：

  ​		cessctl find price -c /root/cess_client.yaml



### (B)Check remaining space

* instruction：

  ​		Chain query current account purchased storage space usage (used and remaining)

* usage：

  ​		cessctl find space

* example：

  ​		cessctl find space -c /root/cess_client.yaml



### (C)Query file information

* instruction：

  ​		Chain query all file information that has been uploaded by the current account (sorting, keyword retrieval...)

* usage：

  ​		cessctl find file <fileid>

  ​		If fileid is vaild, output all file list information

* example：

  ​		Query single file information:cessctl find file 1483720947931287552 -c /root/cess_client.yaml

  ​		Query file list information:cessctl find file -c /root/cess_client.yaml



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
