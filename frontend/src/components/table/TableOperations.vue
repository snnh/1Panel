<template>
    <FuTableOperationActions
        v-if="cardRow"
        :buttons="buttons"
        :row="cardRow"
        :ellipsis="ellipsis"
        :extra="2"
        :trigger="resolvedTrigger"
        :dropdown-style="dropdownStyle"
    />
    <el-table-column
        v-else
        v-bind="$attrs"
        :label="label"
        :width="resolvedWidth"
        :min-width="resolvedMinWidth"
        :align="align"
        :fixed="resolvedFixed"
    >
        <template #default="{ row }">
            <FuTableOperationActions
                :buttons="buttons"
                :row="row"
                :ellipsis="columnEllipsis"
                :trigger="resolvedTrigger"
                :dropdown-style="dropdownStyle"
            />
        </template>
    </el-table-column>
</template>

<script setup lang="ts">
import { computed, type PropType } from 'vue';
import { useMediaQuery } from '@vueuse/core';

import FuTableOperationActions from './TableOperationActions.vue';
import type { FuTableOperationButton } from './shared';
import { useGlobalStore } from '@/composables/useGlobalStore';

defineOptions({ name: 'FuTableOperations' });

type DropdownTrigger = 'hover' | 'click' | 'contextmenu';
type DropdownTriggerValue = DropdownTrigger | DropdownTrigger[];

const normalizeWidth = (value?: string | number) => {
    if (value === undefined || value === null || value === '') {
        return undefined;
    }
    if (typeof value === 'number') {
        return value;
    }
    const trimmed = value.trim();
    if (!trimmed || trimmed === 'auto') {
        return undefined;
    }
    return trimmed;
};

const props = defineProps({
    buttons: {
        type: Array as PropType<FuTableOperationButton[]>,
        default: () => [],
    },
    label: {
        type: String,
        default: '',
    },
    width: {
        type: [Number, String],
        default: undefined,
    },
    minWidth: {
        type: [Number, String],
        default: undefined,
    },
    align: {
        type: String,
        default: 'center',
    },
    ellipsis: {
        type: Number,
        default: 2,
    },
    trigger: {
        type: [String, Array] as PropType<DropdownTriggerValue>,
        default: undefined,
    },
    fixed: {
        type: [Boolean, String],
        default: undefined,
    },
    fix: {
        type: [Boolean, String],
        default: undefined,
    },
    maxHeight: {
        type: [Number, String],
        default: undefined,
    },
    cardRow: {
        type: Object,
        default: undefined,
    },
});

const hasFinePointer = useMediaQuery('(hover: hover) and (pointer: fine)');
const resolvedTrigger = computed<DropdownTriggerValue>(
    () => props.trigger ?? (hasFinePointer.value ? 'hover' : 'click'),
);

// 手机屏幕宽度有限：固定右列会盖住其它列的表头与内容（且列宽常占视口一大半），
// 因此移动端一律取消固定；操作改由「更多」下拉承载，列宽按单个下拉按钮估算。
const { isMobile } = useGlobalStore();

const resolvedFixed = computed(() => {
    if (isMobile.value) {
        return false;
    }
    if (props.fixed !== undefined) {
        return props.fixed;
    }
    if (typeof props.fix === 'string') {
        return props.fix;
    }
    return props.fix ? 'right' : false;
});

const columnEllipsis = computed(() => (isMobile.value ? 0 : props.ellipsis));

const hasDynamicShow = computed(() => {
    return props.buttons.some((button) => typeof button.show === 'function');
});

const staticVisibleCount = computed(() => {
    return props.buttons.filter((button) => button.show !== false).length;
});

const estimatedVisibleCount = computed(() => {
    if (staticVisibleCount.value === 0) {
        return 0;
    }

    const renderCap = Math.max(props.ellipsis, 0) + 1;
    if (hasDynamicShow.value) {
        return renderCap;
    }

    return Math.min(staticVisibleCount.value, renderCap);
});

const estimatedWidth = computed(() => {
    const padding = 35;
    const buttonWidth = 58;
    if (isMobile.value) {
        // 移动端操作全部收进「更多」下拉，忽略页面写死的列宽/最小宽度
        return padding + buttonWidth;
    }
    const buttonsWidth = padding + estimatedVisibleCount.value * buttonWidth + buttonWidth;
    const minWidth = normalizeWidth(props.minWidth);
    if (typeof minWidth === 'number') {
        return Math.max(buttonsWidth, minWidth);
    }
    if (typeof minWidth === 'string') {
        const parsed = Number(minWidth.replace('px', ''));
        if (!Number.isNaN(parsed)) {
            return Math.max(buttonsWidth, parsed);
        }
    }
    return buttonsWidth;
});

const resolvedWidth = computed(() => {
    if (isMobile.value) {
        return estimatedWidth.value;
    }
    if (props.width === 'auto') {
        return undefined;
    }
    return normalizeWidth(props.width) ?? estimatedWidth.value;
});

const resolvedMinWidth = computed(() => {
    if (isMobile.value) {
        return undefined;
    }
    return props.width === 'auto' ? normalizeWidth(props.minWidth) : undefined;
});

const dropdownStyle = computed(() => {
    if (props.maxHeight === undefined || props.maxHeight === null || props.maxHeight === '') {
        return undefined;
    }
    const maxHeight = typeof props.maxHeight === 'number' ? `${props.maxHeight}px` : props.maxHeight;
    return {
        maxHeight,
        overflowY: 'auto',
    };
});
</script>
