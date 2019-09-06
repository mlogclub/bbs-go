<template>
  <quill-editor :value="value" :options="options"
                @blur="onEditorBlur($event)"
                @focus="onEditorFocus($event)"
                @ready="onEditorReady($event)"
                @input="onInput"></quill-editor>
</template>

<script>

import 'quill/dist/quill.core.css';
import 'quill/dist/quill.snow.css';
import 'quill/dist/quill.bubble.css';
import { quillEditor } from 'vue-quill-editor';
import hljs from 'highlight.js';

export default {
  name: 'HtmlEditor',
  components: { quillEditor },
  props: {
    value: {
      type: String,
      default: '',
    },
    height: {
      type: Number,
      default: 350,
    },
  },
  data() {
    return {
      options: {
        theme: 'snow',
        // theme: 'bubble',
        placeholder: '请输入内容',
        modules:
            {
              toolbar: [
                ['bold', 'italic', 'underline', 'strike'],
                ['blockquote', 'code-block'],
                [{ list: 'ordered' }, { list: 'bullet' }],
                ['link', 'image', 'video'],
                ['clean'],
              ],
              syntax: {
                highlight: text => hljs.highlightAuto(text).value,
              },
            },
      },
    };
  },
  mounted() {
  },
  methods: {
    onEditorBlur(editor) {
      this.$emit('blur', editor);
    },
    onEditorFocus(editor) {
      this.$emit('focus', editor);
    },
    onEditorReady(editor) {
      this.$emit('ready', editor);
    },
    onInput(_value) {
      this.$emit('input', _value);
    },
  },
};
</script>

<style scoped>

</style>
