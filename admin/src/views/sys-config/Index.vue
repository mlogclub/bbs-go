<template>
  <section class="page-container">
    <el-tabs value="first">
      <el-tab-pane label="通用配置" name="first">
        <div class="config">
          <el-form label-width="100px">
            <el-form-item label="网站名称">
              <el-input v-model="config.siteTitle" type="text" placeholder="网站名称"></el-input>
            </el-form-item>

            <el-form-item label="网站描述">
              <el-input
                v-model="config.siteDescription"
                type="textarea"
                autosize
                placeholder="网站描述"
              ></el-input>
            </el-form-item>

            <el-form-item label="网站关键字">
              <el-select
                v-model="config.siteKeywords"
                style="width:100%"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="网站关键字"
              ></el-select>
            </el-form-item>

            <el-form-item label="推荐标签">
              <el-select
                v-model="config.recommendTags"
                style="width:100%"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="推荐标签"
              ></el-select>
            </el-form-item>

            <el-form-item label="论坛导航">
              <el-select
                v-model="config.bbsNavTags"
                style="width:100%"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="论坛导航标签，用于显示在讨论区侧边栏"
              ></el-select>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>
      <el-tab-pane label="导航配置" name="second" class="nav-panel">
        <draggable v-model="config.siteNavs" draggable=".nav" handle=".nav-sort-btn" class="navs">
          <div v-for="(nav, index) in config.siteNavs" :key="index" class="nav">
            <el-row :gutter="20">
              <el-col :span="1">
                <i class="iconfont icon-sort nav-sort-btn" />
              </el-col>
              <el-col :span="10">
                <el-input v-model="nav.title" type="text" size="small" placeholder="标题"></el-input>
              </el-col>
              <el-col :span="11">
                <el-input v-model="nav.url" type="text" size="small" placeholder="链接"></el-input>
              </el-col>
              <el-col :span="2">
                <el-button
                  type="danger"
                  icon="el-icon-delete"
                  circle
                  size="small"
                  @click="delNav(index)"
                ></el-button>
              </el-col>
            </el-row>
          </div>
        </draggable>
        <div class="add-nav">
          <el-tooltip class="item" effect="dark" content="点击按钮添加导航" placement="top">
            <el-button type="primary" icon="el-icon-plus" circle @click="addNav"></el-button>
          </el-tooltip>
        </div>
      </el-tab-pane>
    </el-tabs>

    <div style="margin-top: 20px;">
      <el-button type="primary" :loading="loading" @click="save">保存</el-button>
    </div>
  </section>
</template>

<script>
import HttpClient from "@/apis/HttpClient";
import draggable from "vuedraggable";

export default {
  name: "List",
  components: {
    draggable
  },
  data() {
    return {
      config: {
        siteTitle: "",
        siteDescription: "",
        siteKeywords: [],
        siteNavs: [
          // {
          //   title: "xxx",
          //   url: "/topics"
          // }
        ],
        recommendTags: [],
        bbsNavTags: []
      },
      loading: false
    };
  },
  mounted() {
    this.load();
  },
  methods: {
    load() {
      const me = this;
      HttpClient.get("/api/admin/sys-config/all")
        .then(data => {
          for (let i = 0; i < data.length; i++) {
            const item = data[i];
            if (!me.config.hasOwnProperty(item.key)) {
              continue;
            }
            switch (item.key) {
              case "siteKeywords":
              case "siteNavs":
              case "recommendTags":
              case "bbsNavTags":
                try {
                  this.config[item.key] = JSON.parse(item.value);
                } catch (err) {
                  console.error(err);
                }
                break;
              default:
                this.config[item.key] = item.value;
                break;
            }
          }
        })
        .catch(rsp => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    save() {
      console.log(this.config);
      const me = this;
      me.loading = true;
      HttpClient.post("/api/admin/sys-config/save", {
        config: JSON.stringify(this.config)
      })
        .then(() => {
          me.loading = false;
          me.$message({ message: "提交成功", type: "success" });
          me.load();
        })
        .catch(rsp => {
          me.loading = false;
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    addNav() {
      if (!this.config.siteNavs) {
        this.config.siteNavs = [];
      }
      this.config.siteNavs.push({
        title: "",
        url: ""
      });
    },
    delNav(index) {
      if (!this.config.siteNavs) {
        return;
      }
      this.config.siteNavs.splice(index, 1);
    }
  }
};
</script>

<style scoped lang="scss">
.config {
  padding: 10px 0;
}
.nav-panel {
  .navs {
    border: 1px solid #ddd;
    border-radius: 5px;
    .nav {
      padding: 5px 5px;
      margin: 0;

      &:not(:last-child) {
        border-bottom: 1px solid #ddd;
      }

      .nav-sort-btn {
        font-size: 21px;
        font-weight: 700;
        cursor: pointer;
        float: right;
      }
    }
  }

  .add-nav {
    margin-top: 20px;
    text-align: center;
  }
}
</style>
