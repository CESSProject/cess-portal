
#!/bin/bash
#-----------------------------------------------------------------------------#
# Modify the following configuration items according to your actual situation #
#↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓#
##Log information output directory
boardPath=''
##Cess chain communication address
cessRpcAddr=''
##Faucet address
faucetAddress=''
##Wallet private key
idAccountPhraseOrSeed=''
##Wallet public key
walletAddress=''
##The storage location of the file encryption password entered when uploading the file,using an absolute address
keyPath=''
##The path to download the file,using an absolute address
installPath=''



sudo sed -i "s|boardPath : \"\"|boardPath : \"${boardPath}\"|g" ./cess_client.yaml
sudo sed -i "s|cessRpcAddr : \"\"|cessRpcAddr : \"${cessRpcAddr}\"|g" ./cess_client.yaml
sudo sed -i "s|faucetAddress : \"\"|faucetAddress : \"${faucetAddress}\"|g" ./cess_client.yaml
sudo sed -i "s|idAccountPhraseOrSeed : \"\"|idAccountPhraseOrSeed : \"${idAccountPhraseOrSeed}\"|g" ./cess_client.yaml
sudo sed -i "s|walletAddress : \"\"|walletAddress : \"${walletAddress}\"|g" ./cess_client.yaml
sudo sed -i "s|keyPath : \"\"|keyPath : \"${keyPath}\"|g" ./cess_client.yaml
sudo sed -i "s|installPath : \"\"|installPath : \"${installPath}\"|g" ./cess_client.yaml

rm -rf /usr/bin/cessctl
rm -rf /etc/cess.d/
mkdir /etc/cess.d/
cp ./cess_client.yaml /etc/cess.d/
chmod 777 ./cessctl
mv ./cessctl /usr/bin/
cessctl -h
