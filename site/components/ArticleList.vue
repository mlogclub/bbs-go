<template>
  <div class="article-list">
    <div
      v-for="article in articles"
      :key="article.articleId"
      class="article-item"
    >
      <nuxt-link class="article-title" :to="'/article/' + article.articleId">{{
        article.title
      }}</nuxt-link>

      <div class="article-summary">
        {{ article.summary }}
      </div>

      <div class="article-meta">
        <div class="article-meta-left">
          <span class="article-meta-item">
            <nuxt-link :to="'/user/' + article.user.id" class="article-author">
              <span>{{ article.user.nickname }}</span>
            </nuxt-link>
            <time
              :datetime="article.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')"
              >发布于 {{ article.createTime | prettyDate }}</time
            >
          </span>
        </div>

        <div class="article-meta-right">
          <div v-if="article.tags && article.tags.length > 0">
            <span
              v-for="tag in article.tags"
              :key="tag.tagId"
              class="article-tag"
            >
              <nuxt-link :to="'/articles/' + tag.tagId" class>{{
                tag.tagName
              }}</nuxt-link>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
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
  },
}
</script>

<style lang="scss" scoped>
.article-list {
  margin: 0 !important;

  .article-item {
    padding: 12px 12px;
    transition: background 0.5s;
    border-radius: 3px;
    background: var(--bg-color);
    line-height: 24px;

    &:not(:last-child) {
      margin-bottom: 10px;
    }

    .article-title {
      font-size: 18px;
      line-height: 30px;
      font-weight: 500;
      color: var(--text-color);
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .article-summary {
      font-size: 14px;
      color: var(--text-color2);
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
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 13px;
      padding-top: 6px;

      .article-meta-left {
        .article-meta-item {
          padding: 0 6px 0 0;
          color: var(--text-color3);

          .article-author {
            font-weight: bold;
            padding: 0 3px;
          }
        }
      }

      .article-meta-right {
        .article-tag {
          height: 22px;
          padding: 0 8px;
          display: inline-flex;
          justify-content: center;
          align-items: center;
          border-radius: 12.5px;
          background: var(--bg-color2);
          border: 1px solid var(--border-color);
          color: var(--text-color3);
          font-size: 12px;

          &:hover {
            color: var(--text-link-color);
            background: var(--bg-color);
            border: 1px solid var(--border-hover-color);
          }

          &:not(:last-child) {
            margin-right: 10px;
          }
        }
      }

      @media screen and (max-width: 1024px) {
        .article-meta-right {
          display: none;
        }
      }
    }
  }
}
</style>
