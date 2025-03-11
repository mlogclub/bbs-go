<template>
  <section class="main">
    <div class="container main-container">
      <div class="main-content no-padding no-bg topics-wrapper">
        <div class="topics-nav">
          <topics-nav />
        </div>
        <div class="topics-main">
          <!-- <div class="topics-main-header">
            <div>全部帖子</div>
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
const route = useRoute();

let nodeId = Number.parseInt(route.params.id) || 0;
let nodeName = "";
if (route.params.id === "newest") {
  nodeId = 0;
  nodeName = "最新";
} else if (route.params.id === "recommend") {
  nodeId = -1;
  nodeName = "推荐";
} else if (route.params.id === "feed") {
  nodeId = -2;
  nodeName = "关注";
} else {
  const { data: node } = await useAsyncData(() =>
    useHttpGet(`/api/topic/node?nodeId=${nodeId}`)
  );
  nodeName = node.value.name;
}

onMounted(() => {
  useEnvStore().setCurrentNodeId(nodeId);
});

useHead({
  title: useSiteTitle(nodeName, "话题"),
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
