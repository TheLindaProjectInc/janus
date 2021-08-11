# Simple VUE project to switch to METRIX network via Metamask

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

### wallet_addEthereumChain
```
// request account access
window.ethereum.request({ method: 'eth_requestAccounts' })
    .then(() => {
        // add chain
        window.ethereum.request({
            method: "wallet_addEthereumChain",
            params: [{
                {
                    chainId: '0x71',
                    chainName: 'Metrix Testnet',
                    rpcUrls: ['https://localhost:23889'],
                    blockExplorerUrls: ['https://testnet-explorer.metrixcoin.com/'],
                    iconUrls: [
                        'https://metrixcoin.com/images/metamask_icon.svg',
                        'https://metrixcoin.com/images/metamask_icon.png',
                    ],
                }
            }],
        }
    });
```

# Known issues
- Metamask requires https for `rpcUrls` so that must be enabled
  - Either directly through Janus with `--https-key ./path --https-cert ./path2` see [SSL](../README.md#ssl)
  - Through the Makefile `make docker-configure-https && make run-janus-https`
  - Or do it yourself with a proxy (eg, nginx)
- Client side transaction signing is not implemented yet
