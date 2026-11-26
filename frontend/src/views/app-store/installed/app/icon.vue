<template>
    <div class="icon">
        <el-avatar @click="$emit('openDetail')" shape="square" :size="77" :src="getAppIconUrl(appKey, undefined, icon)">
            <SvgIcon iconName="p-appstore" />
        </el-avatar>
    </div>
</template>

<script lang="ts" setup>
import { getAppIconUrl } from '@/api/modules/app';

interface Props {
    appKey: string;
    icon?: string;
}
defineProps<Props>();

defineEmits(['openDetail']);
</script>
