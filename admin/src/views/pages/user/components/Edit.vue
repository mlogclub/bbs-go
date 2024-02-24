<template>
  <a-modal
    v-model:visible="config.visible"
    :title="config.title"
    :size="appStore.table.size"
    @cancel="handleCancel"
    @before-ok="handleBeforeOk"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item label="类型" field="type">
        <a-select v-model="form.type" placeholder="用户类型">
          <a-option :value="0" label="用户" />
          <a-option :value="1" label="员工" />
        </a-select>
      </a-form-item>

      <a-form-item label="昵称" field="nickname">
        <a-input v-model="form.nickname" />
      </a-form-item>

      <a-form-item label="头像" field="avatar">
        <image-upload v-model="form.avatar" />
      </a-form-item>

      <a-form-item label="用户名" field="username">
        <a-input v-model="form.username" />
      </a-form-item>

      <a-form-item label="邮箱" field="email">
        <a-input v-model="form.email" />
      </a-form-item>

      <a-form-item label="性别" field="gender">
        <a-select v-model="form.gender">
          <a-option value="Male" label="男" />
          <a-option value="Female" label="女" />
        </a-select>
      </a-form-item>

      <!-- <a-form-item label="生日" field="birthday">
        <a-input v-model="form.birthday" />
      </a-form-item> -->

      <a-form-item label="主页" field="homePage">
        <a-input v-model="form.homePage" />
      </a-form-item>

      <a-form-item label="描述" field="description">
        <a-input v-model="form.description" />
      </a-form-item>

      <a-form-item label="角色" field="roles">
        <a-select v-model="form.roleIds" multiple placeholder="请选择角色">
          <a-option
            v-for="role in roles"
            :key="role.id"
            :value="role.id"
            :label="role.name"
          />
        </a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  import ImageUpload from '@/components/ImageUpload.vue';

  const emit = defineEmits(['ok']);

  const appStore = useAppStore();
  const formRef = ref();
  const config = reactive({
    visible: false,
    isCreate: false,
    title: '',
  });

  const form = ref({
    type: 0,
    username: undefined,
    email: undefined,
    nickname: undefined,
    avatar: undefined,
    gender: undefined,
    // birthday: undefined,
    homePage: undefined,
    description: undefined,
    roleIds: undefined,
  });

  const roles = ref();

  const rules = {
    type: [{ required: true, message: '请选择用户类型' }],
    nickname: [{ required: true, message: '请输入用户昵称' }],
  };

  const show = async () => {
    formRef.value.resetFields();

    config.isCreate = true;
    config.title = '新增';

    try {
      await loadRoles();
    } catch (e: any) {
      useHandleError(e);
    }

    config.visible = true;
  };

  const showEdit = async (id: any) => {
    formRef.value.resetFields();

    config.isCreate = false;
    config.title = '编辑';

    try {
      form.value = await axios.get(`/api/admin/user/${id}`);
      await loadRoles();
    } catch (e: any) {
      useHandleError(e);
    }

    config.visible = true;
  };

  const loadRoles = async () => {
    roles.value = await axios.get('/api/admin/role/roles');
  };

  const handleCancel = () => {
    formRef.value.resetFields();
  };
  const handleBeforeOk = async (done: (closed: boolean) => void) => {
    const validateErr = await formRef.value.validate();
    if (validateErr) {
      done(false);
      return;
    }
    try {
      const url = config.isCreate
        ? '/api/admin/user/create'
        : '/api/admin/user/update';
      await axios.postForm<any>(url, jsonToFormData(form.value));
      useNotificationSuccess('提交成功');
      emit('ok');
      done(true);
    } catch (e: any) {
      useHandleError(e);
      done(false);
    }
  };

  defineExpose({
    show,
    showEdit,
  });
</script>

<style lang="less" scoped></style>
