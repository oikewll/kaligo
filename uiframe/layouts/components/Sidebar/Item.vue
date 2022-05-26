<script>
export default {
    name: "MenuItem",
    functional: true,
    props: {
        icon: {
            type: String,
            default: "",
        },
        title: {
            type: String,
            default: "",
        },
    },
    render(h, context) {
        const { icon, title } = context.props;
        const vnodes = [];

        if (icon) {
            // element-ui图标库
            if (icon.includes("el-icon")) {
                vnodes.push(<i class={[icon, "sub-el-icon"]} />);
            }
            // 增加fontawesome图标支持
            else if(icon.includes('fa')){
                vnodes.push(<i class={[icon, "sub-el-icon"]} />);
            }
            else {
                vnodes.push(<svg-icon icon-class={icon} />);
            }
        }

        if (title) {
            vnodes.push(<span slot="title">{title}</span>);
        }
        return vnodes;
    },
};
</script>

<style scoped>
.sub-el-icon {
    color: currentColor;
    width: 1em;
    height: 1em;
}
</style>
