<template>
    <el-drawer
        :visible.sync="drawer"
        @opened="handleOpened"
        :direction="direction">
        <div slot="title" class="title" @click="handleMsgCenter">消息中心</div>
        <ul class="mod-listmsg">
            <li class="item unread">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item unread">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item unread">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item unread">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题这是消息标题</h5>
                <section class="txt">消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
            <li class="item">
                <h5 class="title">这是消息标题</h5>
                <section class="txt">消息文本推送</section>
            </li>
        </ul>
    </el-drawer>
</template>

<script>
export default {
    name: 'MessageDrawer',
    data(){
        return {
            direction: 'rtl',
        }
    },
    computed: {
        drawer: {
            get(){
                return this.$store.state["app"].msgCenter;
            },
            set(){
                this.$store.commit("app/UPDATE_MSGCENTER");
            }
        },

    },
    methods: {
        handleClose(){
            this.$store.commit("app/UPDATE_MSGCENTER");
        },
        handleOpened(){
            // 打开面板后做一些api交互，请求最新的消息列表
        },
        handleMsgCenter(){
            this.$router.push('/message');
            this.handleClose();
        }
    }
}
</script>

<style lang="less" scoped>
.el-drawer__header{
    margin-bottom: 15px;
    .title{
        cursor: pointer;
        &:hover{
            color: #2086f9;
        }
    }
}
.mod-listmsg{
    padding: 0 20px;
    .item{
        border-top: 1px dotted #e6e6e6;
        padding: 10px 0;
        cursor: pointer;
        &.unread{
            position: relative;
            &:before{
                content: "";
                position: absolute;
                width: 5px;
                height: 5px;
                border-radius: 50px;
                background-color:#2086f9;
                top: 18px;
                left: -10px;
            }
            .title{
                color: #2086f9;
            }
            .txt{
                color: #999;
            }
        }
        .title,
        .txt{
            overflow : hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
        }
        .title{
            font-weight: normal;
            font-size: 14px;
            color: #999;
        }
        .txt{
            color: #c5c5c5;
        }
    }
}
</style>