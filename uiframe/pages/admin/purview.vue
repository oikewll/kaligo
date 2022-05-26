<template>
    <Wrap title="权限管理">
        <template slot="content">
            <section class="item-section" v-for="(item, idx) in menu_list" :key="idx">
                <el-checkbox 
                    v-model="checkAll[idx]" 
                    @change="handleCheckAllChange($event, idx)">{{item.name}}</el-checkbox>
                <el-checkbox-group v-model="checkedArray_1">
                    <el-checkbox v-for="(_item, _idx) in item.children" :label="_item" :key="_idx">
                        <div class="item-inner">
                            <span class="txt">{{_item.name}}</span>
                            <el-checkbox-group v-model="checkedArray_2" @change="handleCheckedChange($event, _idx)">
                                <el-checkbox v-for="(__item, __idx) in _item.children" :label="__item" :key="__idx">{{__item.name}}</el-checkbox>
                            </el-checkbox-group>
                        </div>
                    </el-checkbox>
                </el-checkbox-group>
            </section>
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
            checkedArray_1: [],
            checkedArray_2: [],
            isIndeterminate: []
        }
    },
    computed: {
        menu_list() {
            return this.$store.getters["menu/getMenuSetting"];
        },
    },
    async created() {
    },
    beforeDestroy() {
    },
    methods: {
        handleCheckAllChange(val = true, idx = 0) {
            console.log(val, idx)
            // 如果是true的话，当前array选中就是
            this[`checkedArray_${idx}`] = !!val ? this.menu_list[idx].children : [];
            console.log(this[`checkedArray_${idx}`]);
            // this.isIndeterminate[idx] = false;
        },
        handleCheckedChange(value, idx) {
            return console.log(value.map(item => {
                return item.name
            }));
            let checkedCount = value.length;
            this.checkAll = checkedCount === this.cities.length;
            this.isIndeterminate = checkedCount > 0 && checkedCount < this.cities.length;
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
</style>
