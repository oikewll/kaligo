<template>
    <el-dropdown class="app-global-header-dropdown mr-md select-lang" @command="set_lang">
        <span class="el-dropdown-link">
            <img :src="require(`@/assets/images/flags/${lang}.svg`)" class="icon-flag" />{{maplang[$store.getters["i18n/getLang"] || 'en']}}
        </span>
        <el-dropdown-menu slot="dropdown" class="dropdown-list">
            <el-dropdown-item :command="idx" v-for="(item, idx) in maplang" :key="idx"><img :src="require(`@/assets/images/flags/${idx}.svg`)" class="icon-flag" />{{item}}</el-dropdown-item>
        </el-dropdown-menu>
    </el-dropdown>
</template>

<script>
export default {
    name: "LangChange",
    data: () =>{
        return {
            maplang: {
                'en': 'English',
                'zh-CN': '简体中文',
                'kh': 'ភាសាខ្មែរ',
            }
        }
    },
    computed: {
        lang(){
            return this.$store.getters["i18n/getLang"];
        }
    },
    methods: {
        set_lang(val) {
            this.$store.commit("i18n/set_lang", val);
        }
    }
}
</script>

<style lang="less">
.dropdown-list .el-dropdown-menu__item,
.el-dropdown-link{
    display: flex;
    align-items: center;
}
.icon-flag{
    width: 18px;
    margin-right: 5px;
}
</style>