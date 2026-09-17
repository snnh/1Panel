<template>
    <DrawerPro
        v-model="drawerVisible"
        :header="title"
        size="large"
        :resource="dialogData.title === 'add' ? '' : dialogData.rowData?.name"
        @close="handleClose"
    >
        <el-form ref="formRef" label-position="top" :model="dialogData.rowData" :rules="rules" v-loading="loading">
            <el-row type="flex" justify="center">
                <el-col :span="22">
                    <el-form-item :label="$t('commons.table.name')" prop="name">
                        <el-input
                            :disabled="dialogData.title === 'edit'"
                            clearable
                            v-model.trim="dialogData.rowData!.name"
                        />
                    </el-form-item>
                    <el-form-item :label="$t('toolbox.clam.scanDir')" prop="path">
                        <el-input v-model="dialogData.rowData!.path">
                            <template #prepend>
                                <el-button
                                    icon="Folder"
                                    v-permission
                                    v-node-admin
                                    @click="scanDirRef.acceptParams({ dir: true })"
                                />
                            </template>
                        </el-input>
                    </el-form-item>

                    <el-form-item :label="$t('toolbox.clam.infectedStrategy')" prop="infectedStrategy">
                        <el-radio-group v-model="dialogData.rowData!.infectedStrategy">
                            <el-radio value="none">{{ $t('toolbox.clam.none') }}</el-radio>
                            <el-radio value="remove">{{ $t('commons.button.delete') }}</el-radio>
                            <el-radio value="move">{{ $t('toolbox.clam.move') }}</el-radio>
                            <el-radio value="copy">{{ $t('commons.button.copy') }}</el-radio>
                        </el-radio-group>
                        <span class="input-help">
                            {{ $t('toolbox.clam.' + dialogData.rowData!.infectedStrategy + 'Helper') }}
                        </span>
                    </el-form-item>
                    <el-form-item v-if="hasInfectedDir()" :label="$t('toolbox.clam.infectedDir')" prop="infectedDir">
                        <el-input v-model="dialogData.rowData!.infectedDir">
                            <template #prepend>
                                <el-button
                                    icon="Folder"
                                    v-permission
                                    v-node-admin
                                    @click="infectedDirRef.acceptParams({ dir: true })"
                                />
                            </template>
                        </el-input>
                    </el-form-item>

                    <el-form-item :label="$t('cronjob.timeout')" prop="timeoutItem">
                        <el-input type="number" class="selectClass" v-model.number="dialogData.rowData!.timeoutItem">
                            <template #append>
                                <el-select v-model="dialogData.rowData!.timeoutUnit" style="width: 80px">
                                    <el-option :label="$t('commons.units.second')" value="s" />
                                    <el-option :label="$t('commons.units.minute')" value="m" />
                                    <el-option :label="$t('commons.units.hour')" value="h" />
                                </el-select>
                            </template>
                        </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('commons.table.description')" prop="description">
                        <el-input type="textarea" :rows="3" clearable v-model="dialogData.rowData!.description" />
                    </el-form-item>
                </el-col>
            </el-row>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="drawerVisible = false">{{ $t('commons.button.cancel') }}</el-button>
                <el-button v-permission v-node-admin :disabled="loading" type="primary" @click="onSubmit(formRef)">
                    {{ $t('commons.button.confirm') }}
                </el-button>
            </span>
        </template>
    </DrawerPro>
    <FileList ref="scanDirRef" @choose="loadDir" />
    <FileList ref="infectedDirRef" @choose="loadInfectedDir" />
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Rules } from '@/global/form-rules';
import FileList from '@/components/file-list/index.vue';
import i18n from '@/lang';
import { ElForm } from 'element-plus';
import { MsgSuccess } from '@/utils/message';
import { Toolbox } from '@/api/interface/toolbox';
import { createClam, updateClam } from '@/api/modules/toolbox';
import { splitTimeFromSecond, transferTimeToSecond } from '@/utils/validate';

const scanDirRef = ref();
const infectedDirRef = ref();

interface DialogProps {
    title: string;
    rowData?: Toolbox.ClamInfo;
    getTableList?: () => Promise<any>;
}
const loading = ref();
const title = ref<string>('');
const drawerVisible = ref(false);
const dialogData = ref<DialogProps>({
    title: '',
});

const acceptParams = (params: DialogProps): void => {
    dialogData.value = params;
    if (dialogData.value.rowData!.timeout) {
        let item = splitTimeFromSecond(dialogData.value.rowData!.timeout);
        dialogData.value.rowData.timeoutItem = item.timeItem;
        dialogData.value.rowData.timeoutUnit = item.timeUnit;
    }
    title.value = i18n.global.t('commons.button.' + dialogData.value.title);
    drawerVisible.value = true;
};
const emit = defineEmits<{ (e: 'search'): void }>();

const handleClose = () => {
    drawerVisible.value = false;
};

const rules = reactive({
    name: [Rules.simpleName],
    path: [Rules.requiredInput, Rules.noSpace],
});

type FormInstance = InstanceType<typeof ElForm>;
const formRef = ref<FormInstance>();

const hasInfectedDir = () => {
    return (
        dialogData.value.rowData!.infectedStrategy === 'move' || dialogData.value.rowData!.infectedStrategy === 'copy'
    );
};
const loadDir = async (path: string) => {
    dialogData.value.rowData!.path = path;
};
const loadInfectedDir = async (path: string) => {
    dialogData.value.rowData!.infectedDir = path;
};

const onSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    formEl.validate(async (valid) => {
        if (!valid) return;
        loading.value = true;
        dialogData.value.rowData.timeout = transferTimeToSecond(
            dialogData.value.rowData.timeoutItem + dialogData.value.rowData.timeoutUnit,
        );

        if (dialogData.value.title === 'edit') {
            await updateClam(dialogData.value.rowData)
                .then(() => {
                    loading.value = false;
                    drawerVisible.value = false;
                    MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
                    emit('search');
                })
                .catch(() => {
                    loading.value = false;
                });

            return;
        }

        await createClam(dialogData.value.rowData)
            .then(() => {
                loading.value = false;
                MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
                emit('search');
                drawerVisible.value = false;
            })
            .catch(() => {
                loading.value = false;
            });
    });
};

defineExpose({
    acceptParams,
});
</script>
