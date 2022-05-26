<template>
    <div class="navbar">
        <div class="left-menu">
            <hamburger
                id="hamburger-container"
                :is-active="sidebar.opened"
                class="hamburger-container"
                @toggleClick="handleSideBar"
            />
            <el-menu
                :default-active="activeMenu"
                class="el-menu-demo"
                mode="horizontal"
                @select="handleSelect"
            >
                <el-menu-item v-for="(item, idx) in menuSetting" :index="`${idx}`" :key="idx">
                    {{item.name}}
                </el-menu-item>
                <el-menu-item index="messages" disabled>{{
                    $lang("com.messages")
                }}</el-menu-item>
            </el-menu>
        </div>
        <div class="right-menu">
            <LangChange />
            <el-dropdown
                class="avatar-container right-menu-item hover-effect"
                trigger="click"
            >
                <div class="avatar-wrapper">
                    <el-badge :value="200" :max="99" class="item">
                        <el-avatar
                            :size="30"
                            src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
                            >{{ userInfo.username[0] }}</el-avatar
                        >
                    </el-badge>
                    <span class="name">{{ userInfo.username }}</span>
                    <i class="el-icon-caret-bottom" />
                </div>
                <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item divided @click.native="handleMsgCenter">
                        <el-badge :value="200" :max="99">
                            <span class="txt">{{ $lang("com.messages") }}</span>
                        </el-badge>
                    </el-dropdown-item>
                    <el-dropdown-item divided @click.native="logout">
                        <span class="txt">{{ $lang("com.password") }}</span>
                    </el-dropdown-item>
                    <el-dropdown-item divided @click.native="lockScreen">
                        <span class="txt"
                            >{{ $lang("com.lock-screen") }}(ALT+L)</span
                        >
                    </el-dropdown-item>
                    <el-dropdown-item divided @click.native="logout">
                        <span class="txt">{{ $lang("com.logout") }}</span>
                    </el-dropdown-item>
                </el-dropdown-menu>
            </el-dropdown>
        </div>
    </div>
</template>

<script>
import { mapState } from "vuex";
import hamburger from "./Hamburger";
import LangChange from "./LangChange";

export default {
    data: () => {
        return {
            activeMenu: "collection",
        };
    },
    name: "Navbar",
    computed: {
        ...mapState({
            sidebar: (state) => state.app.sidebar,
            userInfo: (state) => state.user.userInfo,
            menuSetting: (state) => state.menu.menu_setting,
        }),
    },
    components: {
        hamburger,
        LangChange,
    },
    methods: {
        handleSelect(key) {
            this.$store.commit("menu/SET_MAIN_SIDE", key);
        },
        handleSideBar() {
            this.$store.commit("app/UPDATE_SIDEBAR");
        },
        async logout() {
            await this.$store.dispatch("user/logout");
            this.$router.push(`/login?redirect=${this.$route.fullPath}`);
        },
        handleMsgCenter() {
            this.$store.commit("app/UPDATE_MSGCENTER");
        },
        lockScreen() {
            this.$store.commit("app/UPDATE_LOCKSTATE");
        },
    },
};
</script>

<style lang="less" scoped>
.navbar {
    // display: flex;
    // justify-content: space-between;
    // width: 100%;
    height: 50px;
    overflow: hidden;
    position: relative;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    .left-menu {
        float: left;
        display: flex;
        align-items: center;
        height: 100%;
        .el-dropdown-link {
            margin-left: 20px;
            cursor: pointer;
            &:hover {
                color: #409eff;
            }
        }
    }
    .hamburger-container {
        cursor: pointer;
        transition: background 0.3s;
        -webkit-tap-highlight-color: transparent;

        &:hover {
            background: rgba(0, 0, 0, 0.025);
        }
    }

    .breadcrumb-container {
        float: left;
    }

    .errLog-container {
        display: inline-block;
        vertical-align: top;
    }

    .right-menu {
        float: right;
        height: 100%;
        display: flex;
        align-items: center;

        &:focus {
            outline: none;
        }
        .el-dropdown {
            cursor: pointer;
            &:hover {
                color: #2086f9;
            }
        }
        .right-menu-item {
            display: inline-block;
            padding: 0 8px;
            height: 100%;
            font-size: 18px;
            color: #5a5e66;
            vertical-align: text-bottom;

            .name {
                margin: 0 0 0 10px;
                color: #333;
                font-size: 14px;
            }

            &.hover-effect {
                cursor: pointer;
                transition: background 0.3s;

                &:hover {
                    background: rgba(0, 0, 0, 0.025);
                }
            }
        }
        .avatar-container {
            margin-right: 30px;
            display: flex;
            align-items: center;
            display: none;
            .avatar-wrapper {
                display: flex;
                position: relative;
                align-items: center;
                > .item {
                    margin-top: 5px;
                }
                .el-icon-caret-bottom {
                    cursor: pointer;
                    position: absolute;
                    right: -20px;
                    top: 15px;
                    font-size: 12px;
                }
            }
        }
    }
}
.el-dropdown-menu__item--divided {
    margin-top: 0;
}
</style>
