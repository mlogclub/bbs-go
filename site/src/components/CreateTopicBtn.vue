<template>
  <!-- <client-only> -->
  <el-dropdown placement="bottom" trigger="click" @command="handlePostCommand">
    <el-button type="primary" :icon="Plus"> 发表 </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="(item, i) in modules"
          :key="i"
          :command="item.command"
          :class="item.icon"
        >
          {{ item.name }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
  <!-- </client-only> -->
</template>

<script setup>
import { Plus } from "@element-plus/icons-vue";

const modules = ref([
  {
    command: "tweet",
    name: "发动态",
    icon: "iconfont icon-tweet2",
  },
  {
    command: "topic",
    name: "发帖子",
    icon: "iconfont icon-topic",
  },
  {
    command: "article",
    name: "发文章",
    icon: "iconfont icon-article",
  },
]);

function handlePostCommand(cmd) {
  const router = useRouter();
  if (cmd === "topic") {
    router.push("/topic/create");
  } else if (cmd === "tweet") {
    router.push("/topic/create?type=1");
  } else if (cmd === "article") {
    router.push("/article/create");
  }
}
</script>

<style lang="scss" scoped></style>
