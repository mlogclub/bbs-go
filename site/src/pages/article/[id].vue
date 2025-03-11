<template>
  <section class="main">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        文章正在审核中
      </div>
    </div>
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <article class="article-detail">
          <div class="article-header">
            <div class="article-title-wrapper">
              <h1 class="article-title">
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
                  >{{ article.user.nickname }}</nuxt-link
                >发布于
                <time :datetime="useFormatDate(article.createTime)">{{
                  usePrettyDate(article.createTime)
                }}</time>
              </span>
            </div>
          </div>

          <div
            class="article-content content line-numbers"
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
  useHttpGet(`/api/article/${route.params.id}`)
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
.article-detail {
  margin-bottom: 12px;
  padding: 12px;
  border-radius: var(--border-radius);
  background-color: var(--bg-color);
  overflow: hidden;

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

  .article-header {
    padding: 10px 0;
    border-bottom: 1px solid var(--border-color);

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
</style>
