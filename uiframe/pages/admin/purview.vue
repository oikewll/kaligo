<template>
    <Wrap title="权限管理" :fillout="true">
        <section class="item-section" v-for="(item, idx) in permissionList" :key="idx">
            <el-checkbox 
                disabled
                v-model="checkAll[idx]" 
                @change="handleCheckAllChange($event, idx)">{{item.name}}</el-checkbox>
            <el-checkbox-group v-model="arrayCate" @change="handleCheckedChangeCate($event, idx)">
                <el-checkbox disabled v-for="(_item, _idx) in item.children" :label="_item" :key="_idx">
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
    </Wrap>
</template>

<script>
import Wrap from "@/components/Common/Wrap";
import { mapState } from "vuex";

export default {
    name: "admin",
    components: { Wrap },
    data: () => {
        return {
            checkAll: [],
            arrayCate: [],
            arraySubmit: [],
            isIndeterminate: [],
            permissionList: [],
            datalistSubmit: [],
        }
    },
    computed: {
        menu_list() {
            return this.$store.getters["menu/getMenuSetting"];
        },
    },
    async created() {
        this.inistList();
    },
    methods: {
        async inistList(){
            const {data} = await this.$api.user.permission();
            this.permissionList = data.data;

            let childrenList = [];
            data.data.forEach(item=>{
                childrenList = childrenList.concat(item.children)
            })
            let submitArrayList = [];
            childrenList.forEach(item=>{
                submitArrayList = submitArrayList.concat(item.children);
            })
            this.arraySubmit = submitArrayList.filter(item=>{
                return item.permit
            })
        },
        handleCheckAllChange(val = true, idx = 0) {
            console.log(val, idx)
            // 如果是true的话，当前array选中就是
            this.arrayCate = !!val ? this.permissionList[idx].children : [];
            this.arrayCate = Array.from(this.arrayCate);
            console.log(this.arrayCate);
            // this.isIndeterminate[idx] = false;
        },
        // 更改第二列
        handleCheckedChangeCate(value, idx) {
            let checkedCount = value.length;
            this.checkAll[idx] = checkedCount === this.permissionList[idx].children.length;
            this.checkAll = Array.from(this.checkAll);
        },
        handleCheckedChangeSubmit(value, idx) {
            console.log(this.arrayCate)
            console.log(this.arraySubmit)
        },
        async handleSubmit(){
            this.datalistSubmit = this.arraySubmit.map(item=>{
                return `${item.method}-${item.path}`
            }).join(',');
            const {data} = await this.$api.user.permission(this.datalistSubmit);
        }
    }
};
</script>

<style scoped lang="less">
.item-section{
    padding: 20px 0;
    display: flex;
    border-bottom: 1px solid #efefef;
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
