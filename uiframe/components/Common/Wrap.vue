<template>
    <div class="mod-wrap" :class="{'is-fillout': fillout}">
        <h6 class="title" v-if="!!title">{{title}}</h6>
        <slot name="content" />
    </div>
</template>

<script>
export default {
    props: {
        title: {
            default: '',
            type: String,
        },
        fillout: {
            default: false,
            type: Boolean,
        }
    }
}
</script>

<style lang="less" scoped>
.mod-wrap {
    margin-bottom: 20px;
    padding: 30px;
    background-color: #fff;
    display: flex;
    flex-direction: column;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    &.is-fillout{
        height: 100%;
        overflow: hidden;
        overflow-y: auto;
    }
    .title{
        margin-bottom: 20px;
        font-size: 16px;
        font-weight: normal;
        color: #333;
        position: relative;
        padding-left: 10px;
        line-height: 1;
        &::before{
            content: "";
            width: 5px;
            height: 100%;
            background-color: #2086f9;
            position: absolute;
            left: 0;
            top: 0;
        }
    }
}
</style>