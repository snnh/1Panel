import { defineStore } from 'pinia';
import type { StoreDefinition } from 'pinia';
import piniaPersistConfig from '@/config/pinia-persist';
import { GlobalState } from '../interface';
import { DeviceType } from '@/enums/app';
import i18n, { setActiveLocale } from '@/lang';
import { isMasterOnlyPermissionCode, setMasterOnlyPermissionCodes, toManageCode } from '@/utils/permission-codes';
import { clearPageStateCache } from '@/utils/page-state-cache';

const CN_DOCS_URL = 'https://1panel.cn/docs/v2';
const INTL_DOCS_URL = 'https://1panel.pro/docs/v2';

const GlobalStore = defineStore('GlobalState', {
    state: (): GlobalState => ({
        language: i18n.global.locale.value,
        device: DeviceType.Desktop,
        themeConfig: {
            panelName: '',
            primary: '#005eeb',
            theme: 'auto',
            footer: true,
            themeColor: '',
        },
        // ui
        isFullScreen: false,
        openMenuTabs: false,
        menuAccordion: false,
        isLoading: false,
        loadingText: '',
        csrfToken: '',
        // auth
        ignoreCaptcha: true,
        agreeLicense: false,
        isLogin: false,
        entrance: '',
        // context
        hasNewVersion: false,
        lastFilePath: '',
        currentDB: '',
        currentPgDB: '',
        currentRedisDB: '',
        currentMongodbDB: '',
        showEntranceWarn: true,
        defaultNetwork: 'all',
        defaultIO: 'all',
        isOnRestart: false,
        // tags
        isAdmin: false,
        permissions: [],
        masterOnlyPermissions: [],
        isEnterprise: false,
        isIntl: false,
        docWithRegion: true,
        isFxplay: false,
        isOffline: false,
    }),
    getters: {
        isDarkTheme: (state) =>
            state.themeConfig.theme === 'dark' ||
            (state.themeConfig.theme === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches),
        docsUrl: (state) => {
            if (state.docWithRegion) {
                return state.isIntl ? INTL_DOCS_URL : CN_DOCS_URL;
            }
            const lang = state.language.toLowerCase();
            const isChinese = lang === 'zh';
            return isChinese ? CN_DOCS_URL : INTL_DOCS_URL;
        },
        isMobile: (state) => state.device === DeviceType.Mobile,
    },
    actions: {
        setScreenFull() {
            this.isFullScreen = !this.isFullScreen;
        },
        setLogStatus(login: boolean) {
            this.isLogin = login;
        },
        setAuthInfo(payload: { isAdmin: boolean; permissions: string[]; masterOnlyPermissions?: string[] }) {
            this.isAdmin = !!payload.isAdmin;
            this.permissions = payload.permissions || [];
            this.masterOnlyPermissions = payload.masterOnlyPermissions || [];
            setMasterOnlyPermissionCodes(this.masterOnlyPermissions);
        },
        clearAuthInfo() {
            clearPageStateCache();
            this.permissions = [];
            this.masterOnlyPermissions = [];
            this.isAdmin = false;
            setMasterOnlyPermissionCodes([]);
        },
        hasPermission(permission: string) {
            setMasterOnlyPermissionCodes(this.masterOnlyPermissions);
            const normalizedPermission = permission.trim();
            if (!normalizedPermission) {
                return false;
            }
            if (isMasterOnlyPermissionCode(normalizedPermission)) {
                return false;
            }
            if (this.isAdmin) {
                return true;
            }
            if (this.permissions.includes(normalizedPermission)) {
                return true;
            }
            const managePermission = toManageCode(normalizedPermission);
            if (!managePermission) {
                return false;
            }
            if (isMasterOnlyPermissionCode(managePermission)) {
                return false;
            }
            return this.permissions.includes(managePermission);
        },
        async updateLanguage(language: string) {
            const activeLocale = await setActiveLocale(language);
            this.language = activeLocale;
            return activeLocale;
        },
        toggleDevice(value: DeviceType) {
            this.device = value;
        },
    },
    persist: piniaPersistConfig('GlobalState'),
}) as StoreDefinition<'GlobalState', GlobalState, any, any>;

export default GlobalStore;
