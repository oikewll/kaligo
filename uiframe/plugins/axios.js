import axios from "axios";
import qs from "qs";
import Utils from "@/utils";
const { getTimezone } = Utils;

export default function ({ app, store, route, redirect, req }, inject) {
    const $http = axios.create({
        baseURL: process.server
            ? app.$config.server.api
            : app.$config.client.api,
        headers: {
            "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
        },
        params: {
            os: "web",
        },
        timeout: 20000,
        transformRequest: [
            function (data, headers) {
                return qs.stringify(data);
            },
        ],
    });

    $http.interceptors.request.use(
        function (config) {
            // Do something before request is sent
            config.headers["Accept-Language"] =
                store.getters["i18n/getLang"] === "en"
                    ? "en,en;q=0.9"
                    : "zh-CN,zh;q=0.9";
            // config.headers["Authorization"] = `Bearer ${store.state.token || '<NO-ACCESS-TOKEN>'}`;
            // config.headers["Accept-Timezone"] = getTimezone();
            return config;
        },
        function (error) {
            // Do something with request error
            return Promise.reject(error);
        }
    );

    $http.interceptors.response.use(
        function (response) {
            // process.server && console.log(`[${response.status}] ${response.request.path}`,response.data);
            return response;
        },
        function (err) {
            // Do something with response error
            console.error(
                `[${err.response && err.response.status}] ${
                    err.response && err.response.request.path
                }`
            );
            // console.log(err.response && err.response.data);
            return Promise.reject(err?.response?.data);
        }
    );

    inject("axios", $http);
}
