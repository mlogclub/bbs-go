<template>
  <section class="main">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        文章正在审核中
      </div>
    </div>
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <article
          class="article-item article-detail"
          itemscope
          itemtype="http://schema.org/BlogPosting"
        >
          <div class="main-content">
            <div class="article-header">
              <div class="article-title-wrapper">
                <h1 class="article-title" itemprop="headline">
                  {{ article.title }}
                </h1>
                <div class="article-manage-menu">
                  <article-manage-menu :article="article" />
                </div>
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
                    :datetime="useFormatDate(article.createTime)"
                    itemprop="datePublished"
                    >{{ usePrettyDate(article.createTime) }}</time
                  >
                </span>
              </div>
            </div>

            <div
              class="article-content content line-numbers"
              itemprop="articleBody"
              v-html="article.content"
            ></div>

            <!--节点、标签-->
            <div class="article-tags">
              <nuxt-link
                v-for="tag in article.tags"
                :key="tag.id"
                :to="'/articles/tag/' + tag.id"
                class="article-tag"
                >#{{ tag.name }}</nuxt-link
              >
            </div>
          </div>
        </article>

        <!-- 评论 -->
        <comment
          :entity-id="article.id"
          :comment-count="article.commentCount"
          entity-type="article"
        />
      </div>
      <div class="right-container">
        <user-info :user="article.user" />
      </div>
    </div>
  </section>
</template>

<script setup>
const route = useRoute();
const { data: article, error } = await useAsyncData(() =>
  useMyFetch(`/api/article/${route.params.id}`)
);

if (error.value) {
  // error.value.cause
  // error.value.message
  throw createError({
    statusCode: 500,
    message: error.value.message || "你访问的页面发生错误!",
  });
}

useHead({
  title: useSiteTitle(article.value.title),
});

const isPending = computed(() => {
  return article.value.status === 2;
});
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

  // 详情页
  &.article-detail {
    .article-header {
      padding: 10px 0;
      border-bottom: 1px solid var(--border-color);
    }

    .article-title-wrapper {
      display: flex;
      .article-title {
        width: 100%;
        color: var(--text-color);
        font-weight: normal;
        overflow: hidden;
        text-overflow: ellipsis;
        font-size: 18px;
        line-height: 30px;
      }
      .article-manage-menu {
        min-width: max-content;
      }
    }

    .article-tags {
      margin-top: 10px;
      .article-tag {
        height: 25px;
        padding: 0 8px;
        display: inline-flex;
        justify-content: center;
        align-items: center;
        border-radius: 12.5px;
        margin-right: 10px;
        background: var(--bg-color2);
        border: 1px solid var(--border-color);
        color: var(--text-color3);
        font-size: 12px;

        &:hover {
          color: var(--text-link-color);
          background: var(--bg-color);
          border: 1px solid var(--border-hover-color);
        }
      }
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
    display: inline-block;
    font-size: 13px;
    padding-top: 6px;

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

      .article-tag {
        height: 25px;
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
