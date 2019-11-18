<template>
  <ul class="article-list">
    <li v-for="(article, index) in articles" :key="article.articleId">
      <div
        v-if="
          showAd &&
            ((articles.length < 3 && index === 1) ||
              (index !== 0 && index % 5 === 0))
        "
      >
        <ins
          class="adsbygoogle"
          style="display:block"
          data-ad-format="fluid"
          data-ad-layout-key="-ig-s+1x-t-q"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="4728140043"
        />
        <script>
          ;(adsbygoogle = window.adsbygoogle || []).push({})
        </script>
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

          <span v-if="article.category" class="article-meta-item">
            <span class="article-tag tag">
              <a :href="'/articles/cat/' + article.category.categoryId">{{
                article.category.categoryName
              }}</a>
            </span>
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
              <a :href="'/articles/tag/' + tag.tagId" class>{{
                tag.tagName
              }}</a>
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
