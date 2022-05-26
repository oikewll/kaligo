<template>
    <div class="mod-table-box">
        <slot name="header" />
        <slot name="container" />
        <slot name="footer" />
    </div>
</template>

<script>
export default {
    name: "TableBox",
    props: {},
};
</script>

<style lang="less" scoped>
@boxshadow: 0 0 10px rgba(0, 0, 0, 0.15);
@colorBlue: rgb(67, 127, 195);

.mod-table-box {
    padding: 15px 20px;
    background-color: #fff;
    box-shadow: @boxshadow;
    display: flex;
    flex-direction: column;
}
</style>
