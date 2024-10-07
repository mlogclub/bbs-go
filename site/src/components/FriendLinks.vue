<template>
  <div v-if="links && links.length" class="widget">
    <div class="widget-header">
      <span>友情链接</span>
      <span class="slot"
        ><nuxt-link to="/links">查看更多&gt;&gt;</nuxt-link></span
      >
    </div>
    <div class="widget-content">
      <ul class="links">
        <li v-for="link in links" :key="link.linkId" class="link">
          <div class="link-logo">
            <img v-if="link.logo" :src="link.logo" />
            <img v-if="!link.logo" src="~/assets/images/net.png" />
          </div>
          <div class="link-content">
            <a
              :href="link.url"
              :title="link.title"
              class="link-title"
              target="_blank"
              >{{ link.title }}</a
            >
            <p class="link-summary">
              {{ link.summary }}
            </p>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
const { data: links } = useAsyncData(() => useMyFetch("/api/link/toplinks"));
</script>

<style scoped lang="scss">
.links {
  .link {
    display: flex;
    height: 62px;
    padding-top: 5px;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color4);
    }

    .link-logo {
      display: inline-block;
      min-width: 50px;
      min-height: 50px;
      img {
        max-width: 50px;
        max-height: 50px;
        border-radius: 50%;
      }
    }

    .link-content {
      display: block;
      margin-left: 5px;

      .link-title {
        font-size: 15px;
        font-weight: 600;
        color: var(--text-link-color);

        overflow: hidden;
        word-break: break-all;
        -webkit-line-clamp: 1;
        text-overflow: ellipsis;
        -webkit-box-orient: vertical;
        display: -webkit-box;
      }

      .link-summary {
        font-size: 14px;
        margin-top: 3px;
        // font-weight: 500;

        overflow: hidden;
        word-break: break-all;
        -webkit-line-clamp: 1;
        text-overflow: ellipsis;
        -webkit-box-orient: vertical;
        display: -webkit-box;
      }
    }
  }
}
</style>
