<template>
  <div class="hello">
    <div v-if="web3Detected">
      <b-button v-if="metrixConnected">Connected to METRIX</b-button>
      <b-button v-else-if="connected" v-on:click="connectToMetrix()">Connect to METRIX</b-button>
      <b-button v-else v-on:click="connectToWeb3()">Connect</b-button>
    </div>
    <b-button v-else>No Web3 detected - Install metamask</b-button>
  </div>
</template>

<script>
let METRIXMainnet = {
  chainId: '0x71',
  chainName: 'Metrix Mainnet',
  rpcUrls: ['https://localhost:23889'],
  blockExplorerUrls: ['https://metrixcoin.com/'],
  iconUrls: [
    'https://metrixcoin.com/images/metamask_icon.svg',
    'https://metrixcoin.com/images/metamask_icon.png',
  ],
};
let METRIXTestNet = {
  chainId: '0x71',
  chainName: 'Metrix Testnet',
  rpcUrls: ['https://localhost:23889'],
  blockExplorerUrls: ['https://testnet-explorer.metrixcoin.com/'],
  iconUrls: [
    'https://metrixcoin.com/images/metamask_icon.svg',
    'https://metrixcoin.com/images/metamask_icon.png',
  ],
};
let config = {
  "0x1": METRIXMainnet,
  // ETH Ropsten
  "0x3": METRIXTestNet,
  // ETH Rinkby
  "0x4": METRIXTestNet,
  // ETH Görli
  "0x5": METRIXTestNet,
  // ETH Kovan
  "0x71": METRIXTestNet,
};
config[METRIXMainnet.chainId] = METRIXMainnet;
config[METRIXTestNet.chainId] = METRIXTestNet;

export default {
  name: 'Web3Button',
  props: {
    msg: String,
    connected: Boolean,
    metrixConnected: Boolean,
  },
  computed: {
    web3Detected: function() {
      return !!this.Web3;
    },
  },
  methods: {
    getChainId: function() {
      return window.ethereum.chainId;
    },
    isOnMetrixChainId: function() {
      let chainId = this.getChainId();
      return chainId == METRIXMainnet.chainId || chainId == METRIXTestNet.chainId;
    },
    connectToWeb3: function(){
      if (this.connected) {
        return;
      }
      let self = this;
      window.ethereum.request({ method: 'eth_requestAccounts' })
        .then(() => {
          console.log("Emitting web3Connected event");
          let metrixConnected = self.isOnMetrixChainId();
          let currentlyMetrixConnected = self.metrixConnected;
          self.$emit("web3Connected", true);
          if (currentlyMetrixConnected != metrixConnected) {
            console.log("ChainID matches METRIX, not prompting to add network to web3, already connected.");
            self.$emit("metrixConnected", true);
          }
        })
        .catch(() => {
          console.log("Connecting to web3 failed", arguments);
        })
    },
    connectToMetrix: function() {
      console.log("Connecting to Metrix, current chainID is", this.getChainId());

      let self = this;
      let metrixConfig = config[this.getChainId()] || METRIXTestNet;
      console.log("Adding network to Metamask", metrixConfig);
      window.ethereum.request({
        method: "wallet_addEthereumChain",
        params: [metrixConfig],
      })
        .then(() => {
          self.$emit("metrixConnected", true);
        })
        .catch(() => {
          console.log("Adding network failed", arguments);
        })
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
