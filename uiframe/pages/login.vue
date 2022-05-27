<template>
    <div class="page-login">
        <img src="@/assets/images/logo.svg" alt="kali-admin" class="logo" />
        <el-form
            :model="formData"
            :rules="rules"
            ref="logiForm"
            label-position="top"
        >
            <el-form-item prop="username">
                <el-input
                    v-model="formData.username"
                    @keyup.enter.native="login"
                    placeholder="输入账号"
                >
                </el-input>
            </el-form-item>
            <el-form-item prop="password">
                <el-input
                    v-model="formData.password"
                    type="password"
                    @keyup.enter.native="login"
                    placeholder="请输入密码"
                >
                </el-input>
            </el-form-item>
            <el-form-item class="btns">
                <el-button type="primary" @click="login" round
                    >立即登录</el-button
                >
            </el-form-item>
        </el-form>
        <footer class="copyright">©2022, Kaligo Tech.</footer>
    </div>
</template>

<script>
export default {
    name: "login",
    data: () => ({
        formData: {
            username: "",
            password: "",
            remember: true,
        },
        rules: {
            username: [
                {
                    required: true,
                    message: "请输入用户名",
                    trigger: "change",
                },
            ],
            password: [
                {
                    required: true,
                    message: "请输入密码",
                    trigger: "change",
                },
            ],
        },
    }),
    methods: {
        login() {
            // 模拟登录
            this.$refs.logiForm.validate(async (valid) => {
                if (valid) {
                    const rs = await this.$api.common.login(this.formData);
                    if (rs) {
                        this.$store.commit("user/setUserInfo", rs);
                        this.$message.success("登录成功");
                        this.$router.push("/");
                    }
                }
            });
        },
    },
};
</script>

<style lang="less" scoped>
.page-login {
    .logo {
        position: absolute;
        width: 60px;
        left: 50%;
        transform: translateX(-50%);
        top: 120px;
        opacity: 0.8;
    }
    & > * {
        position: relative;
        z-index: 1;
    }
    &:before,
    &:after {
        content: "";
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
    }
    &:before {
        z-index: 0;
        background: rgba(0, 0, 0, 0.5);
    }
    &:after {
        background-image: linear-gradient(
            0deg,
            rgba(44, 44, 44, 0.2),
            rgba(224, 23, 3, 0.6)
        );
        z-index: -1;
    }
    position: fixed;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    background: url("@/assets/images/login.jpg") no-repeat center center;
    background-size: cover;

    /deep/.el-form {
        width: 320px;
        margin: 220px auto 80px auto;

        &.el-form--label-top .el-form-item__label {
            padding-bottom: 0;
        }
        .el-input__inner {
            background: rgba(255, 255, 255, 0.1);
            color: #fff;
            border: 0 none;
            border-radius: 50px;
            height: 46px;
        }
        .el-form-item {
            margin-bottom: 30px;
        }
        .el-form-item__error {
            padding-left: 15px;
        }
    }
    .copyright {
        color: #fff;
        text-align: center;
    }
}
</style>
