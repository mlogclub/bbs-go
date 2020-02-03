<template>
  <div class="bbs-admin">
    <el-menu
      @select="handleSelect"
      class="el-menu-demo"
      mode="horizontal"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
    >
      <el-menu-item index="/">
        <span slot="title">
          <el-tooltip effect="dark" content="返回前台" placement="bottom-end">
            <i class="iconfont icon-home" style="color: #fff;"></i>
          </el-tooltip>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin">
        <span slot="title">
          <i class="iconfont icon-dashboard"></i>
          <span>后台</span>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin/nodes">
        <span slot="title">
          <i class="iconfont icon-tags"></i>
          <span>节点</span>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin/topics">
        <span slot="title">
          <i class="iconfont icon-topic"></i>
          <span>话题</span>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin/articles">
        <span slot="title">
          <i class="iconfont icon-article"></i>
          <span>文章</span>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin/comments">
        <span slot="title">
          <i class="iconfont icon-comment"></i>
          <span>评论</span>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin/users">
        <span slot="title">
          <i class="iconfont icon-username"></i>
          <span>用户</span>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin/links">
        <span slot="title">
          <i class="iconfont icon-link"></i>
          <span>链接</span>
        </span>
      </el-menu-item>
      <el-menu-item index="/admin/settings">
        <span slot="title">
          <i class="iconfont icon-setting"></i>
          <span>设置</span>
        </span>
      </el-menu-item>
    </el-menu>

    <nuxt />
  </div>
</template>

<script>
import utils from '~/common/utils'
import '~/plugins/element-ui'

export default {
  middleware: ['global', 'authenticated'],
  methods: {
    handleSelect(index, keyPath) {
      utils.linkTo(index)
    }
  },
  head() {
    return {
      title: 'BBS-GO管理',
      link: [
        {
          rel: 'stylesheet',
          href: '//unpkg.com/element-ui/lib/theme-chalk/index.css'
        }
      ]
    }
  }
}
</script>

<style lang="scss">
.bbs-admin {
  .page-container {
    margin: 10px !important;
    padding: 10px !important;
    background: #fff;

    .toolbar {
      border-bottom: 1px solid #e9ecf0;
      margin: 0 0 5px 0;
      padding: 10px 10px 0 10px;

      .el-form-item {
        margin-bottom: 10px;
      }
    }

    .pagebar {
      text-align: right;
      margin-top: 5px;
    }
  }
}
</style>
