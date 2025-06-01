import { App } from 'vue';
import MEditor from './components/MEditor.vue';

export {
  MEditor
};

export default {
  install: (app: App) => {
    app.component('MEditor', MEditor);
  }
}; 