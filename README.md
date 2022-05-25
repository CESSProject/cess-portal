# **CESS-Portal**

cess-portal is the client of the cess project. By using some simple commands of cess-portal, you can easily realize a series of operations such as purchasing space, querying space, uploading/downloading files, and querying file information on the Linux system.

# **Build Code**

If you don't have git software on your machine, please install it first

```shell
yum install git -y
```

First you need to download the cess-portal project from GitHub

```shell
git clone https://github.com/CESSProject/cess-portal.git
```

Then run the build.sh(On Linux) or build.bat(On Windows) script file in the ‘install_ctl‘ folder,You can compile this project on any system,Before downloading, please install golang on the system and the version must be over 1.17.

```shell
##Compile with script
cd /cess-portal/install_ctl

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

```shell
##The log file of the client's operation output, your operation results will be recorded in the output.log file under this file
boardPath='/root'
##The mailing address of the CESS chain
cessRpcAddr='ws://xxx.xx.xx.xxx:9949/'
##tCESS pick-up tap address
faucetAddress='http://xx.xxx.xx.xx:9708/transfer'
##Memo Seed for Wallet
idAccountPhraseOrSeed='lazy funny invest opinion jaguar romance anger return glare flat lift clap'
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

* :one: Centos
* :two: Go 1.17 and above

## **Command group**

| command group name | subcommand name | features                                                     |
| ------------------ | --------------- | ------------------------------------------------------------ |
| query              | price           | Query the current storage price per MB                       |
| query              | space           | Query currently used space, purchased space, remaining space |
| query              | file            | Query file or file list                                      |
| file               | upload          | upload files                                                 |
| file               | download        | download file                                                |
| file               | delete          | delete file                                                  |
| file               | decrypt         | decrypt encrypted files                                      |
| purchase           | storage         | buy storage                                                  |
| purchase           | free            | Get coins from the faucet                                    |



## **Global command**

-h,--help:Get the specific operation method of the command line

-c,--config:Absolute path, the address of the configuration file; used when not defined:/etc/cess.d/cess_client.yaml



## **Configuration file**

boardPath:Absolute path, the data kanban location of the result output; if not defined, output to: /etc/cess.d file.

cessRpcAddr:Chain interaction address, the address that interacts with the chain.

faucetAddress:Faucet address, the address to get coins from the faucet.

idAccountPhraseOrSeed:Account private key, which is used as the user's mnemonic when signing transactions.

walletAddress:The wallet public key address, the owner id of the file when uploading the file metadata.



## **Operate example**

### (A)Query storage unit price

* instruction:

  ​		Chain query and displays the current lease storage space Price (Unit: TCess / GB)

* usage:

  ​		cessctl query price

* example:

  ​		cessctl query price



### (B)Check remaining space

* instruction:

  ​		Chain query current account purchased storage space usage (used and remaining)

* usage:

  ​		cessctl query space

* example:

  ​		cessctl query space



### (C)Query file information

* instruction:

  ​		Chain query all file information that has been uploaded by the current account (sorting, keyword retrieval...)

* usage:

  ​		cessctl query file <file id>

  ​		If fileid is vaild, output all uploaded file information of the user in the configuration file

* example:

  ​		Query single file information:cessctl query file 1483720947931287552

  ​		Query file list information:cessctl query file



### (D)Upload files

* instruction:

  ​		Send local source files to scheduling nodes

* usage:

  ​		cessctl file upload <file path> <backups>

  ​		file path:The absolute path of the file, not a folder

  ​		backups:The number of backups of uploaded files in the CESS system. The more the number of backups, the more secure it is and the more space it consumes

* example:

  ​		cessctl file upload /root/test.txt 3



### (E)Download file

* instruction:

  ​		Download file based on fileid,Download to the installPath path in the configuration file

* usage:

  ​		cessctl file download <file id>

  ​		fileid:The unique id of the file

* example:

  ​		cessctl file download 1483720947931287552



### (F)Buy space

* instruction:

  ​		Send on-chain transactions, buy space

* usage:

  ​		cessctl purchase storage <space quantity> <space duration> <expected price>

  ​		space quantity:The number of expansion capacity, unit: 1/1GB

  ​		space duration:You want to buy space for several time units, unit: 1/1month

  ​		expected price:The maximum acceptable price for buying space, in cess; if it is empty, all prices are accepted

* example:

  ​		expected price 20cess:cessctl purchase storage 1 1 20

  ​		All price accepted:cessctl purchase storage 1 1



### (Y)Tap to get tokens

* instruction:

  ​		Get a certain amount of tokens through the faucet service

* usage:

  ​		cessctl purchase free <wallet address>

  ​		address:wallet address

* example:

  ​		cessctl purchase free cXjsAyird2dizRjmHML9Eqxp1MGodGdEUHv8rjr7z56Dv5A7C



### (T)File delete

* instruction:

  ​		Delete file meta information.

* usage:

  ​		cessctl file delete <file id>

  ​		fileid:file unique id

* example:

  ​		cessctl file delete 1506154108548026368

###  (L)File decrypt

* instruction:

  ​		Decrypt the files that have not been decrypted, and the decrypted files will be stored in the 'installPath' path of the configuration file.

* usage:

  ​		cessctl file decrypt  <file path>

  ​		filepath:path to the file that needs to be decrypted

* example:

  ​		cessctl file decrypt /root/test.txt
