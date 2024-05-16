<template>
  <div class="article-list">
    <div v-for="article in articles" :key="article.id" class="article-item">
      <div class="article-item-main">
        <div class="article-info">
          <nuxt-link
            class="article-title"
            :to="'/article/' + article.id"
            target="_blank"
            >{{ article.title }}</nuxt-link
          >

          <div class="article-summary">
            {{ article.summary }}
          </div>
        </div>

        <div class="article-meta">
          <div class="article-meta-left">
            <span class="article-meta-item">
              <nuxt-link
                :to="'/user/' + article.user.id"
                class="article-author"
              >
                <span>{{ article.user.nickname }}</span>
              </nuxt-link>
              <time :datetime="useFormatDate(article.createTime)"
                >发布于 {{ usePrettyDate(article.createTime) }}</time
              >
            </span>
          </div>

          <div class="article-meta-right">
            <div v-if="article.tags && article.tags.length > 0">
              <nuxt-link
                v-for="tag in article.tags"
                :key="tag.id"
                class="article-tag"
                :to="'/articles/tag/' + tag.id"
                >{{ tag.name }}</nuxt-link
              >
            </div>
          </div>
        </div>
      </div>
      <div v-if="article.cover" class="article-item-cover">
        <img :src="article.cover.url" />
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  articles: {
    type: Array,
    default() {
      return [];
    },
    required: false,
  },
});
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

    display: flex;
    align-items: center;

    .article-item-main {
      width: 100%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;

      .article-info {
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
      }

      .article-meta {
        display: flex;
        // justify-content: space-between;
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
          @media screen and (max-width: 768px) {
            & {
              display: none;
            }
          }

          margin-left: 10px;

          .article-tag {
            padding: 2px 8px;
            justify-content: center;
            align-items: center;
            border-radius: 12.5px;
            margin-right: 10px;
            background: var(--bg-color2);
            border: 1px solid var(--border-color2);
            color: var(--text-color3);
            font-size: 12px;

            &:hover {
              color: var(--text-link-color);
              background: var(--bg-color);
            }
          }
        }
      }
    }

    .article-item-cover {
      display: flex;
      margin-left: 6px;
      img {
        min-width: 140px;
        min-height: 90px;
        width: 140px;
        height: 90px;
        object-fit: cover;

        @media screen and (max-width: 768px) {
          & {
            min-width: 110px;
            min-height: 80px;
            width: 110px;
            height: 80px;
          }
        }
      }
    }
  }
}
</style>
