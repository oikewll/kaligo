<template>
    <section class="app-main">
        <nuxt :key="key" />
    </section>
</template>

<script>
export default {
    name: "AppMain",
    computed: {
        key() {
            return this.$route.path;
        },
    },
};
</script>

<style lang="scss" scoped>
.app-main {
    /* 50= navbar  50  */
    height: calc(100vh - 108px);
    padding: 15px;
    position: relative;
    overflow: hidden;
    box-sizing: content-box;
    overflow-y: auto;
    & > * {
        padding: 15px;
        background-color: #fff;
        display: flex;
        flex-direction: column;
        box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
        // height: 100%;
    }
}

.fixed-header + .app-main {
    padding-top: 50px;
}

.hasTagsView {
    .app-main {
        /* 84 = navbar + tags-view = 50 + 34 */
        min-height: calc(100vh - 84px);
    }

    .fixed-header + .app-main {
        padding-top: 84px;
    }
}
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
    .fixed-header {
        padding-right: 15px;
    }
}
</style>
