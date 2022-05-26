<template>
    <div :class="['app-wrapper', ...classObj, `theme-${theme}`, { 'login-main': $route.path.match('/login') }]">
        <Sidebar v-if="!$route.path.match('/login')" />
        <div class="main-container">
            <Navbar v-if="!$route.path.match('/login')" />
            <Tabnav v-if="!$route.path.match('/login')" />
            <AppMain />
            <MessageDrawer />
        </div>
        <LockScreen v-if="lockScreen" />
    </div>
</template>

<script>
import Sidebar from "./components/Sidebar";
import AppMain from "./components/AppMain";
import Navbar from "./components/Navbar";
import Tabnav from "./components/Tabnav";
import MessageDrawer from "./components/MessageDrawer";
import LockScreen from "./components/LockScreen";
import { mapState } from "vuex";

export default {
    name: "Layout",
    components: {
        AppMain,
        Navbar,
        Sidebar,
        Tabnav,
        MessageDrawer,
        LockScreen,
    },
    computed: {
        ...mapState({
            sidebar: state => state.app.sidebar,
            lockScreen: state => state.app.lockScreen,
            theme: state => state.app.theme,
        }),
        classObj() {
            return {
                hideSidebar: !this.sidebar.opened,
                openSidebar: this.sidebar.opened,
                withoutAnimation: this.sidebar.withoutAnimation,
            };
        },
    },
    methods: {
        handleClickOutside() {
            this.$store.dispatch("app/closeSideBar", {
                withoutAnimation: false,
            });
        },
    },
    async created(){
        const {data} = await this.$api.common.init();
        this.$store.commit("menu/SET_MENUDATA", data.data);
    }
};
</script>

<style lang="less" scoped>
@import "../assets/less/mixin.less";
@import "../assets/less/variables.less";

.app-wrapper {
    position: relative;
    height: 100%;
    width: 100%;

    &.mobile.openSidebar {
        position: fixed;
        top: 0;
    }
}

.drawer-bg {
    background: #000;
    opacity: 0.3;
    width: 100%;
    top: 0;
    height: 100%;
    position: absolute;
    z-index: 999;
}

.fixed-header {
    position: fixed;
    top: 0;
    right: 0;
    z-index: 9;
    width: calc(100% - @sideBarWidth);
    transition: width 0.28s;
}

.hideSidebar .fixed-header {
    width: calc(100% - 54px);
}

.mobile .fixed-header {
    width: 100%;
}
</style>
