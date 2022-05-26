<template>
    <section class="mod-datepicker">
        <div class="wrap-left">
            <el-date-picker
                v-model="valuedate"
                type="daterange"
                align="right"
                unlink-panels
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                :picker-options="pickerOptions"
            >
            </el-date-picker>
            <el-button
                type="primary"
                size="large"
                round
                icon="el-icon-search"
                @click="emmitEvent('queryAnalyticsData', 'filter', valuedate)"
                >{{ $lang("com.search") }}</el-button
            >
            <el-button
                @click="exportSelect('filter')"
                type="primary"
                size="medium"
                plain
                round
                icon="el-icon-printer"
                :loading="loadingExportAnalysis"
                >{{ $lang("com.export") }}</el-button
            >
        </div>
        <div class="wrap-right">
            <el-button
                @click="exportTodayForm"
                type="warning"
                size="medium"
                plain
                round
                icon="el-icon-printer"
                :loading="loadingExportToday"
            >
                导出今日新增到桌面
            </el-button>
        </div>
    </section>
</template>

<script>
export default {
    name: "Datepicker",
    date: () => {
        return {
            valuedate: "",
            pickerOptions: {
                shortcuts: [
                    {
                        text: "最近一周",
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(
                                start.getTime() - 3600 * 1000 * 24 * 7
                            );
                            picker.$emit("pick", [start, end]);
                        },
                    },
                    {
                        text: "最近一个月",
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(
                                start.getTime() - 3600 * 1000 * 24 * 30
                            );
                            picker.$emit("pick", [start, end]);
                        },
                    },
                    {
                        text: "最近三个月",
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(
                                start.getTime() - 3600 * 1000 * 24 * 90
                            );
                            picker.$emit("pick", [start, end]);
                        },
                    },
                ],
            },
        };
    },
    props: {
        "handle-emitter": Function,
        loadingExportAnalysis: {
            default: false,
            type: Boolean,
        },
        loadingExportToday: {
            default: false,
            type: Boolean,
        },
    },
    methods: {
        emmitEvent(methods = "", data) {
            this.$emit("emmit-event", methods, data);
        },
    },
};
</script>

<style lang="less" scoped>
.mod-datepicker {
    box-shadow: @boxshadow;
    margin-bottom: 20px;
    background-color: #fff;
    padding: 15px;
    display: flex;
    justify-content: space-between;
}
</style>
