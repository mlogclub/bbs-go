import Vue from 'vue'
// import 'element-ui/lib/theme-chalk/index.css'

import {
  Dialog,
  Button,
  Alert,
  Dropdown,
  DropdownMenu,
  DropdownItem,
  Progress,
  Loading,
  MessageBox,
  Message,
  Notification,
} from 'element-ui'

Vue.use(Dialog)
Vue.use(Button)
Vue.use(Alert)
Vue.use(Dropdown)
Vue.use(DropdownMenu)
Vue.use(DropdownItem)
Vue.use(Progress)

Vue.use(Loading.directive)

Vue.prototype.$loading = Loading.service
Vue.prototype.$msgbox = MessageBox
Vue.prototype.$alert = MessageBox.alert
Vue.prototype.$confirm = MessageBox.confirm
Vue.prototype.$prompt = MessageBox.prompt
Vue.prototype.$notify = Notification
Vue.prototype.$message = Message
