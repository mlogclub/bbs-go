<template>
  <el-dropdown
    v-if="modules.length > 0"
    placement="bottom"
    trigger="click"
    @command="handlePostCommand"
  >
    <el-button type="primary" :icon="Plus"> 发表 </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="(item, i) in modules"
          :key="i"
          :command="item.command"
        >
          <i class="iconfont" :class="item.icon"></i>
          <span>{{ item.name }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { Plus } from "@element-plus/icons-vue";

const configStore = useConfigStore();

let modules = [];
if (configStore.config.modules.tweet) {
  modules.push({
    command: "tweet",
    name: "发动态",
    icon: "icon-tweet2",
  });
}
if (configStore.config.modules.topic) {
  modules.push({
    command: "topic",
    name: "发帖子",
    icon: "icon-topic",
  });
}
if (configStore.config.modules.article) {
  modules.push({
    command: "article",
    name: "发文章",
    icon: "icon-article",
  });
}

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
