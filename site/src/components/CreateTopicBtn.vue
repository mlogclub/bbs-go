<template>
  <el-dropdown
    v-if="modules.length > 0"
    placement="bottom"
    trigger="click"
    @command="handlePostCommand"
  >
    <el-button type="primary" :icon="Plus">
      {{ $t("common.createBtn.create") }}
    </el-button>
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
const { t } = useI18n();

const configStore = useConfigStore();

let modules = [];
if (configStore.config.modules.tweet) {
  modules.push({
    command: "tweet",
    name: t("common.createBtn.tweet"),
    icon: "icon-tweet2",
  });
}
if (configStore.config.modules.topic) {
  modules.push({
    command: "topic",
    name: t("common.createBtn.topic"),
    icon: "icon-topic",
  });
}
if (configStore.config.modules.article) {
  modules.push({
    command: "article",
    name: t("common.createBtn.article"),
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
