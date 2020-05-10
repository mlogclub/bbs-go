<template>
  <div v-if="user" class="widget">
    <div class="widget-header">
      个人信息
    </div>
    <div class="widget-content">
      <img :src="user.avatar" class="img-avatar" />
      <div class="nickname">
        <a :href="'/user/' + user.id">{{ user.nickname }}</a>
      </div>
      <div v-if="user.description" class="description">
        <p>{{ user.description }}</p>
      </div>
      <div v-if="user.type === 1">
        <img
          :src="'https://open.weixin.qq.com/qr/code?username=' + user.username"
        />
      </div>
    </div>
  </div>
</template>
``
<script>
export default {
  props: {
    user: {
      type: Object,
      default() {
        return null
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.nickname {
  font-size: 18px;
  font-weight: bold;
  a {
    color: #3273dc;
  }
}
.img-avatar {
  margin-top: 5px;
  border: 1px dotted #eeeeee;
  border-radius: 5%;
}
.description {
  font-size: 14px;
  padding: 10px 15px;
  border: 1px dotted #eeeeee;
  border-left: 3px solid #eeeeee;
  background-color: #fbfbfb;
}
</style>
