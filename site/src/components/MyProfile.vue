<template>
  <div class="widget">
    <div class="widget-header">
      <span>个人资料</span>
      <div class="slot">
        <nuxt-link to="/user/profile">编辑资料</nuxt-link>
      </div>
    </div>
    <div class="widget-content stable">
      <div class="str">
        <div class="slabel">昵称</div>
        <div class="svalue">{{ user.nickname }}</div>
      </div>
      <div class="str">
        <div class="slabel">签名</div>
        <div class="svalue">
          {{ user.description }}
        </div>
      </div>
      <div v-if="user.homePage" class="str">
        <div class="slabel">主页</div>
        <div class="svalue">
          <a :href="user.homePage" target="_blank" rel="nofollow">{{
            user.homePage
          }}</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "MyProfile",
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style scoped></style>
