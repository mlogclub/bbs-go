<template>
  <section class="main">
    <div class="container">
      <user-profile :user="user" />

      <div class="container main-container right-main size-320">
        <user-center-sidebar :user="user" />
        <div class="right-container">
          <div class="tabs-warp">
            <div class="tabs">
              <ul>
                <li :class="{ 'is-active': activeTab === 'topics' }">
                  <nuxt-link :to="'/user/' + user.id">
                    <span class="icon is-small">
                      <i class="iconfont icon-topic" aria-hidden="true" />
                    </span>
                    <span>话题</span>
                  </nuxt-link>
                </li>
                <li :class="{ 'is-active': activeTab === 'articles' }">
                  <nuxt-link :to="'/user/' + user.id + '/articles'">
                    <span class="icon is-small">
                      <i class="iconfont icon-article" aria-hidden="true" />
                    </span>
                    <span>文章</span>
                  </nuxt-link>
                </li>
              </ul>
            </div>

            <div>
              <load-more-async
                v-slot="{ results }"
                url="/api/topic/user/topics"
                :params="{ userId: user.id }"
              >
                <topic-list :topics="results" :show-avatar="false" />
              </load-more-async>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const route = useRoute();
const user = await useHttpGet(`/api/user/${route.params.userId}`);
const activeTab = ref("topics");
useHead({
  title: useSiteTitle("个人主页", user.nickname),
});
</script>

<style lang="scss" scoped>
.tabs-warp {
  background-color: var(--bg-color);
  padding: 0 10px 10px;

  .tabs {
    margin-bottom: 5px;
  }

  .more {
    text-align: right;
  }
}
</style>
