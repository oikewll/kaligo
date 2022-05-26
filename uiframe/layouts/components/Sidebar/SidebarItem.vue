<template>
    <el-submenu
        v-if="item.children && item.children.length"
        ref="subMenu"
        :index="`${item.id}`"
        popper-append-to-body
    >
        <template slot="title">
            <Item
                v-if="item.name"
                :icon="item.icon"
                :title="item.name"
            />
        </template>
        <sidebar-item
            v-for="(child, idx) in item.children"
            :key="idx"
            :is-nest="true"
            :item="child"
            :base-path="child.path"
            class="nest-menu"
        />
    </el-submenu>
    <el-menu-item 
        v-else
        :index="`${item.path}`"
        :class="{'submenu-title-noDropdown': !isNest }"
        @click="handleSetBar(item)"
    >
        <Item :icon="item.icon" :title="item.name" />
    </el-menu-item>
</template>

<script>
import path from "path";
import Item from "./Item";

export default {
    name: "SidebarItem",
    components: { Item },
    props: {
        // route object
        item: {
            type: Object,
            required: true,
        },
        isNest: {
            type: Boolean,
            default: false,
        },
        basePath: {
            type: String,
            default: "",
        },
    },
    data: () => ({ onlyOneChild: true }),
    methods: {
        handleSetBar(item = {}){
            this.$router.push(item.path);
            this.$store.commit('menu/MAIN_TABS_ADD', item);
        },
        hasOneShowingChild(children = [], parent) {
            const showingChildren = children.filter((item) => {
                if (item.hidden) {
                    return false;
                } else {
                    this.onlyOneChild = item;
                    return true;
                }
            });

            if (showingChildren.length === 1) {
                return true;
            }

            if (showingChildren.length === 0) {
                this.onlyOneChild = {
                    ...parent,
                    path: "",
                    noShowingChildren: true,
                };
                return true;
            }

            return false;
        },
        resolvePath(routePath) {
            return path.resolve(this.basePath, routePath);
        },
    },
};
</script>
