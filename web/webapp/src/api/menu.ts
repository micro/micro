const Menu = [
    {
        title: "homePage",
        group: "apps",
        icon: "home",
        name: "home"
    },
    {
        title: "cliTerminal",
        group: "apps",
        icon: "tune",
        name: "cli"
    },
    {
        title: "registryInfo",
        component: "apps",
        icon: "cloud"
    },
    {
        title: "callService",
        component: "apps",
        icon: "train"
    },
    {divider: true}
];
// reorder menu
Menu.forEach((item: any) => {
    if (item.items) {
        item.items.sort((x: any, y: any) => {
            let textA = x.title.toUpperCase();
            let textB = y.title.toUpperCase();
            return textA < textB ? -1 : textA > textB ? 1 : 0;
        });
    }
});

export default Menu;
