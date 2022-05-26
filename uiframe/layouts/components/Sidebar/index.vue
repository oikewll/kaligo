<template>
    <div class="has-logo sidebar-container">
        <logo :collapse="isCollapse" />
        <el-scrollbar wrap-class="scrollbar-wrapper">
            <el-menu
                :default-active="activeMenu"
                :collapse="isCollapse"
                background-color="transparent"
                text-color="#bfcbd9"
                :unique-opened="false"
                active-text-color="#409eff"
                :collapse-transition="false"
                mode="vertical"
            >
                <template v-for="(route, idx) in menu_sidebar">
                    <SidebarItem
                        v-if="!route.show"
                        :key="idx"
                        :item="route"
                        :base-path="route.path"
                    />
                </template>
            </el-menu>
            <SwitchTheme />
        </el-scrollbar>
    </div>
</template>

<script>
import { mapState } from "vuex";
import Logo from "./Logo";
import SidebarItem from "./SidebarItem";
import SwitchTheme from "./Switch";

export default {
    name: "Sidebar",
    data: ()=>{
        return {}
    },
    components: {
        SidebarItem, 
        Logo,
        SwitchTheme
    },
    computed: {
        ...mapState({
            sidebar: (state) => state.app.sidebar,
            menu_sidebar: (state) => state.menu.menu_sidebar,
        }),
        activeMenu() {
            const route = this.$route;
            const { meta, path } = route;
            // if set path, the sidebar will highlight the path you set
            if (meta.activeMenu) {
                return meta.activeMenu;
            }
            return path;
        },
        isCollapse() {
            return !this.sidebar.opened;
        },
    },
};
</script>

<style lang="less" scoped>
.scrollbar-wrapper{
    position: relative;
}
</style>