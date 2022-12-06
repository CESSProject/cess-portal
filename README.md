# **CESS-Portal**

cess-portal is the client of the cess project. By using some simple commands of cess-portal, you can easily realize a series of operations such as purchasing space, querying space, uploading/downloading files, and querying file information on the Linux system.

# Prerequisites

* one: linux system
* two: Go 1.19 and above

# **Build Code**

If you don't have git software on your machine, please install it first

```sh
yum install git -y
```

First you need to download the cess-portal project from GitHub

```sh
git clone https://github.com/CESSProject/cess-portal.git
```

Then run the build.sh(On Linux) or build.bat(On Windows) script file in the ‘install_ctl‘ folder,You can compile this project on any system,Before downloading, please install golang on the system and the version must be over 1.19.

```sh
#Compile with simple commands
cd /cess-portal/
go build portal.go
```

# **Config**

```sh
cd /cess-portal/

vim conf.toml
```

Let me introduce the content of the configuration file of cess-portal.

```toml
#The rpc address of the chain node
RpcAddr           = "wss://testnet-rpc0.cess.cloud/ws/"
#Phrase or seed for wallet account
AccountSeed       = "virtual field alert rapid wasp snap logic exact useless together stay settle"
#wallet account of cess 
AccountId = "cXjTYBWUY63uFG2t3ahAhmLtChz3WdBfXrDn4XaQY45pKLZBK"
```

Please edit the configuration of the above file, press the ESC key on the keyboard and enter': wq', then press the Enter key on keyboard for save it.
# **Getting Started**

## **Command group**

| command group name | subcommand name | features                                                     |
| ------------------ | --------------- | ------------------------------------------------------------ |
| query              | fstate           | Query state of the specified file in cess system  |
| query              | space           | Query space info of your account |
| query              | files            | Query file list in the specified bucket |
| query              | buckets          | Query bucket list of your account |
| file               | upload          | upload file |
| file               | download        | download file |
| file               | delete          | delete file |
| bucket             | create          | create new bucket for your account |
| bucket             | delete          | delete the specified bucket from your account |
| space              | purchase        | purchase storage space |
| space              | auth            | authorize purchased space for your account |
| space              | cancel          | cancel space authorization |


## **Global command**

-h,--help:Get the specific operation method of the command line

-c,--config:Absolute path, the address of the configuration file;

## **Operate example**

### 1.Query storage space info
```sh
./protal query space 
```
### 2.Query state of the specified file by file id
```sh
./protal query fstate 1e0ffe8a980aed71fc4f69f830076af19a5865194f0befb5b61475f1a18b9936
```
### 3.Query file list in the specified bucket
```sh
./protal query files "bucket-name"
```
### 4.Query bucket list
```sh
./protal query buckets
```
### 5.Upload file
```sh
./protal file upload "/opt/test_file" # either absolute or relative path
```
### 6.Download file by file id
```sh
./protal file download 1e0ffe8a980aed71fc4f69f830076af19a5865194f0befb5b61475f1a18b9936 ./data/cache # specify save path
```
### 7.Delete file by file id
```sh
./protal file delete 1e0ffe8a980aed71fc4f69f830076af19a5865194f0befb5b61475f1a18b9936
```
### 8.Create bucket
```sh
./protal bucket create "bucket-name"
```
### 9.Delete bucket
```sh
./protal bucket delete "bucket-name"
# Files in bucket will be deleted together
```
### 10.Purchase storage space
```sh
./protal space purchase 1 # purchase 1 GiB space,unit(GiB)
# Currently, the space only needs to be purchased once
```
### 11.Authorize purchased space
```sh
./protal space auth
# Authorize all the space purchased by the account, otherwise it will be unavailable
```
### 12.Cancel space authorization
```sh
./protal space cancel
# Make user space unavailable
```