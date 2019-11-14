<template>
  <ul class="article-list">
    <li v-for="(article, index) in articles" :key="article.articleId">
      <div
        v-if="showAd && ((articles.length < 3 && index === 1) || (index !== 0 && index % 5 === 0))"
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
          (adsbygoogle = window.adsbygoogle || []).push({});
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
            <a :href="'/user/' + article.user.id">{{ article.user.nickname }}</a>
          </span>

          <span v-if="article.category" class="article-meta-item">
            <a
              :href="'/articles/cat/' + article.category.categoryId"
            >{{ article.category.categoryName }}</a>
          </span>

          <span v-for="tag in article.tags" :key="tag.tagId" class="article-meta-item">
            <a :href="'/articles/tag/' + tag.tagId">{{ tag.tagName }}</a>
          </span>

          <span class="article-meta-item">
            <time itemprop="datePublished">{{ article.createTime | prettyDate }}</time>
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
      default: function () {
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

<style lang="scss" scoped>
.article-list {
  margin: 0 !important;

  li {
    padding: 8px 0;
    zoom: 1;
    position: relative;
    overflow: hidden;

    &:not(:last-child) {
      border-bottom: 1px dashed #f2f2f2;
    }
  }

  .article-item {
    overflow: hidden;
    zoom: 1;
    line-height: 24px;
  }

  article {
    .article-title {
      a {
        font-size: 18px;
        line-height: 30px;
        font-weight: 500;
        color: #17181a;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    .article-summary {
      font-size: 14px;
      color: rgba(0, 0, 0, 0.7);
      overflow: hidden;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 3;
      text-align: justify;
      padding-top: 6px;
      word-break: break-all;
      text-overflow: ellipsis;
    }

    .article-meta {
      display: inline-block;
      font-size: 13px;
      padding-top: 6px;

      .article-meta-item {
        padding: 0 6px 0 0;
      }

      a {
        color: #3273dc;
      }

      span {
        color: #999;
      }
    }
  }
}
</style>
