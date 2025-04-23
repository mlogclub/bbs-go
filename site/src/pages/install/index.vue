<template>
  <div class="install-container">
    <h1 class="title is-3 has-text-centered mb-5">BBS-GO 安装引导</h1>

    <!-- 步骤1: 欢迎页面 -->
    <div v-if="currentStep === 'welcome'" class="step-item">
      <div class="step-title">
        <div class="step-number">1</div>
        <h3 class="step-header">欢迎使用</h3>
      </div>
      <div class="content">
        <p>
          欢迎使用
          BBS-GO，这是一个基于Go语言开发的社区系统。在开始安装前，请确保您已经准备好以下内容：
        </p>
        <ol>
          <li>MySQL数据库（5.7+），并且已创建好空数据库</li>
          <li>已经规划好社区名称和描述</li>
          <li>已准备好管理员账号和密码</li>
        </ol>
        <p>接下来，我们将引导您完成安装过程。</p>
      </div>
      <div class="field is-grouped is-grouped-right">
        <div class="control">
          <button class="button is-primary" @click="goToStep('database')">
            下一步
          </button>
        </div>
      </div>
    </div>

    <!-- 步骤2: 数据库配置 -->
    <div v-if="currentStep === 'database'" class="step-item">
      <div class="step-title">
        <div class="step-number">2</div>
        <h3 class="step-header">数据库配置</h3>
      </div>
      <div>
        <div class="field">
          <label class="label"
            >数据库主机 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-host"
              placeholder="数据库主机地址（必填）"
              v-model="dbConfig.host"
              required
              @blur="validateDbForm"
            />
          </div>
          <p class="help">
            通常为localhost或127.0.0.1，如使用远程数据库请填写相应的主机地址
          </p>
        </div>
        <div class="field">
          <label class="label"
            >数据库端口 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-port"
              placeholder="数据库端口（必填）"
              v-model="dbConfig.port"
              required
              @blur="validateDbForm"
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >数据库名称 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-name"
              placeholder="请输入数据库名称（必填）"
              v-model="dbConfig.database"
              required
              @blur="validateDbForm"
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >数据库用户名 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-user"
              placeholder="请输入数据库用户名（必填）"
              v-model="dbConfig.username"
              required
              @blur="validateDbForm"
            />
          </div>
        </div>
        <div class="field">
          <label class="label">数据库密码</label>
          <div class="control">
            <input
              class="input"
              type="password"
              id="db-password"
              placeholder="请输入数据库密码"
              v-model="dbConfig.password"
            />
          </div>
        </div>
        <div v-if="dbError" class="notification is-warning is-light">
          {{ dbError }}
        </div>
        <div v-if="dbSuccess" class="notification is-success is-light">
          {{ dbSuccess }}
        </div>
        <div class="field is-grouped">
          <div class="control">
            <button
              class="button is-info"
              @click="testDbConnection"
              :disabled="testingConnection || !isDbFormValid"
            >
              {{ testingConnection ? "测试中" : "测试连接" }}
              <span v-if="testingConnection" class="loader"></span>
            </button>
          </div>
          <div class="control is-pulled-right ml-auto">
            <button
              class="button is-primary"
              @click="goToStep('site')"
              :disabled="testingConnection || !isDbFormValid"
            >
              {{ testingConnection ? "测试连接中" : "下一步" }}
              <span v-if="testingConnection" class="loader"></span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤3: 站点信息 -->
    <div v-if="currentStep === 'site'" class="step-item">
      <div class="step-title">
        <div class="step-number">3</div>
        <h3 class="step-header">站点信息</h3>
      </div>
      <div>
        <div class="field">
          <label class="label"
            >站点名称 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="site-name"
              placeholder="请输入站点名称（必填）"
              v-model="siteInfo.title"
              required
              @blur="validateSiteName"
            />
          </div>
        </div>
        <div class="field">
          <label class="label">站点描述</label>
          <div class="control">
            <textarea
              class="textarea"
              id="site-desc"
              rows="3"
              placeholder="请输入站点描述"
              v-model="siteInfo.description"
            ></textarea>
          </div>
        </div>
        <div v-if="siteError" class="notification is-warning is-light">
          {{ siteError }}
        </div>
        <div class="field is-grouped">
          <div class="control">
            <button class="button is-info" @click="goToStep('database')">
              上一步
            </button>
          </div>
          <div class="control is-pulled-right ml-auto">
            <button
              class="button is-primary"
              @click="goToStep('admin')"
              :disabled="!siteInfo.title"
            >
              下一步
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤4: 管理员设置 -->
    <div v-if="currentStep === 'admin'" class="step-item">
      <div class="step-title">
        <div class="step-number">4</div>
        <h3 class="step-header">管理员设置</h3>
      </div>
      <div>
        <div class="field">
          <label class="label"
            >管理员用户名 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="admin-username"
              placeholder="请输入管理员用户名（必填）"
              v-model="adminInfo.username"
              required
              @blur="validateAdminForm"
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >管理员密码 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="password"
              id="admin-password"
              placeholder="请输入管理员密码（必填）"
              v-model="adminInfo.password"
              required
              @blur="validateAdminForm"
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >确认密码 <span class="has-text-danger">*</span></label
          >
          <div class="control">
            <input
              class="input"
              type="password"
              id="admin-password-confirm"
              placeholder="请再次输入管理员密码（必填）"
              v-model="adminInfo.passwordConfirm"
              required
              @blur="validateAdminForm"
            />
          </div>
        </div>
        <div v-if="adminError" class="notification is-warning is-light">
          {{ adminError }}
        </div>
        <div class="field is-grouped">
          <div class="control">
            <button class="button is-info" @click="goToStep('site')">
              上一步
            </button>
          </div>
          <div class="control is-pulled-right ml-auto">
            <button
              class="button is-primary"
              @click="confirmInstall"
              :disabled="
                !adminInfo.username ||
                !adminInfo.password ||
                !adminInfo.passwordConfirm ||
                adminInfo.password !== adminInfo.passwordConfirm
              "
            >
              安装
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤5: 安装中 -->
    <div v-if="currentStep === 'install'" class="step-item">
      <div class="step-title">
        <div class="step-number">5</div>
        <h3 class="step-header">安装中</h3>
      </div>
      <progress
        class="progress"
        :class="installStatus.type === 'alert-danger' ? 'is-danger' : 'is-info'"
        :value="installProgress"
        max="100"
      ></progress>
      <div
        class="notification"
        :class="installStatus.type === 'alert-danger' ? 'is-danger' : 'is-info'"
      >
        {{ installStatus.message }}
      </div>
    </div>

    <!-- 步骤6: 安装完成 -->
    <div v-if="currentStep === 'complete'" class="step-item">
      <div class="step-title">
        <div class="step-number">6</div>
        <h3 class="step-header">安装完成</h3>
      </div>
      <div class="notification is-success">恭喜您，BBS-GO 已成功安装！</div>
      <p>现在您可以开始使用您的社区系统了。</p>
      <div class="has-text-centered mt-5">
        <a href="/" class="button is-success mr-2">进入首页</a>
        <a href="/admin" class="button is-primary">进入管理后台</a>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: "empty",
});

// 步骤控制
const currentStep = ref("welcome");

// 数据库配置
const dbConfig = ref({
  host: "localhost",
  port: "3306",
  database: "",
  username: "",
  password: "",
});
const dbError = ref("");
const testingConnection = ref(false);
const dbSuccess = ref("");

// 验证数据库表单是否填写完整
const isDbFormValid = computed(() => {
  return (
    dbConfig.value.host &&
    dbConfig.value.port &&
    dbConfig.value.database &&
    dbConfig.value.username
  );
});

// 验证数据库表单并显示提示
const validateDbForm = () => {
  if (!isDbFormValid.value) {
    dbError.value = "请填写完整的数据库信息";
    dbSuccess.value = "";
    return false;
  }
  return true;
};

// 监听数据库表单字段变化
watch(
  [
    () => dbConfig.value.host,
    () => dbConfig.value.port,
    () => dbConfig.value.database,
    () => dbConfig.value.username,
  ],
  () => {
    if (isDbFormValid.value) {
      dbError.value = "";
    }
  }
);

// 站点信息
const siteInfo = ref({
  title: "",
  description: "",
});
const siteError = ref("");

// 验证站点名称
const validateSiteName = () => {
  if (!siteInfo.value.title) {
    siteError.value = "请填写站点名称";
    return false;
  }
  siteError.value = "";
  return true;
};

// 监听站点名称变化
watch(
  () => siteInfo.value.title,
  (newVal) => {
    if (newVal) {
      siteError.value = "";
    }
  }
);

// 管理员信息
const adminInfo = ref({
  username: "",
  password: "",
  passwordConfirm: "",
});
const adminError = ref("");

// 验证管理员表单
const validateAdminForm = () => {
  if (!adminInfo.value.username) {
    adminError.value = "请填写管理员用户名";
    return false;
  }
  if (!adminInfo.value.password) {
    adminError.value = "请填写管理员密码";
    return false;
  }
  if (!adminInfo.value.passwordConfirm) {
    adminError.value = "请确认管理员密码";
    return false;
  }
  if (adminInfo.value.password !== adminInfo.value.passwordConfirm) {
    adminError.value = "两次输入的密码不一致";
    return false;
  }
  adminError.value = "";
  return true;
};

// 监听管理员表单字段变化
watch(
  [
    () => adminInfo.value.username,
    () => adminInfo.value.password,
    () => adminInfo.value.passwordConfirm,
  ],
  () => {
    if (
      adminInfo.value.username &&
      adminInfo.value.password &&
      adminInfo.value.passwordConfirm &&
      adminInfo.value.password === adminInfo.value.passwordConfirm
    ) {
      adminError.value = "";
    }
  }
);

// 安装状态
const installProgress = ref(0);
const installStatus = ref({
  message: "准备安装...",
  type: "is-info",
});

const { data: status } = await useAsyncData(() =>
  useHttpGet(`/api/install/status`)
);

if (status.value.installed) {
  const router = useRouter();
  router.push("/");
}

// 步骤导航
const goToStep = async (step) => {
  if (step === "site" && currentStep.value === "database") {
    // 验证数据库配置
    if (!validateDbForm()) {
      return;
    }

    // 测试数据库连接
    dbError.value = "";
    dbSuccess.value = "";
    testingConnection.value = true;

    try {
      await useHttpPost("/api/install/test_db_connection", {
        body: {
          host: dbConfig.value.host,
          port: dbConfig.value.port,
          database: dbConfig.value.database,
          username: dbConfig.value.username,
          password: dbConfig.value.password,
        },
      });

      // 连接成功才进入下一步
      dbSuccess.value = "数据库连接成功";
      testingConnection.value = false;
      currentStep.value = step;
    } catch (error) {
      dbError.value = error.message
        ? "数据库连接失败: " + error.message
        : "数据库连接失败";
      testingConnection.value = false;
      return; // 连接失败不进入下一步
    }

    return; // 已处理步骤切换，不执行后续代码
  }

  if (step === "admin" && currentStep.value === "site") {
    // 验证站点信息
    if (!validateSiteName()) {
      return;
    }
    siteError.value = "";
  }

  currentStep.value = step;
};

// 测试数据库连接
const testDbConnection = async () => {
  if (!validateDbForm()) {
    return;
  }

  dbError.value = "";
  dbSuccess.value = "";
  testingConnection.value = true;

  try {
    await useHttpPost("/api/install/test_db_connection", {
      body: {
        host: dbConfig.value.host,
        port: dbConfig.value.port,
        database: dbConfig.value.database,
        username: dbConfig.value.username,
        password: dbConfig.value.password,
      },
    });

    dbSuccess.value = "数据库连接成功";
  } catch (error) {
    dbError.value = error.message
      ? "数据库连接失败: " + error.message
      : "数据库连接失败";
  } finally {
    testingConnection.value = false;
  }
};

// 执行安装
const confirmInstall = () => {
  // 验证管理员信息
  if (!validateAdminForm()) {
    return;
  }

  // 切换到安装步骤
  currentStep.value = "install";

  // 准备安装数据
  const installData = {
    siteTitle: siteInfo.value.title,
    siteDescription: siteInfo.value.description,
    dbConfig: {
      host: dbConfig.value.host,
      port: dbConfig.value.port,
      database: dbConfig.value.database,
      username: dbConfig.value.username,
      password: dbConfig.value.password,
    },
    username: adminInfo.value.username,
    password: adminInfo.value.password,
  };

  // 初始化进度
  installProgress.value = 10;
  installStatus.value = {
    message: "正在连接数据库...",
    type: "is-info",
  };

  // 进度条动画
  const progressInterval = setInterval(() => {
    if (installProgress.value < 90) {
      installProgress.value += 5;
    }
  }, 500);

  // 执行安装请求
  fetch("/api/install/install", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(installData),
  })
    .then((response) => response.json())
    .then((data) => {
      clearInterval(progressInterval);

      if (data.success) {
        // 安装成功
        installProgress.value = 100;
        installStatus.value = {
          message: "安装完成！",
          type: "is-info",
        };

        // 显示完成页面
        setTimeout(() => {
          currentStep.value = "complete";
        }, 1000);
      } else {
        // 安装失败
        installProgress.value = 100;
        installStatus.value = {
          message: "安装失败: " + (data.message || "未知错误"),
          type: "is-danger",
        };
      }
    })
    .catch((error) => {
      clearInterval(progressInterval);

      // 显示错误
      installProgress.value = 100;
      installStatus.value = {
        message: "安装请求失败: " + error,
        type: "is-danger",
      };
    });
};
</script>

<style lang="scss" scoped>
body {
  background-color: #f5f5f5;
}
.install-container {
  max-width: 800px;
  margin: 50px auto;
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  padding: 30px;
}
.step-item {
  padding: 20px 0;
  /* border-bottom: 1px solid #eee; */
}
.step-title {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}
.step-number {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: #3273dc;
  color: #fff;
  text-align: center;
  line-height: 30px;
  margin-right: 10px;
  font-weight: bold;
}
.step-header {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
}
.ml-auto {
  margin-left: auto !important;
}
.loader {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  animation: spin 2s linear infinite;
  display: inline-block;
  margin-left: 10px;
}
.button {
  color: #fff;
}
@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
