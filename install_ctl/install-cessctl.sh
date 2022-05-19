
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

if [ -f "/usr/bin/cessctl" ]; then
  rm -rf /usr/bin/cessctl
fi
if [ -e "/etc/cess.d/" ]; then
  rm -rf /etc/cess.d/
fi
mkdir /etc/cess.d/
cp ./cess_client.yaml /etc/cess.d/
chmod 777 ./cessctl
mv ./cessctl /usr/bin/
cessctl -h

sudo sed -i "s|boardPath : \"\"|boardPath : \"${boardPath}\"|g" /etc/cess.d/cess_client.yaml
sudo sed -i "s|cessRpcAddr : \"\"|cessRpcAddr : \"${cessRpcAddr}\"|g" /etc/cess.d/cess_client.yaml
sudo sed -i "s|faucetAddress : \"\"|faucetAddress : \"${faucetAddress}\"|g" /etc/cess.d/cess_client.yaml
sudo sed -i "s|idAccountPhraseOrSeed : \"\"|idAccountPhraseOrSeed : \"${idAccountPhraseOrSeed}\"|g" /etc/cess.d/cess_client.yaml
sudo sed -i "s|walletAddress : \"\"|walletAddress : \"${walletAddress}\"|g" /etc/cess.d/cess_client.yaml
sudo sed -i "s|keyPath : \"\"|keyPath : \"${keyPath}\"|g" /etc/cess.d/cess_client.yaml
sudo sed -i "s|installPath : \"\"|installPath : \"${installPath}\"|g" /etc/cess.d/cess_client.yaml
