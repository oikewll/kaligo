<template>
    <div class="mod-lockscreen">
        <time class="countdown">{{ timeFormatter }}</time>
        <div class="wrap">
            <el-avatar
                :size="60"
                src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
            ></el-avatar>
            <p class="name">Admin</p>
            <el-form>
                <el-form-item prop="account">
                    <el-input
                        v-model="password"
                        placeholder="请输入密码解锁"
                        type="password"
                    >
                    </el-input>
                </el-form-item>
                <el-form-item class="btns">
                    <el-button type="primary" @click="login" round
                        >立即登录</el-button
                    >
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
export default {
    name: "lockscreen",
    data() {
        return {
            password: "",
            timeFormatter: "--:--:--",
        };
    },
    created() {
        setInterval(() => this.countdown(), 1000);
    },
    methods: {
        login() {
            this.$store.commit("app/UPDATE_LOCKSTATE");
        },
        countdown() {
            this.timeFormatter = this.$dayjs(new Date()).format("HH:mm:ss");
        },
    },
};
</script>

<style lang="less" scoped>
.mod-lockscreen {
    &:after {
        content: "";
        position: absolute;
        background-color: rgba(0, 0, 0, 0.35);
        left: 0;
        top: 0;
        z-index: -1;
        width: 100%;
        height: 100%;
    }
    background: url("@/assets/images/bg-lock.jpeg") no-repeat center center;
    background-size: cover;
    position: fixed;
    z-index: 1002;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    .countdown {
        font-size: 60px;
        line-height: 1;
        text-shadow: 0 0 8px rgba(0, 0, 0, 0.6);
        color: #fff;
        display: block;
        margin-bottom: 50px;
    }
    .wrap {
        display: flex;
        align-items: center;
        justify-items: center;
        flex-direction: column;
        padding: 50px 70px;
        background-color: rgba(255, 255, 255, 0.35);
        position: relative;
        border-radius: 5px;
        .el-avatar {
            position: absolute;
            left: 50%;
            top: 0;
            transform: translate3D(-50%, -50%, 0);
        }
        .name {
            font-size: 18px;
            color: #fff;
            margin-bottom: 20px;
            text-shadow: 0 0 8px rgba(0, 0, 0, 0.6);
        }
        .el-form {
            width: 250px;
        }
        /deep/ .el-input__inner {
            background: rgba(255, 255, 255, 0.8);
            color: #999;
            border: 0 none;
            border-radius: 50px;
            height: 46px;
        }
    }
}
</style>
