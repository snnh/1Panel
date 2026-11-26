import { RouteRecordRaw } from 'vue-router';
import { DeviceType } from '@/enums/app';
import type { PanelCacheSettings } from '@/utils/panel-cache';
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
    panelCache: PanelCacheSettings;
    cacheUserID: string;
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
    isFxplay: boolean;
    currentNode: string;
    docWithRegion: boolean;
    isOffline: boolean;
}

export interface MenuState {
    isCollapse: boolean;
    menuList: RouteRecordRaw[];
    isMenuLoaded: boolean;
    menuRequestId: number;
    systemVersion: string;
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
