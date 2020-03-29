import Vue from 'vue'

// Vue.directive('highlight', function(el) {
//   const blocks = el.querySelectorAll('pre code')
//   blocks.forEach((block) => {
//     hljs.highlightBlock(block)
//   })
// })

Vue.directive('paste', function(el, binding) {
  el.onpaste = function(event) {
    binding.value(event)
  }
})

// 复制粘贴指令
// Vue.directive('paste', {
//   bind(el, binding, vnode) {
//     el.addEventListener('paste', function(event) {
//       debugger
//       // 这里直接监听元素的粘贴事件
//       binding.value(event)
//     })
//   }
// })

// 拖拽释放指令
Vue.directive('drag', {
  bind(el, binding, vnode) {
    // 因为拖拽还包括拖动时的经过事件，离开事件，和进入事件，放下事件，浏览器对于拖拽的默认事件的处理是打开拖进来的资源，所以要先对这三个事件进行默认事件的禁止
    el.addEventListener('dragenter', function(event) {
      event.stopPropagation()
      event.preventDefault()
    })
    el.addEventListener('dragover', function(event) {
      event.stopPropagation()
      event.preventDefault()
    })
    el.addEventListener('dragleave', function(event) {
      event.stopPropagation()
      event.preventDefault()
    })
    el.addEventListener('drop', function(event) {
      // 这里阻止默认事件，并绑定事件的对象，用来在组件上返回事件对象
      event.stopPropagation()
      event.preventDefault()
      binding.value(event)
    })
  }
})
