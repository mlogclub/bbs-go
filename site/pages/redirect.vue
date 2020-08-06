<template>
  <section class="main">
    <div class="container">
      <div class="main-body">
        <div
          style="text-align: center; vertical-align: center; margin-top: 100px;"
        >
          <div>
            <img src="~/assets/images/logo.png" style="max-width: 100px;" />
          </div>
          <div style="margin-top: 20px;">
            <a :href="url" rel="nofollow"
              >即将跳往站外地址，点击该链接继续跳转&gt;&gt;</a
            >
            <adsbygoogle ad-slot="1742173616" />
          </div>

          <div class="columns recommend">
            <div class="column">
              <div
                v-if="recommendArticles && recommendArticles.length"
                class="widget"
              >
                <div class="widget-header">
                  推荐文章
                </div>
                <div class="widget-content">
                  <ul>
                    <li v-for="a in recommendArticles" :key="a.articleId">
                      <a
                        :href="'/article/' + a.articleId"
                        :title="a.title"
                        target="_blank"
                        >{{ a.title }}</a
                      >
                    </li>
                  </ul>
                </div>
              </div>
            </div>
            <div class="column">
              <div
                v-if="recommendTopics && recommendTopics.length"
                class="widget"
              >
                <div class="widget-header">
                  推荐话题
                </div>
                <div class="widget-content">
                  <ul>
                    <li v-for="t in recommendTopics" :key="t.topicId">
                      <a
                        :href="'/topic/' + t.topicId"
                        :title="t.title"
                        target="_blank"
                        >{{ t.title }}</a
                      >
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  async asyncData({ $axios, query }) {
    const [recommendArticles, recommendTopics] = await Promise.all([
      $axios.get('/api/article/recommend'),
      $axios.get('/api/topic/recommend'),
    ])
    return {
      url: query.url,
      recommendArticles,
      recommendTopics,
    }
  },
}
</script>

<style lang="scss" scoped>
.recommend {
  text-align: left;
  margin: 0;
}
</style>
