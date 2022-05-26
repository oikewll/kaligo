/**
 * 捕获vue的错误，不触发nuxt的错误补货逻辑（跳页面）
 */
import Vue from 'vue';
export default function({app}) {
    Vue.config.errorHandler = function(error,component,fun){
        console.error(error);
        return false;
    }
}