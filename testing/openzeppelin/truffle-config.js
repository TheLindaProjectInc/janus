module.exports = {
  test_directory: "./openzeppelin-contracts/test",
  migrations: "./openzeppelin-contracts/migrations",
  contracts_directory: "./openzeppelin-contracts/contracts",
  contracts_build_directory: "./openzeppelin-contracts/build/output",
  mocha: {
    reporter: "mocha-spec-json-output-reporter",
    reporterOptions: {
      fileName: "output.json",
    },
  },
  networks: {
    development: {
      host: "127.0.0.1",
      port: 23889, //Switch to 23888 for local HTTP Server, look at Makefile run-janus
      network_id: "*",
      gasPrice: "0x1388",
    },
    testing: {
      host: "127.0.0.1",
      port: 23888,
      network_id: "*",
      gasPrice: "0x1388",
    },
    docker: {
      host: "janus",
      port: 23889,
      network_id: "*",
      gasPrice: "0x1388",
    },
    ganache: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*",
    },
    testnet: {
      host: "206.189.117.117",
      port: 23889,
      network_id: "*",
      from: "0xb89a4201258da334e3cd6d49047715fbf8a0e386",
      gasPrice: "0x1388",
    },
  },
  compilers: {
    solc: {
      version: "^0.8.0",
      settings: {
        optimizer: {
          enabled: true,
          runs: 1,
        },
      },
    },
  },
};
