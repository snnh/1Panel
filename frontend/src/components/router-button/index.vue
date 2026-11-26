<template>
    <div>
        <el-card class="router_card p-1 sm:p-0">
            <div class="flex w-full flex-col justify-start sm:items-center items-start sm:justify-between sm:flex-row">
                <el-radio-group :model-value="activeName" @update:model-value="updateActiveName" @change="handleChange">
                    <el-radio-button
                        class="router_card_button"
                        :class="{ 'opacity-50': button.disabled }"
                        :label="button.label"
                        :value="button.label"
                        :disabled="button.disabled"
                        v-for="(button, index) in buttonArray"
                        size="large"
                        :key="index"
                    >
                        <el-badge :value="button.count" v-if="button.count" is-dot>
                            <span>{{ button.label }}</span>
                        </el-badge>
                    </el-radio-button>
                </el-radio-group>
                <div class="flex flex-col gap-2 sm:flex-row">
                    <slot name="route-button"></slot>
                </div>
            </div>
        </el-card>
    </div>
</template>

<script lang="ts" setup>
import { openPage } from '@/utils/router';
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { hasPermissionMetaAccess, hasRouteAccess } from '@/utils/rbac';

defineOptions({ name: 'RouterButton' });

const props = defineProps({
    buttons: {
        type: Array<RouterButton>,
        required: true,
    },
    showExpiresAt: {
        type: Boolean,
        default: false,
    },
    activePath: {
        type: String,
        default: '',
    },
});

const router = useRouter();
const buttonArray = computed(() => {
    return props.buttons.filter((button) => {
        if (!hasPermissionMetaAccess(button.permission)) {
            return false;
        }
        if (button.path || button.name) {
            const route = router.resolve(button.path ? { path: button.path } : { name: button.name });
            return route.matched.length === 0 || hasRouteAccess(route);
        }
        return true;
    });
});

const activeName = ref('');

const updateActiveName = (label: string | number | boolean) => {
    if (!props.activePath) activeName.value = String(label);
};

const handleChange = async (label: string) => {
    const btn = buttonArray.value.find((btn) => btn.label === label);
    if (!btn || btn.disabled) {
        syncActiveName();
        return;
    }
    try {
        if (btn.path) await openPage({ path: btn.path });
        else if (btn.name) await openPage({ name: btn.name });
    } finally {
        syncActiveName();
    }
};

onMounted(() => {
    syncActiveName();
});

watch(
    () => [
        props.activePath || router.currentRoute.value.path,
        buttonArray.value.map((button) => button.label).join('|'),
    ],
    () => {
        syncActiveName();
    },
);

function syncActiveName() {
    if (!buttonArray.value.length) {
        activeName.value = '';
        return;
    }
    const route = props.activePath ? router.resolve(props.activePath) : router.currentRoute.value;
    const button = buttonArray.value.find((item) =>
        item.path ? route.path.startsWith(item.path) : item.name && route.name === item.name,
    );
    activeName.value = button?.label ?? buttonArray.value[0].label;
}
</script>

<style lang="scss" scoped>
.router_card {
    --el-card-padding: 0;

    :deep(.el-card__body) {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
}

.router_card_button {
    :deep(.el-radio-button__inner) {
        min-width: 100px;
        height: 100%;
        background-color: var(--panel-button-active) !important;
        box-shadow: none !important;
        outline: none !important;
        border: 2px solid transparent !important;
        color: var(--el-text-color-regular) !important;
    }

    :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
        color: var(--panel-button-text-color) !important;
        background-color: var(--panel-button-bg-color) !important;
        border-color: var(--panel-color-primary) !important;
        border-radius: 4px;
    }
}
</style>
