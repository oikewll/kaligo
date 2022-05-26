<template>
    <div class="mod-darkswitch">
        <div v-if="!!sidebarStatus" class="list-switch" :class="{trigger: theme === 'dark'}">
            <span class="item" :class="{active: theme === 'light'}" @click="handleSwitch('light')"><i class="el-icon el-icon-sunny" />{{$lang('com.light')}}</span>
            <span class="item" :class="{active: theme === 'dark'}" @click="handleSwitch('dark')"><i class="el-icon el-icon-moon" />{{$lang('com.dark')}}</span>
        </div>
        <el-tooltip v-else class="mini-switch" :content="`${theme === 'light' ? 'Open' : 'Close'} dark mode`" placement="top">
            <el-switch
                v-model="theme"
                active-value="dark"
                inactive-value="light">
            </el-switch>
        </el-tooltip>
    </div>
</template>

<script>
export default {
    name: "Darkswitch",
    data: () =>{
        return {}
    },
    computed: {
        sidebarStatus(){
            return this.$store.getters["app/getSidebar"].opened;
        },
        theme: {
            get(){
                return this.$store.getters["app/getTheme"];
            },
            set(val){
                this.$store.commit('app/UPDATE_THEME', val);
            }
        }
    },
    methods: {
        handleSwitch(val){
            this.$store.commit('app/UPDATE_THEME', val);
        }
    }
}
</script>

<style lang="less">
.mod-darkswitch{
    position: absolute;
    bottom: 20px;
    font-size: 14px;
    left: 50%;
    transform: translateX(-50%);
    text-align: center;
    .list-switch{
        display: flex;
        background-color: #67727f;
        white-space: nowrap;
        height: 32px;
        border-radius: 500px;
        padding: 4px;
        position: relative;
        width: 150px;
        .item{
            transition: all 0.2s ease-in-out;
            display: inline-block;
            height: 24px;
            line-height: 24px;
            padding: 0 10px;
            border-radius: 100px;
            color: #2f4156;
            cursor: pointer;
            position: relative;
            z-index: 2;
            .el-icon{
                margin-right: 5px;
            }
            &.active{
                color: #bfcbd9;
            }
        }
        &.trigger{
            &::after{
                left: 72px;
            }
        }
        &::after{
            content: "";
            transition: left 0.3s ease;
            position: absolute;
            height: 24px;
            width: 72px;
            background-color: #2f4156;
            z-index: 1;
            left: 4px;
            top: 4px;
            border-radius: 100px;
        }
    }
}
</style>