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
          <!--
          <button class="button is-primary">
            <span class="icon">
              <i class="iconfont icon-add"></i>
            </span>
            <span>发表</span>
          </button>
          -->
        </slot>
      </span>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item
          v-for="(item, i) in modules"
          :key="i"
          :command="item.command"
          :icon="item.icon"
          >{{ item.name }}</el-dropdown-item
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
  computed: {
    config() {
      return this.$store.state.config.config
    },
    modules() {
      const modules = []
      for (let i = 0; i < this.config.modules.length; i++) {
        const item = this.config.modules[i]
        if (item.enabled) {
          const command = item.module
          let icon = ''
          let name = ''
          if (item.module === 'tweet') {
            icon = 'iconfont icon-tweet2'
            name = '发动态'
          } else if (item.module === 'topic') {
            icon = 'iconfont icon-topic'
            name = '发帖子'
          } else if (item.module === 'article') {
            icon = 'iconfont icon-article'
            name = '发文章'
          }
          modules.push({
            command,
            icon,
            name,
          })
        }
      }
      return modules
    },
  },
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
