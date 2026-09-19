import { getSettingBaseInfo } from '@/api/modules/setting';
import { useTheme } from '@/global/use-theme';
import { GlobalStore } from '@/store';

let switchThemeFn: (() => void) | undefined;

const switchTheme = () => {
    if (!switchThemeFn) {
        switchThemeFn = useTheme().switchTheme;
    }
    switchThemeFn();
};

export function initFavicon() {
    const globalStore = GlobalStore();
    document.title = globalStore.themeConfig.panelName;
    const link = document.createElement('link');
    link.rel = 'shortcut icon';
    link.type = 'image/x-icon';
    link.href = '/public/favicon.png';
    document.head.appendChild(link);
}

export async function loadBaseDataFromDB() {
    const globalStore = GlobalStore();
    const res = await getSettingBaseInfo();
    document.title = res.data.panelName;
    globalStore.entrance = res.data.securityEntrance;
    globalStore.openMenuTabs = res.data.menuTabs === 'Enable';
    globalStore.menuAccordion = res.data.menuAccordion === 'Enable';
    if (res.data.theme) {
        globalStore.themeConfig.theme = res.data.theme;
    }
    switchTheme();
    initFavicon();
}

export function applyLocalThemeSettings() {
    switchTheme();
    initFavicon();
}
