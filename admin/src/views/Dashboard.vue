<template>
  <div class="index-container">
    <h1 class="index-title">后台首页</h1>
    <div v-if="systeminfo" class="widget">
      <div class="widget-header">服务器信息</div>
      <div class="widget-content">
        <div>OS: {{ systeminfo.os }}</div>
        <div>Arch: {{ systeminfo.arch }}</div>
        <div>Cpu: {{ systeminfo.numCpu }}</div>
        <div>Go version: {{ systeminfo.goversion }}</div>
        <div>Hostname: {{ systeminfo.hostname }}</div>
      </div>
    </div>
    <div class="widget">
      <div class="widget-header">开发团队信息</div>
      <div class="widget-content">
        <div>
          Official Site:
          <a href="https://mlog.club/" target="_blank">https://mlog.club/</a>
        </div>
        <div>
          Dev Team:
          <a href="https://github.com/mlogclub/" target="_blank">码农俱乐部</a>
        </div>
        <div>
          Github:
          <a href="https://github.com/mlogclub/bbs-go" target="_blank"
            >https://github.com/mlogclub/bbs-go</a
          >
        </div>
        <div>
          Gitee:
          <a href="https://gitee.com/mlogclub/bbs-go" target="_blank"
            >https://gitee.com/mlogclub/bbs-go</a
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      systeminfo: null,
    };
  },
  mounted() {
    this.load();
  },
  methods: {
    async load() {
      this.systeminfo = await this.axios.get("/api/admin/common/systeminfo");
    },
  },
};
</script>

<style lang="scss" scoped>
.index-container {
  margin: 20px;
  font-size: 14px;
  font-family: inherit;
  color: inherit;
  font-weight: 500;
  line-height: 1.8;

  .widget-content {
    margin-top: 0px;
  }

  .index-title {
    font-size: 24px;
    font-weight: 700;
  }

  .card {
    padding: 10px;
    position: relative;
    display: flex;
    flex-direction: column;
    min-width: 0;
    word-wrap: break-word;
    background-clip: border-box;
    background-color: #fff;
    border: 1px solid rgba(0, 0, 0, 0.125);
    border-radius: 0.25rem;
  }
}
</style>
