<template>
    <el-popover
        v-model:visible="visible"
        trigger="click"
        placement="bottom-start"
        :width="400"
        popper-class="terminal-connection-popover"
        @before-enter="loadConnections"
    >
        <template #reference>
            <el-button
                class="terminal-connection-add"
                @click.stop
                @keydown.stop
                icon="Plus"
                text
                :aria-label="$t('terminal.createConn')"
            />
        </template>
        <div class="terminal-connection-menu">
            <div class="terminal-connection-actions">
                <el-button text class="terminal-connection-action" :disabled="connecting" @click="onNewSsh">
                    <el-icon><Plus /></el-icon>
                    {{ $t('terminal.createConn') }}
                </el-button>
                <el-button text class="terminal-connection-action" :disabled="connecting" @click="connectLocal()">
                    <el-icon><House /></el-icon>
                    {{ $t('terminal.localhost') }}
                </el-button>
            </div>
            <template v-if="connectionTree.length > 0 || loadingConnections">
                <el-input
                    v-model="connectionFilter"
                    size="small"
                    clearable
                    prefix-icon="Search"
                    :placeholder="$t('commons.button.search')"
                    :aria-label="$t('terminal.createConn')"
                />
                <el-tree
                    ref="treeRef"
                    v-loading="loadingConnections"
                    node-key="id"
                    default-expand-all
                    :expand-on-click-node="false"
                    :data="connectionTree"
                    :filter-node-method="filterConnection"
                    :empty-text="$t('commons.msg.noneData')"
                    class="terminal-connection-tree"
                >
                    <template #default="{ data }">
                        <span v-if="data.kind === 'group'" class="terminal-connection-group">
                            {{ data.label }}
                        </span>
                        <el-button
                            v-else
                            text
                            class="terminal-connection-item"
                            :disabled="connecting"
                            :title="data.label"
                            @click.stop="connectItem(data)"
                        >
                            <span class="terminal-connection-label">{{ data.label }}</span>
                        </el-button>
                    </template>
                </el-tree>
            </template>
        </div>
    </el-popover>
    <HostDialog
        ref="hostDialogRef"
        @on-conn-terminal="onHostCreated"
        @on-new-local="onLocalConfigured"
        @load-host-tree="loadHosts"
    />
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import { ElTree } from 'element-plus';
import i18n from '@/lang';
import { getHostTree, testByID, testLocalConn } from '@/api/modules/terminal';
import type { Host } from '@/api/interface/host';
import HostDialog from '@/components/terminal/host-create.vue';
import type { TerminalConnectionOptions } from './types';

const props = defineProps<{
    openSession: (options: TerminalConnectionOptions) => Promise<void>;
}>();
const visible = defineModel<boolean>({ default: false });
const hostDialogRef = ref<InstanceType<typeof HostDialog>>();
const onNewSsh = () => {
    if (connecting.value) return;
    visible.value = false;
    hostDialogRef.value?.acceptParams({ isLocal: false });
};
const onHostCreated = (title: string, wsID: number) => connect(wsID, title);
const connectLocal = () => connect(0, i18n.global.t('terminal.localhost'));
const onLocalConfigured = () => connect(0, i18n.global.t('terminal.localhost'));

interface ConnectionTreeItem {
    id: string;
    label: string;
    kind: 'group' | 'host';
    children?: ConnectionTreeItem[];
    wsID?: number;
}

const hostTree = ref<Array<Host.HostTree>>([]);
const treeRef = ref<InstanceType<typeof ElTree>>();
const connectionFilter = ref('');
const loadingConnections = ref(false);
const connecting = ref(false);
const connectionTree = computed<ConnectionTreeItem[]>(() => {
    return hostTree.value.map((group): ConnectionTreeItem => ({
        id: `host-group-${group.id}`,
        label: group.label === 'Default' ? i18n.global.t('commons.table.default') : group.label,
        kind: 'group',
        children: (group.children || []).map((host) => ({
            id: `host-${host.id}`,
            label: host.label,
            kind: 'host',
            wsID: host.id,
        })),
    }));
});
const loadHosts = async () => {
    hostTree.value = [];
    const res = await getHostTree({});
    hostTree.value = res.data || [];
};
const loadConnections = async () => {
    loadingConnections.value = true;
    try {
        await Promise.allSettled([loadHosts()]);
    } finally {
        loadingConnections.value = false;
    }
};
watch([connectionFilter, connectionTree], async () => {
    await nextTick();
    treeRef.value?.filter(connectionFilter.value);
});
const filterConnection = (value: string, data: ConnectionTreeItem) => {
    const filter = value.trim().toLowerCase();
    return !filter || data.label.toLowerCase().includes(filter);
};
const connectItem = (item: ConnectionTreeItem) => {
    if (item.kind === 'host' && item.wsID !== undefined) return connect(item.wsID, item.label);
};

const connect = async (wsID: number, title: string) => {
    if (connecting.value) return;
    connecting.value = true;
    visible.value = false;
    try {
        if (wsID === 0) {
            const res = await testLocalConn();
            if (!res.data) {
                hostDialogRef.value?.acceptParams({ isLocal: true });
                return;
            }
            await props.openSession({ title, wsID });
            return;
        }
        const res = await testByID(wsID);
        await props.openSession({
            title,
            wsID,
            error: res.data ? '' : 'Authentication failed. Please check the host information!',
        });
    } finally {
        connecting.value = false;
    }
};

defineExpose({ connectLocal });
</script>

<style lang="scss">
.terminal-connection-popover {
    max-width: calc(100vw - 24px);
    box-sizing: border-box;
}
</style>

<style scoped lang="scss">
.terminal-connection-add {
    width: 32px;
    height: 32px;
    margin: 0 4px;
    padding: 0;
    border-radius: 6px;
    color: var(--el-text-color-regular);

    &:hover {
        color: var(--el-color-primary);
    }
}

.terminal-connection-menu {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.terminal-connection-actions {
    display: flex;
    gap: 8px;
}

.terminal-connection-action {
    flex: 1;
    min-width: 0;
    height: 30px;
    margin: 0;
    padding: 0 6px;
    font-size: 13px;
    background-color: var(--el-fill-color-light);

    .el-icon {
        margin-right: 8px;
        color: var(--el-text-color-secondary);
    }
}

.terminal-connection-tree {
    max-height: min(192px, 35vh);
    min-height: 32px;
    overflow: auto;

    :deep(.el-tree-node__content) {
        height: 32px;
        border-radius: 4px;
    }
}

.terminal-connection-group {
    font-size: 12px;
    font-weight: 500;
    color: var(--el-text-color-secondary);
}

.terminal-connection-item {
    flex: 1;
    min-width: 0;
    height: 30px;
    margin: 0;
    padding: 0 8px 0 0;
    font-size: 13px;
    font-weight: normal;

    :deep(> span) {
        display: flex;
        align-items: center;
        gap: 12px;
        min-width: 0;
        width: 100%;
    }
}

.terminal-connection-label {
    flex: 1;
    text-align: left;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>
