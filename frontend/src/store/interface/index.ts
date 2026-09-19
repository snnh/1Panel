import { RouteRecordRaw } from 'vue-router';
import { DeviceType } from '@/enums/app';
export interface ThemeConfigProp {
    panelName: string;
    primary: string;
    theme: string; // dark | bright ｜ auto
    footer: boolean;
    themeColor: string;
}

export interface GlobalState {
    language: string; // zh | en | tw
    device: DeviceType;
    themeConfig: ThemeConfigProp;
    // ui
    isFullScreen: boolean;
    openMenuTabs: boolean;
    menuAccordion: boolean;
    isLoading: boolean;
    loadingText: string;
    // auth
    ignoreCaptcha: boolean;
    agreeLicense: boolean;
    isLogin: boolean;
    entrance: string;
    csrfToken: string;
    // context
    hasNewVersion: boolean;
    lastFilePath: string;
    currentDB: string;
    currentPgDB: string;
    currentRedisDB: string;
    currentMongodbDB: string;
    showEntranceWarn: boolean;
    defaultNetwork: string;
    defaultIO: string;
    isOnRestart: boolean;
    // tags
    isAdmin: boolean;
    permissions: string[];
    masterOnlyPermissions: string[];
    isEnterprise: boolean;
    isIntl: boolean;
    docWithRegion: boolean;
    isFxplay: boolean;
    isOffline: boolean;
}

export interface MenuState {
    isCollapse: boolean;
    menuList: RouteRecordRaw[];
    withoutAnimation: boolean;
}

export interface TerminalState {
    showTerminalButton: boolean;
    lineHeight: number;
    letterSpacing: number;
    fontSize: number;
    fontFamily: string;
    backgroundColor: string;
    foregroundColor: string;
    cursorBlink: string;
    cursorStyle: string;
    scrollback: number;
    scrollSensitivity: number;
}
