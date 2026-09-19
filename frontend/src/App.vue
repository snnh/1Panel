<template>
    <el-config-provider :locale="i18nLocale" :button="config" size="default">
        <router-view v-if="isRouterAlive" />
    </el-config-provider>
</template>

<script setup lang="ts">
import { reactive, computed, ref, nextTick, provide } from 'vue';
import { useGlobalStore } from '@/composables/useGlobalStore';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import en from 'element-plus/es/locale/lang/en';
import { useTheme } from '@/global/use-theme';
useTheme();

const { language } = useGlobalStore();
const config = reactive({
    autoInsertSpace: false,
});

const i18nLocale = computed(() => {
    if (language.value === 'zh') return zhCn;
    return en;
});

const isRouterAlive = ref(true);

const reload = () => {
    isRouterAlive.value = false;
    nextTick(() => {
        isRouterAlive.value = true;
    });
};
provide('reload', reload);
</script>

<style scoped lang="scss"></style>
