<template>
  <div>
    <el-menu
      class="el-menu-demo"
      mode="horizontal"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
    >
      <el-menu-item index="1">
        <i class="iconfont icon-dashboard"></i>
        <span slot="title"><nuxt-link to="/admin">后台</nuxt-link></span>
      </el-menu-item>
      <el-menu-item index="2">
        <i class="iconfont icon-tags"></i>
        <span slot="title"><nuxt-link to="/admin/nodes">节点</nuxt-link></span>
      </el-menu-item>
      <el-menu-item index="3">
        <i class="iconfont icon-topic"></i>
        <span slot="title"><nuxt-link to="/admin/topics">话题</nuxt-link></span>
      </el-menu-item>
      <el-menu-item index="4">
        <i class="iconfont icon-article"></i>
        <span slot="title"
          ><nuxt-link to="/admin/articles">文章</nuxt-link></span
        >
      </el-menu-item>
      <el-menu-item index="5">
        <i class="iconfont icon-comment"></i>
        <span slot="title"
          ><nuxt-link to="/admin/comments">评论</nuxt-link></span
        >
      </el-menu-item>
      <el-menu-item index="6">
        <i class="iconfont icon-username"></i>
        <span slot="title"><nuxt-link to="/admin/users">用户</nuxt-link></span>
      </el-menu-item>
      <el-menu-item index="7">
        <i class="iconfont icon-setting"></i>
        <span slot="title"
          ><nuxt-link to="/admin/settings">设置</nuxt-link></span
        >
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
  }
}
</style>
