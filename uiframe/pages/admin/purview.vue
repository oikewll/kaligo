<template>
    <Wrap title="权限管理" :fillout="true">
        <template slot="content">
            <section class="item-section" v-for="(item, idx) in permissionList" :key="idx">
                <el-checkbox 
                    v-model="checkAll[idx]" 
                    @change="handleCheckAllChange($event, idx)">{{item.name}}</el-checkbox>
                <el-checkbox-group v-model="arrayCate" @change="handleCheckedChangeCate($event, idx)">
                    <el-checkbox v-for="(_item, _idx) in item.children" :label="_item" :key="_idx">
                        <div class="item-inner">
                            <span class="txt">{{_item.name}}</span>
                            <el-checkbox-group v-model="arraySubmit" @change="handleCheckedChangeSubmit($event, _idx)">
                                <el-checkbox v-for="(__item, __idx) in _item.children" :label="__item" :key="__idx">{{__item.name}}</el-checkbox>
                            </el-checkbox-group>
                        </div>
                    </el-checkbox>
                </el-checkbox-group>
            </section>
            <footer class="footer-section">
                <el-button type="primary" round size="medium" @click="handleSubmit">提交</el-button>
                <el-button type="warning" round size="medium" plain>取消</el-button>
            </footer>
        </template>
    </Wrap>
</template>

<script>
import Wrap from "@/components/Common/Wrap";
import { mapState } from "vuex";

export default {
    name: "index",
    components: { Wrap },
    data: () => {
        return {
            checkAll: [],
            arrayCate: [],
            arraySubmit: [],
            isIndeterminate: [],
            permissionList: [],
            datalist: [],
        }
    },
    computed: {
        menu_list() {
            return this.$store.getters["menu/getMenuSetting"];
        },
    },
    async created() {
        const {data} = await this.$api.user.permission();
        this.permissionList = data.data;
    },
    beforeDestroy() {
    },
    methods: {
        handleCheckAllChange(val = true, idx = 0) {
            console.log(val, idx)
            // 如果是true的话，当前array选中就是
            this[`checkedArray_${idx}`] = !!val ? this.permissionList[idx].children : [];
            console.log(this[`checkedArray_${idx}`]);
            // this.isIndeterminate[idx] = false;
        },
        handleCheckedChangeCate(value, idx) {
            let checkedCount = value.length;
            this.checkAll[idx] = checkedCount === this.permissionList[idx].children.length;
            console.log(idx, this.checkAll[idx])
        },
        handleCheckedChangeSubmit(value, idx) {
            console.log(this.arrayCate)
            console.log(this.arraySubmit)
            this.datalist = this.arraySubmit.map(item=>{
                return `${item.method}-${item.path}`
            }).join(',');
        },
        async handleSubmit(){
            const {data} = await this.$api.user.permission(this.datalist);
        }
    }
};
</script>

<style scoped lang="less">
.item-section{
    padding: 20px 0;
    display: flex;
}
.el-checkbox-group{
    padding-left: 15px;
    display: flex;
    flex-direction: column;
}
.item-inner{
    display: flex;
    .txt{
        margin: 0 10px 0 0;
    }
    .el-checkbox{
        margin-bottom: 10px;
    }
}
.footer-section{
    border-top: 1px solid #ececec;
    padding-top: 15px;
    text-align: center;
    // position: absolute;
    width: 100%;
    left: 0;
    // right: 15px;
    bottom: 0;
    background-color: #fff;
    // box-shadow: 0 0 15px rgba(0, 0, 0, 0.25);
    // height: 45px;
    // z-index: 2;
    // display: flex;
    // align-items: center;
    // justify-content: center;
}
</style>
