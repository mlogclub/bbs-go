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
                v-model="config.bbsNavTagIds"
                style="width:100%"
                multiple
                filterable
                remote
                placeholder="论坛导航标签，用于显示在讨论区侧边栏"
                :remote-method="loadAutocompleteTags"
                :loading="autocompleteTagLoading"
              >
                <!-- 已经选择的 -->
                <template v-if="config.bbsNavTags && config.bbsNavTags.length">
                  <el-option
                    v-for="tag in config.bbsNavTags"
                    :key="tag.id"
                    :value="tag.tagId"
                    :label="tag.tagName"
                  ></el-option>
                </template>

                <!-- 远程搜索的 -->
                <template v-if="autocompleteTags && autocompleteTags.length">
                  <el-option
                    v-for="tag in autocompleteTags"
                    :key="tag.id"
                    :value="tag.tagId"
                    :label="tag.tagName"
                  ></el-option>
                </template>
              </el-select>
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
      config: {},
      loading: false,
      autocompleteTags: [],
      autocompleteTagLoading: false
    };
  },
  mounted() {
    this.load();
  },
  methods: {
    async load() {
      try {
        this.config = await HttpClient.get("/api/admin/sys-config/all");
      } catch (err) {
        this.$notify.error({ title: "错误", message: err.message });
      }
    },
    async save() {
      this.loading = true;
      try {
        await HttpClient.post("/api/admin/sys-config/save", {
          config: JSON.stringify({
            siteTitle: this.config.siteTitle,
            siteDescription: this.config.siteDescription,
            siteKeywords: this.config.siteKeywords,
            siteNavs: this.config.siteNavs,
            recommendTags: this.config.recommendTags,
            bbsNavTags: this.config.bbsNavTagIds
          })
        });
        this.$message({ message: "提交成功", type: "success" });
        this.load();
      } catch (err) {
        this.$notify.error({ title: "错误", message: err.message });
      } finally {
        this.loading = false;
      }
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
    },
    async loadAutocompleteTags(query) {
      this.autocompleteTagLoading = true;
      this.autocompleteTags = [];
      try {
        const list = await HttpClient.get("/api/admin/tag/autocomplete", {
          keyword: query
        });
        
        if (list && list.length) {
          const me = this;
          this.autocompleteTags = list.filter(item => {
            if (!me.config.bbsNavTagIds || me.config.bbsNavTagIds.length === 0) {
              return true;
            }
            return me.config.bbsNavTagIds.indexOf(item.tagId) === -1;
          });
        }
      } catch (err) {
        console.log(err);
      } finally {
        this.autocompleteTagLoading = false;
      }
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
