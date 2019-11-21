<template>
  <div class="left-container">
    <div class="widget">
      <div class="widget-header">
        <a :href="createTopicUrl" class="button is-success post">
          <i class="iconfont icon-topic" />&nbsp;
          <strong>发表主题</strong>
        </a>
      </div>
      <div class="widget-content">
        <ul class="m-nav">
          <li>
            <a :class="{ active: currentTagId == 0 }" href="/topics">
              <i class="iconfont icon-topic" />
              <span>全部话题</span>
            </a>
          </li>
          <li v-for="tag in bbsNavTags" :key="tag.tagId">
            <a
              :class="{ active: currentTagId == tag.tagId }"
              :href="'/topics/tag/' + tag.tagId"
            >
              <i class="iconfont icon-topic" />
              <span>{{ tag.tagName }}</span>
            </a>
          </li>
        </ul>
      </div>
    </div>

    <weixin-gzh />

    <!-- 展示广告190x90 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:90px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="9345305153"
    />
    <script>
      ;(adsbygoogle = window.adsbygoogle || []).push({})
    </script>

    <!-- 展示广告190x480 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:480px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="3438372357"
    />
    <script>
      ;(adsbygoogle = window.adsbygoogle || []).push({})
    </script>
  </div>
</template>

<script>
import WeixinGzh from '~/components/WeixinGzh'

export default {
  components: { WeixinGzh },
  props: {
    currentTagId: {
      type: Number,
      default: 0
    }
  },
  computed: {
    createTopicUrl() {
      let url = '/topic/create'
      if (this.currentTagId) {
        url += '?tagId=' + this.currentTagId
      }
      return url
    },
    bbsNavTags() {
      return this.$store.state.config.config.bbsNavTags
    }
  }
}
</script>

<style lang="scss" scoped></style>
