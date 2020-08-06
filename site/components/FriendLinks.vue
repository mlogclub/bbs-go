<template>
  <div v-if="links && links.length" class="widget">
    <div class="widget-header">
      <span>友情链接</span>
      <span class="slot"><a href="/links">查看更多&gt;&gt;</a></span>
    </div>
    <div class="widget-content">
      <ul class="links">
        <li v-for="link in links" :key="link.linkId" class="link">
          <div class="link-logo">
            <img v-if="link.logo" :src="link.logo" />
            <img v-if="!link.logo" src="~/assets/images/net.png" />
          </div>
          <div class="link-content">
            <a
              :href="link.url"
              :title="link.title"
              class="link-title"
              target="_blank"
              >{{ link.title }}</a
            >
            <p class="link-summary">
              {{ link.summary }}
            </p>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FriendLinks',
  props: {
    links: {
      type: Array,
      default() {
        return null
      },
    },
  },
}
</script>

<style scoped lang="scss"></style>
