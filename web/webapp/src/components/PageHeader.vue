<template>
    <v-layout row class="align-center layout px-4 pt-4 app--page-header">
        <div class="page-header-left">
            <h3 class="pr-3">{{title}}</h3>
        </div>
        <v-breadcrumbs :items="breadcrumbs" divider=">">
        </v-breadcrumbs>
        <v-spacer></v-spacer>

    </v-layout>
</template>

<script lang="ts">
    import menu from "@/api/menu";

    export default {
        data() {
            return {
                title: ""
            };
        },
        computed: {
            breadcrumbs: function () {
                let breadcrumbs = [
                    {
                        disabled: false,
                        icon: "home",
                        href: "/"
                    }
                ];
                menu.forEach(item => {
                    if (item.items) {
                        let child = item.items.find(i => {
                            return i.component === this.$route.name;
                        });

                        if (child) {
                            let p = {
                                text: item.title,
                                disabled: false,
                                href: item.path
                            };

                            let c = {
                                text: child.title,
                                disabled: false,
                                href: child.path
                            };

                            breadcrumbs.push(p);
                            breadcrumbs.push(c);
                        }
                    } else {
                        if (item.name === this.$route.name) {
                            let p = {
                                text: item.title,
                                disabled: false,
                                href: item.path
                            };

                            // this.title = item.title;
                            breadcrumbs.push(p);
                        }
                    }
                });
                return breadcrumbs;
            }
        }
    };
</script>
