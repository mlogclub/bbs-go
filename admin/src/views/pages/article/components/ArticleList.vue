<template>
  <div class="articles">
    <template v-if="results && results.length">
      <div v-for="item in results" :key="item.id" class="article-item">
        <div class="article-header">
          <a
            class="article-title"
            :href="useSiteUrl(`/article/${item.id}`)"
            target="_blank"
            >{{ item.title }}</a
          >

          <div class="article-status">
            <a-tag v-if="item.status === 1" color="red" size="mini"
              >已删除</a-tag
            >
            <a-tag v-if="item.status === 2" color="blue" size="mini"
              >待审核</a-tag
            >
          </div>
        </div>

        <div class="article-item-info">
          <div class="article-item-main">
            <div class="article-summary">
              {{ item.summary }}
            </div>

            <div class="article-meta">
              <div class="article-meta-left">
                <a
                  :href="useSiteUrl(`/user/${item.user.id}`)"
                  class="article-meta-item"
                  target="_blank"
                >
                  <span>{{ item.user.nickname }}</span>
                </a>
                <span class="article-meta-item"
                  >发布于{{ useFormatDate(item.createTime) }}</span
                >
              </div>

              <div
                v-if="item.tags && item.tags.length > 0"
                class="article-tags"
              >
                <a-tag v-for="tag in item.tags" :key="tag.id" size="mini">{{
                  tag.name
                }}</a-tag>
              </div>
            </div>
          </div>
          <div v-if="item.cover" class="article-item-cover">
            <a-image
              :src="item.cover.url"
              width="150"
              height="90"
              fit="cover"
            />
          </div>
        </div>

        <div class="article-item-actions">
          <a-popconfirm
            v-if="item.status === 0 || item.status === 2"
            content="是否确定删除？"
            @ok="deleteSubmit(item)"
          >
            <a-button
              class="action-item"
              size="mini"
              type="primary"
              status="warning"
              >删除</a-button
            >
          </a-popconfirm>
          <a-popconfirm
            v-if="item.status === 2"
            content="是否确定取消删除？"
            @ok="auditSubmit(item)"
          >
            <a-button
              class="action-item"
              size="mini"
              type="primary"
              status="success"
              >审核通过</a-button
            >
          </a-popconfirm>
        </div>
      </div>
    </template>
    <a-empty v-else />
  </div>
</template>

<script setup>
  defineProps({
    results: {
      type: Array,
      default() {
        return [];
      },
    },
  });

  const emits = defineEmits(['change']);

  const deleteSubmit = async (row) => {
    try {
      await axios.form(
        '/api/admin/article/delete',
        jsonToFormData({ id: row.id })
      );
      useNotificationSuccess('删除成功');
      emits('change');
    } catch (e) {
      useHandleError(e);
    }
  };
  const auditSubmit = async (row) => {
    try {
      await axios.form(
        '/api/admin/article/audit',
        jsonToFormData({ id: row.id })
      );
      useNotificationSuccess('审核成功');
      emits('change');
    } catch (e) {
      useHandleError(e);
    }
  };
</script>

<style lang="less" scoped>
  .articles {
    display: flex;
    flex-direction: column;
    row-gap: 10px;
    .article-item {
      padding: 12px 0;
      border-bottom: 1px solid var(--color-border-1);

      .article-header {
        margin-bottom: 6px;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .article-status {
          right: 0;
          display: flex;
          align-items: center;
          column-gap: 10px;
        }

        .article-title {
          font-size: 18px;
          font-weight: 500;
          line-height: 24px;
          color: var(--color-neutral-10);
          overflow: hidden;
          text-overflow: ellipsis;
          text-decoration: none;
        }
      }

      .article-item-info {
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .article-item-main {
          .article-summary {
            font-size: 14px;
            line-height: 24px;
            color: var(--color-neutral-8);
            overflow: hidden;
            display: -webkit-box;
            -webkit-box-orient: vertical;
            -webkit-line-clamp: 3;
            text-align: justify;
            word-break: break-all;
            text-overflow: ellipsis;
          }

          .article-meta {
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 13px;

            .article-meta-left {
              display: flex;
              align-items: center;
              column-gap: 6px;
              .article-meta-item {
                color: var(--color-neutral-8);
              }
            }

            .article-tags {
              display: flex;
              align-items: center;
              column-gap: 6px;
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

      .article-item-actions {
        margin-top: 10px;
        display: flex;
        align-items: center;
        justify-content: start;
        column-gap: 10px;
      }
    }
  }
</style>
