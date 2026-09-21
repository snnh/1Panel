<template>
    <div>
        <LayoutContent v-loading="loading" :title="$t('logs.login')">
            <template #search>
                <LogRouter current="LoginLog" />
            </template>
            <template #leftToolBar>
                <el-button v-permission type="primary" plain @click="onClean()">
                    {{ $t('logs.deleteLogs') }}
                </el-button>
            </template>
            <template #rightToolBar>
                <el-select v-model="searchStatus" @change="search()" clearable class="p-w-200">
                    <template #prefix>{{ $t('commons.table.status') }}</template>
                    <el-option :label="$t('commons.table.all')" value=""></el-option>
                    <el-option :label="$t('commons.status.success')" value="Success"></el-option>
                    <el-option :label="$t('commons.status.failed')" value="Failed"></el-option>
                </el-select>
                <TableSearch @search="search()" v-model:searchName="searchInfo" />
                <TableViewSwitch v-model="viewMode" storage-key="log-login" />
                <TableRefresh @search="search()" />
                <TableSetting title="login-log-refresh" @search="search()" />
            </template>
            <template #main>
                <ComplexTable
                    :pagination-config="paginationConfig"
                    :data="data"
                    @search="search"
                    :heightDiff="330"
                    v-model:view-mode="viewMode"
                >
                    <el-table-column :label="$t('logs.loginIP')" prop="ip" card-type="name" />
                    <el-table-column
                        v-if="isEnterprise"
                        :label="$t('commons.login.username')"
                        prop="user"
                        card-type="content"
                    />
                    <el-table-column :label="$t('logs.loginAddress')" prop="address" card-type="content" />
                    <el-table-column
                        :label="$t('logs.loginAgent')"
                        show-overflow-tooltip
                        prop="agent"
                        card-type="content-full"
                    />
                    <el-table-column :label="$t('logs.loginStatus')" prop="status" card-type="status">
                        <template #default="{ row }">
                            <Status :status="row.status" :msg="loadMsg(row.message)" />
                        </template>
                    </el-table-column>
                    <el-table-column
                        prop="createdAt"
                        :label="$t('commons.table.date')"
                        :formatter="dateFormat"
                        show-overflow-tooltip
                        card-type="content"
                    />
                </ComplexTable>
            </template>
        </LayoutContent>
        <ConfirmDialog ref="confirmDialogRef" @confirm="onSubmitClean"></ConfirmDialog>
    </div>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/confirm-dialog/index.vue';
import LogRouter from '@/views/log/router/index.vue';
import { dateFormat } from '@/utils/date';
import { cleanLogs, getLoginLogs } from '@/api/modules/log';
import { onMounted, reactive, ref } from 'vue';
import i18n from '@/lang';
import { MsgSuccess } from '@/utils/message';
import { useGlobalStore } from '@/composables/useGlobalStore';

const { isEnterprise } = useGlobalStore();

const loading = ref();
const data = ref();
const viewMode = ref<'table' | 'card'>('table');
const confirmDialogRef = ref();
const paginationConfig = reactive({
    cacheSizeKey: 'login-log-page-size',
    currentPage: 1,
    pageSize: Number(localStorage.getItem('login-log-page-size')) || 20,
    total: 0,
});
const searchInfo = ref<string>('');
const searchStatus = ref<string>('');

const search = async () => {
    let params = {
        info: searchInfo.value,
        status: searchStatus.value,
        page: paginationConfig.currentPage,
        pageSize: paginationConfig.pageSize,
    };
    loading.value = true;
    await getLoginLogs(params)
        .then((res) => {
            loading.value = false;
            data.value = res.data.items;
            paginationConfig.total = res.data.total;
        })
        .catch(() => {
            loading.value = false;
        });
};

const onClean = async () => {
    let params = {
        header: i18n.global.t('logs.deleteLogs'),
        operationInfo: i18n.global.t('commons.msg.delete'),
        submitInputInfo: i18n.global.t('logs.deleteLogs'),
    };
    confirmDialogRef.value!.acceptParams(params);
};

const loadMsg = (msg: string) => {
    if (msg === 'ErrAuth') {
        return i18n.global.t('commons.login.errorAuthInfo');
    }
    if (msg === 'ErrMFA') {
        return i18n.global.t('commons.login.errorMfaInfo');
    }
    return msg;
};

const onSubmitClean = async () => {
    await cleanLogs({ logType: 'login' });
    search();
    MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
};

onMounted(() => {
    search();
});
</script>
