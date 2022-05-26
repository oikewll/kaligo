<template>
    <el-tabs
        v-model="activeName"
        type="card"
        closable
        @tab-click="handleClick"
        @tab-remove="removeTab"
    >
        <el-tab-pane
            v-for="item in tabs"
            :key="item.name"
            :label="item.name"
            :name="item.name"
        >
        </el-tab-pane>
    </el-tabs>
</template>

<script>
export default {
    name: "Tabnav",
    data: () => {
        return {};
    },
    computed: {
        menu_sidebar() {
            return this.$store.getters["menu/getMenu"];
        },
        tabs: {
            get() {
                return this.$store.state["menu"].main_tabs;
            },
            set(value) {
                this.$store.commit("menu/MAIN_TABS_UPDATE", value);
            },
        },
        activeName: {
            get() {
                return this.$store.state["menu"].main_tabs_active_name;
            },
            set(value) {
                let name = value;
                // tab组件初始化会传入数字
                if (this.tabs[value]) {
                    name = this.tabs[value].name;
                }
                this.$store.commit("menu/SET_MAIN_TABS_ACTIVE_NAME", name);
            },
        },
    },
    watch: {
        $route() {
            this.matchedNav();
        },
    },
    created() {
        this.matchedNav();
    },
    methods: {
        matchedNav() {
            const { path, name } = this.$route;
            let tab = this.findTabByPath(path);
            if (!tab) {
                // if(['/login'].find(item => item === path)) return;
                let nav = this.getNavByPath(path);
                if (!nav || Object.keys(nav).length === 0) return;
                this.$store.commit("menu/MAIN_TABS_ADD", nav);
                this.activeName = nav.name;
            } else {
                this.activeName = tab.name;
            }
        },
        findTabByPath(path) {
            return this.tabs.find((item) => item.path === path);
        },
        getNavByPath(path) {
            try {
                let nav = [];
                // 递归函数
                const Recursion = (source) => {
                    source.forEach((el) => {
                        nav.push(el);
                        el.children && el.children.length > 0
                            ? Recursion(el.children)
                            : ""; // 子级递归
                    });
                };
                Recursion(this.menu_sidebar);
                return nav.find((item) => item.path === path);
            } catch (err) {
                console.error(err);
                return {};
            }
        },
        handleClick(panel) {
            const { path } = this.tabs[panel.index];
            this.$router.push(path);
        },
        removeTab(targetName) {
            let tabs = this.tabs;
            let activeName = this.activeName;
            let activeTab;
            if (activeName === targetName) {
                tabs.forEach((tab, index) => {
                    if (tab.name === targetName) {
                        let nextTab = tabs[index + 1] || tabs[index - 1];
                        if (nextTab) {
                            activeName = nextTab.name;
                            activeTab = nextTab;
                        }
                    }
                });
            }
            if (activeTab) {
                this.$router.push(activeTab.path || "/");
            }
            this.activeName = activeName;
            this.tabs = tabs.filter((tab) => tab.name !== targetName);
        },
    },
};
</script>

<style lang="less" scoped>
/deep/ .el-tabs__nav-next,
/deep/ .el-tabs__nav-prev {
    line-height: 28px;
}
/deep/ .el-tabs__item {
    height: 28px;
    line-height: 28px;
    font-size: 12px;
}
.el-tabs--card {
    background-color: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    /deep/ .el-tabs__header {
        .el-tabs__item {
            padding-left: 10px;
            padding-right: 10px;
            color: #666;
            &.is-active {
                color: #2086f9;
                padding-left: 10px !important;
                padding-right: 10px !important;
            }
            &.is-closable:hover {
                padding-left: 10px;
                padding-right: 10px;
                color: #2086f9;
            }
            .el-icon-close {
                width: 14px;
            }
            &:last-child,
            &:nth-child(2) {
                padding-left: 10px !important;
                padding-right: 10px !important;
            }
        }
        margin-bottom: 0;
        overflow: hidden;
        border: 0 none;
    }
}
</style>
