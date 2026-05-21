module.exports = {
  apps: [
    {
      name: "bbs-go",
      script: "./bbs-go",
      env: {
        BBSGO_ENV: "prod",
      },
      // 如需使用集群模式，可根据机器配置开启：
      // exec_mode: "cluster",
      // instances: "max",
      // max_memory_restart: "512M",
    },
  ],
};
