#!/bin/sh
repeat_until_success () {
    echo Running command - "$@"
    i=0
    until $@
    do
        echo Command failed with exit code - $?
        if [ $i -gt 10 ]; then
            echo Giving up running command - "$@"
            return
        fi
        sleep 1
        echo Retrying
        i=`expr $i + 1`
    done
    echo Command finished successfully
}

#import private keys and then prefund them
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd importprivkey "cUzdgbfXTQv4SzYDCwRbZspoHZnsShXf5jt58Gqa4uy1gzWFxoDd" address1 # addr=mZ1SSGGtAav5b5rCgf4x5SphNLLv6EVtMT hdkeypath=m/88'/0'/1'
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd importprivkey "cW7jfi3QYQQAXzf1JPwa2ZtuyuVoyBNE8AD4UjVYcCFnpWYeXoyf" address2 # addr=meRTSNhCRNDmdNkvMumNVvEiZAresrzVbV hdkeypath=m/88'/0'/2'
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd importprivkey "cNX6Ccfn3VVAoRThjxm7RkZwDaJKLTABwoWk8EpN3CD9g6g1GfJN" address3 # addr=mPSNBRriZPyKRXPwXDPovGtx9pgFm4Erjs
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd importprivkey "cNLYNcaNQKYkDSpWaFT8b1uEjdcFjwQCyqDrnfyr2A9yh2ZBZQAQ" address4 # addr=mVtLccbttd4vCZCrf7gKCyZoeZWDjQvrS7
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd importprivkey "cN52NgLAheqLfw5mj8FEx6qtZoApbcjxMZXYg7UHMocPTHy6nV9s" address5 # addr=mNNs437u1qc6DVUbeGGpY7HsJPDRB8JuB4
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd importprivkey "cVrEQTSkp9o3vNCiLmjpG5QhqERRA3HkWLDPCS6WuHsGYx9E2jYj" address6 # addr=mNNiiJsBnPUZu4NwNtPaK9zyxD5ghV1By8
echo Finished importing accounts
echo Seeding accounts
# address1
echo Seeding mZ1SSGGtAav5b5rCgf4x5SphNLLv6EVtMT
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd generatetoaddress 1000 mZ1SSGGtAav5b5rCgf4x5SphNLLv6EVtMT
# address2
echo Seeding meRTSNhCRNDmdNkvMumNVvEiZAresrzVbV
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd generatetoaddress 1000 meRTSNhCRNDmdNkvMumNVvEiZAresrzVbV
# address3
echo Seeding mPSNBRriZPyKRXPwXDPovGtx9pgFm4Erjs
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd generatetoaddress 500 mPSNBRriZPyKRXPwXDPovGtx9pgFm4Erjs
# address4
echo Seeding mVtLccbttd4vCZCrf7gKCyZoeZWDjQvrS7
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd generatetoaddress 250 mVtLccbttd4vCZCrf7gKCyZoeZWDjQvrS7
# address5
echo Seeding mNNs437u1qc6DVUbeGGpY7HsJPDRB8JuB4
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd generatetoaddress 100 mNNs437u1qc6DVUbeGGpY7HsJPDRB8JuB4
# address6
echo Seeding mNNiiJsBnPUZu4NwNtPaK9zyxD5ghV1By8
repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd generatetoaddress 100 mNNiiJsBnPUZu4NwNtPaK9zyxD5ghV1By8
# playground pet shop dapp
# echo Seeding 0xCca81b02942D8079A871e02BA03A3A4a8D7740d2
# repeat_until_success metrix-cli -rpcuser=metrix -rpcpassword=testpasswd generatetoaddress 2 qcDWPLgdY9pTv3cKLkaMPvqjukURH3Qudy
echo Finished importing and seeding accounts
