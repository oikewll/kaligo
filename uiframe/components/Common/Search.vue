<template>
    <div class="mod-search">
        <div class="divider">
            <el-divider></el-divider>
            <el-button
                round
                size="mini"
                icon="el-icon-search"
                @click="searchBox = !searchBox"
                >{{
                    !searchBox ? $lang("com.open") : $lang("com.close")
                }}</el-button
            >
        </div>
        <div class="wrap-top" v-if="!!searchBox">
            <el-input v-model="inputSearch" placeholder="请输入内容"></el-input>
            <el-input placeholder="请输入内容"></el-input>
            <el-input placeholder="请输入内容"></el-input>
            <el-button type="primary">搜索</el-button>
            <el-button type="warning" plain>导出</el-button>
        </div>
    </div>
</template>

<script>
export default {
    name: "Search",
    data: () => {
        return {
            searchBox: false,
            inputSearch: "",
        };
    },
};
</script>

<style lang="less">
.mod-search {
    .divider {
        display: flex;
        align-items: center;
    }
    margin-bottom: 15px;
    .el-divider--horizontal {
        margin: 0;
    }
}
</style>
