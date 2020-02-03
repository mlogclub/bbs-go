<template>
  <div class="bbs-admin">
    <el-menu
      class="el-menu-demo"
      mode="horizontal"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
    >
      <el-menu-item index="1">
        <span slot="title">
          <nuxt-link to="/admin">
            <i class="iconfont icon-dashboard"></i>
            <span>后台</span>
          </nuxt-link>
        </span>
      </el-menu-item>
      <el-menu-item index="2">
        <span slot="title">
          <nuxt-link to="/admin/nodes">
            <i class="iconfont icon-tags"></i>
            <span>节点</span>
          </nuxt-link>
        </span>
      </el-menu-item>
      <el-menu-item index="3">
        <span slot="title">
          <nuxt-link to="/admin/topics">
            <i class="iconfont icon-topic"></i>
            <span>话题</span>
          </nuxt-link>
        </span>
      </el-menu-item>
      <el-menu-item index="4">
        <span slot="title">
          <nuxt-link to="/admin/articles">
            <i class="iconfont icon-article"></i>
            <span>文章</span>
          </nuxt-link>
        </span>
      </el-menu-item>
      <el-menu-item index="5">
        <span slot="title">
          <nuxt-link to="/admin/comments">
            <i class="iconfont icon-comment"></i>
            <span>评论</span>
          </nuxt-link>
        </span>
      </el-menu-item>
      <el-menu-item index="6">
        <span slot="title">
          <nuxt-link to="/admin/users">
            <i class="iconfont icon-username"></i>
            <span>用户</span>
          </nuxt-link>
        </span>
      </el-menu-item>
      <el-menu-item index="7">
        <span slot="title">
          <nuxt-link to="/admin/links"
            ><i class="iconfont icon-link"></i>
            <span>链接</span>
          </nuxt-link>
        </span>
      </el-menu-item>
      <el-menu-item index="8">
        <span slot="title">
          <nuxt-link to="/admin/settings"
            ><i class="iconfont icon-setting"></i>
            <span>设置</span>
          </nuxt-link>
        </span>
      </el-menu-item>
    </el-menu>

    <nuxt />
  </div>
</template>

<script>
export default {
  head() {
    return {
      title: 'BBS-GO管理'
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
