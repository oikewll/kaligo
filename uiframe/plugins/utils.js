import Utils from "~/utils";

export default function ({ app, store }, inject) {
    inject("utils", Utils);
}