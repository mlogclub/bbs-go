import htmlEditor from 'html-editor';

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(htmlEditor);
}); 