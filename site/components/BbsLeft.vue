<template>
  <div class="m-left">
    <a class="button is-primary post" :href="createTopicUrl">
      <i class="iconfont icon-topic" />&nbsp;
      <strong>发表主题</strong>
    </a>
    <ul class="m-nav">
      <li>
        <a :class="{'active': currentTagId == 0}" href="/topics">
          <i class="iconfont icon-topic" />
          <span class="Button-label">全部话题</span>
        </a>
      </li>
      <li>
        <a :class="{'active': currentTagId == 16310}" href="/topics/tag/16310">
          <i class="iconfont icon-topic" />
          <span class="Button-label">分享</span>
        </a>
      </li>
      <li>
        <a :class="{'active': currentTagId == 15533}" href="/topics/tag/15533">
          <i class="iconfont icon-topic" />
          <span class="Button-label">提问</span>
        </a>
      </li>
      <li>
        <a :class="{'active': currentTagId == 15542}" href="/topics/tag/15542">
          <i class="iconfont icon-topic" />
          <span class="Button-label">开源</span>
        </a>
      </li>
      <li>
        <a :class="{'active': currentTagId == 37}" href="/topics/tag/37">
          <i class="iconfont icon-topic" />
          <span class="Button-label">Go语言</span>
        </a>
      </li>
      <li>
        <a :class="{'active': currentTagId == 2}" href="/topics/tag/2">
          <i class="iconfont icon-topic" />
          <span class="Button-label">Java</span>
        </a>
      </li>
      <li>
        <a :class="{'active': currentTagId == 4035}" href="/topics/tag/4035">
          <i class="iconfont icon-topic" />
          <span class="Button-label">JavaScript</span>
        </a>
      </li>
      <li>
        <a :class="{'active': currentTagId == 191}" href="/topics/tag/191">
          <i class="iconfont icon-topic" />
          <span class="Button-label">Python</span>
        </a>
      </li>
    </ul>

    <!-- 展示广告190x90 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:90px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="9345305153"
    />
    <script>
      (adsbygoogle = window.adsbygoogle || []).push({});
    </script>

    <!-- 展示广告190x480 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:480px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="3438372357"
    />
    <script>
      (adsbygoogle = window.adsbygoogle || []).push({});
    </script>
  </div>
</template>

<script>
export default {
  props: {
    currentTagId: {
      type: Number,
      default: 0
    }
  },
  computed: {
    createTopicUrl: function () {
      let url = '/topic/create'
      if (this.currentTagId) {
        url += '?tagId=' + this.currentTagId
      }
      return url
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
