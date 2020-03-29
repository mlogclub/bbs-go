import Vue from 'vue'

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
