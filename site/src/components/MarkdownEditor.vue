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
      :preview="true"
      :language="language"
      :footers="[]"
    >
    </MdEditor>
  </client-only>
</template>

<script setup>
import { MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/style.css";

const language = computed(() => {
  const { locale } = useI18n();
  return locale.value;
});

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
    default: "",
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
      return useUploadImage(file);
    })
  );
  callback(res.map((item) => item.url));
}
</script>

<style lang="scss" scoped></style>
