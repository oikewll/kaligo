<template>
    <Wrap :title="getFormData.name" fillout>
        <template slot="content">
            <div class="mod-operation">
                <div class="search">
                    <el-input type="text" size="mini" placeholder="请输入关键词" />
                    <el-button
                        size="mini"
                        type="primary"
                        >搜索</el-button
                    >
                </div>
                <ul class="btn-list">
                    <li class="item">
                        <el-button
                            size="mini"
                            type="primary"
                            round
                            plain
                            >添加</el-button
                        >
                        <el-button
                            size="mini"
                            type="primary"
                            round
                            plain
                            >删除</el-button
                        >
                        <el-button
                            size="mini"
                            type="primary"
                            round
                            plain
                            >禁用</el-button
                        >
                        <el-button
                            size="mini"
                            type="primary"
                            round
                            plain
                            >激活</el-button
                        >
                        <el-button
                            size="mini"
                            type="primary"
                            round
                            plain
                            >刷新</el-button
                        >
                    </li>
                </ul>
            </div>
            <el-table
               
                :data="userList"
                style="width: 100%">
                <!-- <el-table-column
                    prop="id"
                    sortable
                    fixed
                    label="编号id">
                </el-table-column> -->
                <el-table-column
                    prop="username"
                    label="用户名">
                </el-table-column>
                <el-table-column
                    prop="realname"
                    label="真实姓名">
                </el-table-column>
                <el-table-column
                    prop="email"
                    label="邮箱">
                </el-table-column>
                <el-table-column
                    prop="created_at"
                    label="创建时间">
                </el-table-column>
                <el-table-column label="操作" align="right" width="180">
                    <template slot-scope="scope">
                        <el-button
                            size="mini"
                            type="primary"
                            round
                            plain
                            >编辑</el-button
                        >
                        <el-popconfirm
                            title="确定删除吗？"
                            @confirm="handleDelete(scope.$index, scope.row)"
                        >
                            <el-button
                                slot="reference"
                                size="mini"
                                type="danger"
                                round
                                plain
                                >删除</el-button
                            >
                        </el-popconfirm>
                    </template>
                </el-table-column>
            </el-table>
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
            userList: [{
                created_at: '2016-05-02',
                username: '王小虎',
                realname: '上海ddd',
                email: '@gmail',
            }],
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
                ]
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

        const {data: retUser} = await this.$api.user.getList();
        console.log(retUser)
        this.userList = retUser.data.data;
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
.mod-operation{
    width: 100%;
    display: flex;
    justify-content: space-between;
    padding-bottom: 20px;
    border-bottom: 1px solid #efefef;;
    .search{
        display: flex;
    }
}
</style>
