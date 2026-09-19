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

export async function initFavicon() {
    const globalStore = GlobalStore();
    document.title = globalStore.themeConfig.panelName;
    const favicon = globalStore.themeConfig.favicon;
    const customFaviconUrl = `/api/v2/images/favicon?t=${Date.now()}`;
    const fallbackFavicon = '/public/favicon.png';
    const setLink = (href: string) => {
        let link = document.querySelector("link[rel*='icon']") as HTMLLinkElement;
        if (!link) {
            link = document.createElement('link');
            link.rel = 'shortcut icon';
            link.type = 'image/x-icon';
            document.head.appendChild(link);
        }
        link.href = href;
    };

    if (favicon) {
        const testImg = new Image();
        testImg.onload = () => setLink(customFaviconUrl);
        testImg.onerror = () => setLink(fallbackFavicon);
        testImg.src = customFaviconUrl;
    } else {
        setLink(fallbackFavicon);
    }
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
