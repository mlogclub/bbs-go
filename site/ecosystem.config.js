module.exports = {
  apps: [
    {
      name: "bbs-go-site",
      port: "3000",
      exec_mode: "fork",
      // exec_mode: "cluster",
      // instances: "max",
      script: "./.output/server/index.mjs",
    },
  ],
};
