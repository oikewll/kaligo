import dayjs from "dayjs";

export default ({ app, store }, inject) => {
    inject("dayjs", dayjs);
};
