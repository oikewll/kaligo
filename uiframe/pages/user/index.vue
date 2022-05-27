<template>
    <Wrap title="用户列表" fillout>
        <div class="mod-operation">
            <div class="search">
                <el-input type="text" size="mini" placeholder="账号/昵称" />
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
                        @click="handleAddDialog"
                        round>添加</el-button
                    >
                    <el-button
                        size="mini"
                        type="danger"
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
                        @click="initList()"
                        >刷新</el-button
                    >
                </li>
            </ul>
        </div>
        <el-table 
            :data="dataUser.data"
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
                        @click="handleEdit(scope.row.id)"
                        >编辑</el-button
                    >
                    <el-popconfirm
                        :title="`确定删除${scope.row.id}吗？`"
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
        <div class="pagination" ref="count-page">
            <div class="btn-list"></div>
            <el-pagination
                background
                layout="prev, pager, next"
                :total="dataUser.total"
                :page-size="pageSize"
                @current-change="handlePageChange($event, 'filter-domain')"
            >
            </el-pagination>
        </div>
        <el-dialog
            class="mod-dialog"
            :title="getFormData.name"
            :visible.sync="dialogVisible"
        >
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
                    <template v-if="item.type === 'input'">
                        <el-input
                            v-model="formData[item.field]"
                            :type="item.type"
                            @keyup.enter.native="onSubmitForm"
                            :placeholder="item.placeholder"
                        ></el-input>
                    </template>
                    <template v-else-if="item.type === 'checkbox'">
                        <el-checkbox-group v-model="formData[item.field]">
                            <el-checkbox 
                                v-for="(option, _idx) in item.options" 
                                :key="_idx"
                                :label="option.value" 
                                :disabled="option.disabled"
                            >
                            </el-checkbox>
                        </el-checkbox-group>
                    </template>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <div class="btn-list">
                    <el-button type="primary" :loading="formLoading" @click="onSubmitForm">{{editMod ? '提交修改' : '立即创建'}}</el-button>
                    <el-button @click="dialogVisible = false">取消</el-button>
                </div>
            </div>
        </el-dialog>
    </Wrap>
</template>

<script>
import Wrap from "@/components/Common/Wrap";

export default {
    name: "user",
    components: { Wrap },
    data: () => {
        return {
            editMod: false,
            pageSize: 19,
            dialogVisible: false,
            dataUser: {},
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
    created() {
        this.initList();
    },
    methods: {
        async initList(){
            const {data: retUser} = await this.$api.user.getList({
                size: this.pageSize
            });
            console.log(retUser)
            this.dataUser = retUser.data;
        },
        async handleAddDialog(){
            // 获取动态表单
            const { data } = await this.$api.user.getForm();
            this.getFormData = data.data;

            let formRules = {}; // 校验规则
            let formData = {};  // 提交表单
            data.data.components.forEach((item) => {
                formRules[item.field] = [
                    Object.assign({}, item.validate),
                ];
                formData[item.field] = item.value;
            });
            this.formRules = formRules;
            this.formData = formData;

            this.dialogVisible = true;
        },
        // 编辑状态是，先渲染动态列表的表单，在根据表单的内容拼接表单数据，两步
        async handleEdit(id){
            this.editMod = true;
            const {data: retForm} = await this.$api.user.getForm(id);

            const {data: retUser} = await this.$api.user.get(id);
            retUser.data.groups = retUser.data.groups.split(',');

            this.getFormData = retForm.data;
            let formRules = {}; // 校验规则
            retForm.data.components.forEach((item) => {
                formRules[item.field] = [
                    Object.assign({}, item.validate),
                ];
            });
            this.formRules = formRules;
            this.formData = retUser.data;
            
            this.dialogVisible = true;
        },
        async handleDelete(idx, row){
            const {id} = row;
            const {data} = await this.$api.user.delete(id);
            if(data.code === 0){
                this.$message.success('delete successfully');
                this.initList();
            }
        },
        // 新增和编辑都在这里，动态获取方法和path
        async onSubmitForm() {
            this.formLoading = true;
            this.$refs.formSubmit.validate(async (valid) => {
                // #TODO 验证规则还不完善，要去掉valid
                if (!valid) {
                    const { getFormData, formData, editMod } = this;

                    let ret = null;
                    // formData.groups = formData.groups.join(",");
                    if(editMod){
                        ret = await this.$api.user.put(formData.id, formData);
                    } else{
                        ret = await this.$axios[getFormData["method"].toLocaleLowerCase()](getFormData["path"], formData);
                    }
                    if (ret) {
                        this.$message.success("Successfully");
                        // this.$router.push("/");
                    }
                }
                this.formLoading = false;
                this.dialogVisible = false;
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
        .el-button{
            margin-left: -1px;
        }
    }
}
.pagination {
    white-space: nowrap;
    display: flex;
    justify-content: space-between;
    padding-top: 10px;
    margin-bottom: -5px;
}
</style>
