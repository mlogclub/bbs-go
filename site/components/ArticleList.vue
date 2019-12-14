<template>
  <ul class="article-list">
    <li v-for="(article, index) in articles" :key="article.articleId">
      <div v-if="showAd && (index === 2 || index === 10 || index === 18)">
        <!-- 信息流广告 -->
        <adsbygoogle
          ad-slot="4980294904"
          ad-format="fluid"
          ad-layout-key="-ht-19-1m-3j+mu"
        />
      </div>
      <article class="article-item">
        <div class="article-title">
          <a :href="'/article/' + article.articleId">{{ article.title }}</a>
        </div>

        <div class="article-summary">
          {{ article.summary }}
        </div>

        <div class="article-meta">
          <span class="article-meta-item">
            由
            <a :href="'/user/' + article.user.id" class="article-author"
              >&nbsp;{{ article.user.nickname }}&nbsp;</a
            >发布于
            <time itemprop="datePublished">{{
              article.createTime | prettyDate
            }}</time>
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
              <a :href="'/articles/' + tag.tagId" class>{{ tag.tagName }}</a>
            </span>
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
      required: false
    },
    showAd: {
      type: Boolean,
      default: false
    }
  }
}
</script>

<style lang="scss" scoped></style>
