<template>
  <ul class="article-list">
    <li v-for="(article, index) in articles" :key="article.articleId">
      <div v-if="showAd && index !== 0 && index % 3 === 0">
        <!-- 信息流广告 -->
        <adsbygoogle
          ad-slot="4980294904"
          ad-format="fluid"
          ad-layout-key="-ht-19-1m-3j+mu"
        />
      </div>
      <article
        class="article-item"
        itemscope
        itemtype="http://schema.org/BlogPosting"
      >
        <h1 class="article-title" itemprop="headline">
          <nuxt-link :to="'/article/' + article.articleId">{{
            article.title
          }}</nuxt-link>
        </h1>

        <div class="article-summary" itemprop="description">
          {{ article.summary }}
        </div>

        <div class="article-meta">
          <span class="article-meta-item">
            由
            <nuxt-link
              :to="'/user/' + article.user.id"
              class="article-author"
              itemprop="author"
              itemscope
              itemtype="http://schema.org/Person"
              ><span itemprop="name">{{
                article.user.nickname
              }}</span></nuxt-link
            >发布于
            <time
              :datetime="article.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')"
              itemprop="datePublished"
              >{{ article.createTime | prettyDate }}</time
            >
          </span>

          <span
            v-if="article.tags && article.tags.length > 0"
            class="article-meta-item"
          >
            <span
              v-for="tag in article.tags"
              :key="tag.tagId"
              class="article-tag tag"
            >
              <nuxt-link :to="'/articles/' + tag.tagId" class>{{
                tag.tagName
              }}</nuxt-link>
            </span>
          </span>

          <span v-if="article.status === 2" class="article-meta-item">
            <a
              href="javascript:void(0)"
              style="
                cursor: default;
                text-decoration: none;
                color: green;
                font-size: 12px;
              "
            >
              <i class="iconfont icon-shenhe" />&nbsp;审核中</a
            >
          </span>
        </div>
      </article>
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    articles: {
      type: Array,
      default() {
        return []
      },
      required: false,
    },
    showAd: {
      type: Boolean,
      default: false,
    },
  },
}
</script>

<style lang="scss" scoped></style>
