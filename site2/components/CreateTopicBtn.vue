<template>
  <div>
    <el-dropdown
      placement="bottom"
      trigger="click"
      @command="handlePostCommand"
    >
      <span class="el-dropdown-link">
        <slot>
          <el-button type="primary" icon="el-icon-plus">发表</el-button>
        </slot>
      </span>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item command="tweet" icon="iconfont icon-tweet2"
          >发动态</el-dropdown-item
        >
        <el-dropdown-item command="topic" icon="iconfont icon-topic"
          >发帖子</el-dropdown-item
        >
        <el-dropdown-item command="article" icon="iconfont icon-article"
          >发文章</el-dropdown-item
        >
      </el-dropdown-menu>
    </el-dropdown>
  </div>
</template>
<script>
export default {
  data() {
    return {}
  },
  computed: {},
  methods: {
    handlePostCommand(cmd) {
      if (cmd === 'topic') {
        this.$linkTo('/topic/create')
      } else if (cmd === 'tweet') {
        this.$linkTo('/topic/create?type=1')
      } else if (cmd === 'article') {
        this.$linkTo('/article/create')
      }
    },
  },
}
</script>
<style lang="scss" scoped></style>
