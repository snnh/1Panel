import { getSettingBaseInfo } from '@/api/modules/setting';
import { useTheme } from '@/global/use-theme';
import {
    searchXpackSetting,
    updateXpackSettingByKey as updateXpackSettingByKeyFromExtension,
} from '@/extensions/xpack';
import { GlobalStore } from '@/store';

let switchThemeFn: (() => void) | undefined;

const switchTheme = () => {
    if (!switchThemeFn) {
        switchThemeFn = useTheme().switchTheme;
    }
    switchThemeFn();
};

export function resetXSetting() {
    const globalStore = GlobalStore();
    globalStore.themeConfig.title = '';
    globalStore.themeConfig.logo = '';
    globalStore.themeConfig.logoWithText = '';
    globalStore.themeConfig.favicon = '';
    globalStore.masterAlias = '';
}

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

export async function getXpackSetting() {
    const res = await searchXpackSetting();
    if (!res) {
        initFavicon();
        resetXSetting();
        return;
    }
    initFavicon();
    return res;
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

export async function getXpackSettingForTheme() {
    const globalStore = GlobalStore();
    const res2 = await searchXpackSetting();
    if (res2) {
        globalStore.themeConfig.title = res2.data?.title;
        globalStore.themeConfig.logo = res2.data?.logo;
        globalStore.themeConfig.logoWithText = res2.data?.logoWithText;
        globalStore.themeConfig.favicon = res2.data?.favicon;
        globalStore.themeConfig.loginImage = res2.data?.loginImage;
        globalStore.themeConfig.loginBgType = res2.data?.loginBgType;
        globalStore.themeConfig.loginBackground = res2.data?.loginBackground;
        globalStore.themeConfig.loginBtnLinkColor = res2.data?.loginBtnLinkColor;
        globalStore.themeConfig.themeColor = res2.data?.themeColor;
        globalStore.masterAlias = res2.data.masterAlias;
        if (res2.data?.theme) {
            globalStore.themeConfig.theme = res2.data.theme;
        }
    } else {
        resetXSetting();
    }
    switchTheme();
    initFavicon();
}

export async function updateXpackSettingByKey(key: string, value: string) {
    return updateXpackSettingByKeyFromExtension(key, value);
}
