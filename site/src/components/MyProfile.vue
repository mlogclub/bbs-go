<template>
  <div class="widget">
    <div class="widget-header">
      <span>{{ $t("component.myProfile.title") }}</span>
      <div class="slot">
        <nuxt-link to="/user/profile">{{
          $t("component.myProfile.editProfile")
        }}</nuxt-link>
      </div>
    </div>
    <div class="widget-content stable">
      <div class="str">
        <div class="slabel">{{ $t("component.myProfile.nickname") }}</div>
        <div class="svalue">{{ user.nickname }}</div>
      </div>
      <div class="str">
        <div class="slabel">{{ $t("component.myProfile.description") }}</div>
        <div class="svalue">
          {{ user.description }}
        </div>
      </div>
      <div v-if="user.homePage" class="str">
        <div class="slabel">{{ $t("component.myProfile.homePage") }}</div>
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
