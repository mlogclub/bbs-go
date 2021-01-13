import Vue from 'vue'
// import 'element-ui/lib/theme-chalk/index.css'

import {
  Pagination,
  Dialog,
  Input,
  InputNumber,
  Switch,
  Select,
  Option,
  OptionGroup,
  Button,
  ButtonGroup,
  Table,
  TableColumn,
  Tooltip,
  Form,
  FormItem,
  Tabs,
  TabPane,
  Tag,
  Alert,
  Row,
  Col,
  Upload,
  Dropdown,
  DropdownMenu,
  DropdownItem,
  Popover,
  Progress,
  Loading,
  MessageBox,
  Message,
  Notification,
} from 'element-ui'

Vue.use(Pagination)
Vue.use(Dialog)
Vue.use(Input)
Vue.use(InputNumber)
Vue.use(Switch)
Vue.use(Select)
Vue.use(Option)
Vue.use(OptionGroup)
Vue.use(Button)
Vue.use(ButtonGroup)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Tooltip)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Tabs)
Vue.use(TabPane)
Vue.use(Tag)
Vue.use(Alert)
Vue.use(Row)
Vue.use(Col)
Vue.use(Upload)
Vue.use(Dropdown)
Vue.use(Popover)
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
