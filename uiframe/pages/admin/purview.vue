<template>
    <Wrap title="权限管理" :fillout="true">
        <template slot="content">
            <section class="item-section" 
                v-for="(item, idx) in permissionList" :key="idx">
                <el-checkbox 
                    disabled
                    @change="handleCheckAllChange($event, idx)">{{item.name}}</el-checkbox>
                <el-checkbox-group v-model="arrayCate" @change="handleCheckedChangeCate($event, idx)">
                    <section class="section-inner"
                        v-for="(_item, _idx) in item.children" :key="_idx">
                        <el-checkbox 
                            disabled
                            :label="_item">
                            <div class="item-inner">
                                <span class="txt">{{_item.name}}</span>
                                <el-checkbox-group v-model="arraySubmit" @change="handleCheckedChangeSubmit($event, _idx)">
                                    <el-checkbox v-for="(__item, __idx) in _item.children" :label="__item" :key="__idx">{{__item.name}}</el-checkbox>
                                </el-checkbox-group>
                            </div>
                        </el-checkbox>
                    </section>
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
        },
        async handleSubmit(){
            console.warn('submit data length: '+this.arraySubmit.length)
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
}
.el-checkbox-group{
    padding-left: 15px;
    display: flex;
    flex-direction: column;
    flex: 1;
}
.section-inner{
    margin-bottom: 15px;
    border-bottom: 1px solid #efefef;
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
