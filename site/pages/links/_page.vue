<template>
  <section class="main">
    <div class="container">
      <div class="main-body">
        <!--
        <div class="notice">
          <h1>什么是好博客导航？</h1>
          <p>
            好博客导航是一个收录优质、原创、计算机技术相关博客导航工具。<a
              href="/link/submit"
              >点击提交你的博客链接&gt;&gt;</a
            >
          </p>
          <h1>为什么要做`好博客`导航？</h1>
          <p>
            我在网上看到过很多博客导航，但是收录的博客质量参差不齐，而且没有专业编程相关的技术类型博客导航，有很多优质好博客没有得到很好的展示机会，好博客导航主要就是为了解决一问题，让独立博主能够很好的展示自己，让自己的文章能够帮助更多人，让更多的程序员能够关注到自己喜欢的博客。
          </p>
          <p>
            后续我们还会对所有收录的博客进行分类、打标签，对优质博客进行推荐。
          </p>
        </div>
        -->

        <!-- 展示广告 -->
        <!--        <adsbygoogle ad-slot="1742173616" />-->

        <div class="widget">
          <div class="widget-header">友情链接</div>
          <div class="widget-content">
            <ul class="links">
              <li
                v-for="link in linksPage.results"
                :key="link.linkId"
                class="link"
              >
                <div class="link-logo">
                  <img v-if="link.logo" :src="link.logo" />
                  <img v-if="!link.logo" src="~/assets/images/net.png" />
                </div>
                <div class="link-content">
                  <a
                    :href="'/link/' + link.linkId"
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
            <pagination :page="linksPage.page" url-prefix="/links/" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import Pagination from '~/components/Pagination'
export default {
  components: {
    Pagination
  },
  async asyncData({ $axios, params }) {
    const [linksPage] = await Promise.all([
      $axios.get('/api/link/links', {
        params: {
          page: params.page || 1
        }
      })
    ])
    return {
      linksPage
    }
  },
  head() {
    return {
      title: this.$siteTitle('好博客'),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription()
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  }
}
</script>

<style lang="scss" scoped></style>
