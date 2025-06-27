<template>
  <div class="install-container">
    <h1 class="title is-3 has-text-centered mb-5">
      {{ $t("pages.install.title") }}
    </h1>

    <!-- 步骤1: 欢迎页面 -->
    <div v-if="currentStep === 'welcome'" class="step-item">
      <div class="step-title">
        <div class="step-number">1</div>
        <h3 class="step-header">{{ $t("pages.install.step.welcome") }}</h3>
      </div>
      <div class="content">
        <p>{{ $t("pages.install.welcome.description") }}</p>
        <ol>
          <li>{{ $t("pages.install.welcome.requirements.mysql") }}</li>
          <li>{{ $t("pages.install.welcome.requirements.site") }}</li>
          <li>{{ $t("pages.install.welcome.requirements.admin") }}</li>
        </ol>
        <p>{{ $t("pages.install.welcome.guide") }}</p>
        <!-- 语言选择合并到欢迎页面 -->
        <div class="mt-5">
          <p class="mb-4">{{ $t("pages.install.language.description") }}</p>
          <div class="field">
            <div class="control">
              <label class="radio">
                <input type="radio" v-model="selectedLanguage" value="en-US" />
                {{ $t("pages.install.language.english") }}
              </label>
            </div>
            <div class="control">
              <label class="radio">
                <input type="radio" v-model="selectedLanguage" value="zh-CN" />
                {{ $t("pages.install.language.chinese") }}
              </label>
            </div>
          </div>
        </div>
      </div>
      <div class="field is-grouped is-grouped-right">
        <div class="control">
          <button class="button is-primary" @click="gotoStep('database')">
            {{ $t("pages.install.buttons.next") }}
          </button>
        </div>
      </div>
    </div>

    <!-- 步骤2: 数据库配置 -->
    <div v-if="currentStep === 'database'" class="step-item">
      <div class="step-title">
        <div class="step-number">3</div>
        <h3 class="step-header">{{ $t("pages.install.step.database") }}</h3>
      </div>
      <div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.database.host") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-host"
              :placeholder="$t('pages.install.database.hostPlaceholder')"
              v-model="dbConfig.host"
              required
            />
          </div>
          <p class="help">
            {{ $t("pages.install.database.hostHelp") }}
          </p>
        </div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.database.port") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-port"
              :placeholder="$t('pages.install.database.portPlaceholder')"
              v-model="dbConfig.port"
              required
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.database.name") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-name"
              :placeholder="$t('pages.install.database.namePlaceholder')"
              v-model="dbConfig.database"
              required
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.database.username") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="db-user"
              :placeholder="$t('pages.install.database.usernamePlaceholder')"
              v-model="dbConfig.username"
              required
            />
          </div>
        </div>
        <div class="field">
          <label class="label">{{
            $t("pages.install.database.password")
          }}</label>
          <div class="control">
            <input
              class="input"
              type="password"
              id="db-password"
              :placeholder="$t('pages.install.database.passwordPlaceholder')"
              v-model="dbConfig.password"
            />
          </div>
        </div>
        <div v-if="dbError" class="notification is-warning is-light">
          {{ dbError }}
        </div>
        <div v-if="dbSuccess" class="notification is-primary is-light">
          {{ dbSuccess }}
        </div>
        <div class="field is-grouped">
          <div class="control">
            <button class="button is-info" @click="gotoStep('welcome')">
              {{ $t("pages.install.buttons.previous") }}
            </button>
          </div>
          <div class="control">
            <button
              class="button is-info"
              @click="testDbConnection"
              :disabled="testingConnection || !isDbFormValid"
            >
              {{
                testingConnection
                  ? $t("pages.install.database.testing")
                  : $t("pages.install.database.testConnection")
              }}
              <span v-if="testingConnection" class="loader"></span>
            </button>
          </div>
          <div class="control is-pulled-right ml-auto">
            <button
              class="button is-primary"
              @click="gotoStep('site')"
              :disabled="testingConnection || !isDbFormValid"
            >
              {{
                testingConnection
                  ? $t("pages.install.database.testing")
                  : $t("pages.install.buttons.next")
              }}
              <span v-if="testingConnection" class="loader"></span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤4: 站点信息 -->
    <div v-if="currentStep === 'site'" class="step-item">
      <div class="step-title">
        <div class="step-number">4</div>
        <h3 class="step-header">{{ $t("pages.install.step.site") }}</h3>
      </div>
      <div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.site.title") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="site-name"
              :placeholder="$t('pages.install.site.titlePlaceholder')"
              v-model="siteInfo.title"
              required
              @blur="validateSiteName"
            />
          </div>
        </div>
        <div class="field">
          <label class="label">{{
            $t("pages.install.site.description")
          }}</label>
          <div class="control">
            <textarea
              class="textarea"
              id="site-desc"
              rows="3"
              :placeholder="$t('pages.install.site.descriptionPlaceholder')"
              v-model="siteInfo.description"
            ></textarea>
          </div>
        </div>
        <div v-if="siteError" class="notification is-warning is-light">
          {{ siteError }}
        </div>
        <div class="field is-grouped">
          <div class="control">
            <button class="button is-info" @click="gotoStep('database')">
              {{ $t("pages.install.buttons.previous") }}
            </button>
          </div>
          <div class="control is-pulled-right ml-auto">
            <button
              class="button is-primary"
              @click="gotoStep('admin')"
              :disabled="!siteInfo.title"
            >
              {{ $t("pages.install.buttons.next") }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤5: 管理员设置 -->
    <div v-if="currentStep === 'admin'" class="step-item">
      <div class="step-title">
        <div class="step-number">5</div>
        <h3 class="step-header">{{ $t("pages.install.step.admin") }}</h3>
      </div>
      <div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.admin.username") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="text"
              id="admin-username"
              :placeholder="$t('pages.install.admin.usernamePlaceholder')"
              v-model="adminInfo.username"
              required
              @blur="validateAdminForm"
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.admin.password") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="password"
              id="admin-password"
              :placeholder="$t('pages.install.admin.passwordPlaceholder')"
              v-model="adminInfo.password"
              required
              @blur="validateAdminForm"
            />
          </div>
        </div>
        <div class="field">
          <label class="label"
            >{{ $t("pages.install.admin.confirmPassword") }}
            <span class="has-text-danger">{{
              $t("pages.install.common.required")
            }}</span></label
          >
          <div class="control">
            <input
              class="input"
              type="password"
              id="admin-password-confirm"
              :placeholder="
                $t('pages.install.admin.confirmPasswordPlaceholder')
              "
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
            <button class="button is-info" @click="gotoStep('site')">
              {{ $t("pages.install.buttons.previous") }}
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
              {{ $t("pages.install.buttons.install") }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤6: 安装中 -->
    <div v-if="currentStep === 'install'" class="step-item">
      <div class="step-title">
        <div class="step-number">6</div>
        <h3 class="step-header">{{ $t("pages.install.step.install") }}</h3>
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

    <!-- 步骤7: 安装完成 -->
    <div v-if="currentStep === 'complete'" class="step-item">
      <div class="step-title">
        <div class="step-number">7</div>
        <h3 class="step-header">{{ $t("pages.install.step.complete") }}</h3>
      </div>
      <div class="notification is-success">
        {{ $t("pages.install.complete.congratulations") }}
      </div>
      <p>{{ $t("pages.install.complete.description") }}</p>
      <div class="has-text-centered mt-5">
        <a href="/" class="button is-success mr-2">{{
          $t("pages.install.complete.enterSite")
        }}</a>
        <a href="/admin" class="button is-primary">{{
          $t("pages.install.complete.enterAdmin")
        }}</a>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: "empty",
});

const { t, setLocale } = useI18n();

// 步骤控制
const currentStep = ref("welcome");

// 语言选择
const selectedLanguage = ref("en-US");

// 监听语言变化
watch(selectedLanguage, async (newLang) => {
  await setLocale(newLang);
});

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
    dbError.value = t("pages.install.database.validationError");
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
    siteError.value = t("pages.install.site.validationError");
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
    adminError.value = t("pages.install.admin.usernameError");
    return false;
  }
  if (!adminInfo.value.password) {
    adminError.value = t("pages.install.admin.passwordError");
    return false;
  }
  if (!adminInfo.value.passwordConfirm) {
    adminError.value = t("pages.install.admin.confirmPasswordError");
    return false;
  }
  if (adminInfo.value.password !== adminInfo.value.passwordConfirm) {
    adminError.value = t("pages.install.admin.passwordMismatchError");
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
  message: t("pages.install.install.preparing"),
  type: "is-info",
});

const { data: status } = await useMyFetch(`/api/install/status`);

if (status.value.installed) {
  const router = useRouter();
  router.push("/");
}

// 步骤导航
const gotoStep = async (step) => {
  if (step === "site" && currentStep.value === "database") {
    if (!validateDbForm()) {
      return;
    }
    dbError.value = "";
    dbSuccess.value = "";
    testingConnection.value = true;
    try {
      await useHttpPost("/api/install/test_db_connection", {
        host: dbConfig.value.host,
        port: dbConfig.value.port,
        database: dbConfig.value.database,
        username: dbConfig.value.username,
        password: dbConfig.value.password,
      });
      dbSuccess.value = t("pages.install.database.connectSuccess");
      testingConnection.value = false;
      currentStep.value = step;
    } catch (error) {
      dbError.value = error.message
        ? t("pages.install.database.connectFailed") + ": " + error.message
        : t("pages.install.database.connectFailed");
      testingConnection.value = false;
      return;
    }
    return;
  }
  if (step === "admin" && currentStep.value === "site") {
    if (!validateSiteName()) {
      return;
    }
    siteError.value = "";
  }
  currentStep.value = step;
};

const testDbConnection = async () => {
  if (!validateDbForm()) {
    return;
  }
  dbError.value = "";
  dbSuccess.value = "";
  testingConnection.value = true;
  try {
    await useHttpPost("/api/install/test_db_connection", {
      host: dbConfig.value.host,
      port: dbConfig.value.port,
      database: dbConfig.value.database,
      username: dbConfig.value.username,
      password: dbConfig.value.password,
    });
    dbSuccess.value = t("pages.install.database.connectSuccess");
  } catch (error) {
    dbError.value = error.message
      ? t("pages.install.database.connectFailed") + ": " + error.message
      : t("pages.install.database.connectFailed");
  } finally {
    testingConnection.value = false;
  }
};

const confirmInstall = () => {
  if (!validateAdminForm()) {
    return;
  }
  currentStep.value = "install";
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
    language: selectedLanguage.value,
  };
  installProgress.value = 10;
  installStatus.value = {
    message: t("pages.install.install.connecting"),
    type: "is-info",
  };
  const progressInterval = setInterval(() => {
    if (installProgress.value < 90) {
      installProgress.value += 5;
    }
  }, 500);
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
        installProgress.value = 100;
        installStatus.value = {
          message: t("pages.install.install.completed"),
          type: "is-info",
        };
        setTimeout(() => {
          currentStep.value = "complete";
        }, 1000);
      } else {
        installProgress.value = 100;
        installStatus.value = {
          message:
            t("pages.install.install.failed") +
            ": " +
            (data.message || t("pages.install.install.unknown")),
          type: "is-danger",
        };
      }
    })
    .catch((error) => {
      clearInterval(progressInterval);
      installProgress.value = 100;
      installStatus.value = {
        message: t("pages.install.install.requestFailed") + ": " + error,
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
.radio {
  padding: 8px 0;
}
</style>
