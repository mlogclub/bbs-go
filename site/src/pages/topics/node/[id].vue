<template>
  <section class="main">
    <div class="container main-container">
      <div class="main-content no-padding no-bg topics-wrapper">
        <div class="topics-nav">
          <topics-nav />
        </div>
        <div class="topics-main">
          <!-- <div class="topics-main-header">
            <div>ALL</div>
          </div> -->
          <load-more-async
            v-slot="{ results }"
            url="/api/topic/topics"
            :params="{ nodeId: nodeId }"
          >
            <topic-list :topics="results" show-sticky />
          </load-more-async>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const { t } = useI18n();
const route = useRoute();

let nodeId = Number.parseInt(route.params.id) || 0;
if (route.params.id === "newest") {
  nodeId = 0;
} else if (route.params.id === "recommend") {
  nodeId = -1;
} else if (route.params.id === "feed") {
  nodeId = -2;
}

const { data: node } = await useMyFetch(`/api/topic/node?nodeId=${nodeId}`);

onMounted(() => {
  useEnvStore().setCurrentNodeId(nodeId);
});

useHead({
  title: useSiteTitle(node.value.name, t("pages.topics.title")),
  meta: [
    {
      hid: "description",
      name: "description",
      content: useSiteDescription(),
    },
    { hid: "keywords", name: "keywords", content: useSiteKeywords() },
  ],
});
</script>

<style lang="scss" scoped></style>
