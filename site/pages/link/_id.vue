<template>
  <section class="main">
    <div class="container">
      <div class="main-body">
        <div class="link">
          <div class="logo" />
          <div class="title">
            <img v-if="link.logo" :src="link.logo" />
            <img v-else src="~/assets/images/net.png" />
            {{ link.title }}
          </div>
          <div class="summary">
            {{ link.summary }}
          </div>
          <div class="link">
            博客地址：<a :href="link.url">{{ link.url }}</a>
          </div>
        </div>
        <div style="margin-top: 20px;">
          <!-- 展示广告 -->
          <adsbygoogle ad-slot="1742173616" />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  async asyncData({ $axios, params }) {
    const link = await $axios.get('/api/link/' + params.id)
    return {
      link,
    }
  },
  head() {
    const title = this.link.title + ' - 好博客'
    return {
      title: this.$siteTitle(title),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription(),
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() },
      ],
    }
  },
}
</script>

<style lang="scss" scoped>
.link {
  margin-top: 20px;

  .title {
    display: flex;
    line-height: 100px;
    img {
      width: 100px;
      height: 100px;
      margin-right: 20px;
    }
  }

  .summary {
    font-size: 16px;
    padding: 10px 15px;
    border: 1px dotted #eeeeee;
    border-left: 3px solid #eeeeee;
    background-color: #fbfbfb;
  }
  .link {
    font-size: 16px;

    a {
      color: #3273dc;
    }
  }
}
</style>
