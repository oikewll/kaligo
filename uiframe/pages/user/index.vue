<template>
    <Wrap :title="getFormData.name" fillout>
        <template slot="content">
            <el-form
                ref="formSubmit"
                :model="formData"
                :rules="formRules"
                label-width="80px"
            >
                <el-form-item
                    :label="item.title"
                    v-for="(item, idx) in getFormData.components"
                    :key="idx"
                    :prop="item.field"
                >
                    <el-input
                        v-model="formData[item.field]"
                        :type="item.type"
                        @keyup.enter.native="onSubmitForm"
                        :placeholder="item.placeholder"
                    ></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button
                        type="primary"
                        @click="onSubmitForm"
                        :loading="formLoading"
                        >立即创建</el-button
                    >
                    <el-button @click="history.back()">取消</el-button>
                </el-form-item>
            </el-form>
        </template>
    </Wrap>
</template>

<script>
import Wrap from "@/components/Common/Wrap";

export default {
    name: "user",
    components: { Wrap },
    data: () => {
        return {
            getFormData: {
                name: "",
            },
            formData: {},
            formRules: {
                username: [
                    {
                        required: true,
                        message: "请输入用户名",
                        trigger: "change",
                    },
                ],
                password: [
                    {
                        required: true,
                        message: "请输入密码",
                        trigger: "change",
                    },
                ],
            },
            formLoading: false,
        };
    },
    computed: {},
    async created() {
        const { data } = await this.$api.user.getForm();
        this.getFormData = data.data;

        let formRules = {};
        data.data.components.forEach((item) => {
            formRules[item.field] = [
                Object.assign({ trigger: "change" }, item.validate),
            ];
        });
        this.formRules = formRules;
    },
    methods: {
        async onSubmitForm() {
            this.formLoading = true;
            this.$refs.formSubmit.validate(async (valid) => {
                if (valid) {
                    const { getFormData, formData } = this;
                    const rs = await this.$axios[
                        getFormData["method"].toLocaleLowerCase()
                    ](getFormData["path"], formData);
                    if (rs) {
                        this.$message.success("Successfully");
                        // this.$router.push("/");
                    }
                }
                this.formLoading = false;
            });
        },
    },
};
</script>

<style scoped lang="less">
</style>
