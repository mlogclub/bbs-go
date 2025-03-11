<template>
  <client-only>
    <MdEditor
      v-model="value"
      :theme="$colorMode.preference"
      @onChange="change"
      @onUploadImg="uploadImg"
      :toolbars="toolbars"
      :style="{ height: height }"
      :placeholder="placeholder"
      :preview="false"
    >
      <template #defToolbars>
        <NormalToolbar title="切换为HTML编辑器" @onClick="switchHTML">
          <ArrowLeftRight class="md-editor-icon" />
        </NormalToolbar>
      </template>
    </MdEditor>
  </client-only>
</template>

<script setup>
import { MdEditor, NormalToolbar } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import { ArrowLeftRight } from "lucide-vue-next";

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
  height: {
    type: String,
    default: "400px",
  },
  placeholder: {
    type: String,
    default: "请输入...",
  },
});

const emits = defineEmits(["update:modelValue"]);

const value = ref(props.modelValue);

const toolbars = ref([
  "bold",
  "underline",
  "italic",
  "strikeThrough",
  "-",
  "title",
  "sub",
  "sup",
  "quote",
  "unorderedList",
  "orderedList",
  "task",
  "-",
  "codeRow",
  "code",
  "link",
  "image",
  "table",
  // "mermaid",
  // "katex",
  "-",
  "revoke",
  "next",
  // "save",
  "-",
  // "pageFullscreen",
  "preview",
  // "htmlPreview",
  "catalog",
  "=",
  // 0,
  "fullscreen",
]);

function change(v) {
  emits("update:modelValue", v);
}

async function uploadImg(files, callback) {
  const res = await Promise.all(
    files.map((file) => {
      return new Promise((rev, rej) => {
        const formData = new FormData();
        formData.append("image", file, file.name);
        useHttp("/api/upload", {
          method: "POST",
          body: formData,
        })
          .then((res) => rev(res))
          .catch((error) => rej(error));
      });
    })
  );
  callback(res.map((item) => item.url));
}

const switchHTML = () => {
  const eventBus = useEditorEventBus();
  eventBus.emit("switchHtml");
};
</script>

<style lang="scss" scoped></style>
