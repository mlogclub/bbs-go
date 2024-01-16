<template>
  <div class="topics-main">
    <load-more-async
      v-slot="{ results }"
      url="/api/topic/topics"
      :params="{ nodeId: nodeId }"
    >
      <topic-list :topics="results" />
    </load-more-async>
  </div>
</template>

<script setup>
definePageMeta({
  layout: "topics",
});
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
    useMyFetch(`/api/topic/node?nodeId=${nodeId}`)
  );
  nodeName = node.value.nodeName;
}

const envStore = useEnvStore();
envStore.setCurrentNodeId(nodeId);

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
