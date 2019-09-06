<template>
  <div v-if="activeUsers && activeUsers.length" class="widget">
    <div class="header">
      活跃用户
    </div>
    <div class="content">
      <div v-for="(user, index) in activeUsers" :key="user.id" style="display: flex; margin: 8px;">
        <div v-if="index === 3 || index === 10 || index === 18">
          <ins
            class="adsbygoogle"
            style="display:block"
            data-ad-format="fluid"
            data-ad-layout-key="-ig-s+1x-t-q"
            data-ad-client="ca-pub-5683711753850351"
            data-ad-slot="4728140043"
          />
          <script>
            (adsbygoogle = window.adsbygoogle || []).push({});
          </script>
        </div>
        <a :href="'/user/' + user.id" :title="user.nickname" style="font-weight: bold;">
          <div
            class="avatar has-border is-rounded"
            :style="{backgroundImage:'url(' + user.avatar + ')'}"
          />
        </a>
        <div style="display: block; margin-left: 10px; vertical-align: middle;">
          <div>
            <a
              :href="'/user/' + user.id"
              :title="user.nickname"
              style="font-weight: bold;"
            >{{ user.nickname }}</a>
          </div>
          <div
            style="word-break: break-all; -webkit-line-clamp: 2; overflow: hidden!important; text-overflow: ellipsis; -webkit-box-orient: vertical; display: -webkit-box;"
          >
            {{ user.description }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    activeUsers: {
      type: Array,
      default: function () {
        return []
      }
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
