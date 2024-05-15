module.exports = {
  apps: [
    {
      name: "bbs-go-site",
      port: "3000",
      // exec_mode: "cluster",
      // exec_mode: "fork",
      // instances: "max",
      script: "./.output/server/index.mjs",
    },
  ],
};
