import Vue from 'vue'
import hljs from 'highlight.js'
import 'highlight.js/styles/monokai.css'

Vue.directive('highlight', function(el) {
  const blocks = el.querySelectorAll('pre code')
  blocks.forEach((block) => {
    hljs.highlightBlock(block)
  })
})
