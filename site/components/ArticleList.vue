<template>
  <ul class="article-list">
    <li v-for="article in articles" :key="article.articleId">
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
          <div class="article-meta-left">
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
                :datetime="
                  article.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                "
                itemprop="datePublished"
                >{{ article.createTime | prettyDate }}</time
              >
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

          <div class="article-meta-right">
            <span
              v-if="article.tags && article.tags.length > 0"
              class="article-meta-item"
            >
              <span
                v-for="tag in article.tags"
                :key="tag.tagId"
                class="article-tag"
              >
                <nuxt-link :to="'/articles/' + tag.tagId" class>{{
                  tag.tagName
                }}</nuxt-link>
              </span>
            </span>
          </div>
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
  },
}
</script>

<style lang="scss" scoped>
.article-list {
  margin: 0 !important;

  li {
    padding: 12px 12px;
    display: flex;
    position: relative;
    overflow: hidden;
    transition: background 0.5s;
    border-radius: 3px;
    background: var(--bg-color);

    &:not(:last-child) {
      margin-bottom: 10px;
    }
  }
}

.article-item {
  overflow: hidden;
  zoom: 1;
  line-height: 24px;

  .article-title {
    a {
      font-size: 18px;
      line-height: 30px;
      font-weight: 500;
      color: var(--text-color);
      overflow: hidden;
      text-overflow: ellipsis;
    }
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

        a {
          color: var(--text-link-color);

          &.article-author {
            font-weight: bold;
            padding: 0 3px;
          }
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

  .article-tool {
    display: inline-block;
    margin-right: 5px;
    line-height: 32px;

    & > span {
      margin-left: 5px;

      a {
        font-size: 12px;
        color: var(--text-color3);
        font-weight: 700;

        &:hover {
          text-decoration: underline;
        }

        i {
          font-size: 12px;
          color: var(--text-color);
        }
      }
    }
  }

  .article-content {
    font-size: 15px;
    margin-top: 10px;
    margin-bottom: 10px;

    a.article-share-summary {
      color: var(--text-color);
    }
  }

  .article-footer {
    word-break: break-all;
    background: var(--bg-color);
    padding: 10px;

    &,
    a {
      color: var(--text-color);
      font-size: 14px;
    }
  }
}

.article-related {
  margin-top: 0 !important;

  li {
    // margin: 8px 0;
    padding: 5px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }
  }

  .article-related-title {
    overflow: hidden;
    word-break: break-all;
    text-overflow: ellipsis;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    display: -webkit-box;

    color: var(--text-color2);
    font-size: 14px;

    &:hover {
      color: var(--text-link-color);
    }
  }
}
</style>
