import { reactive, type UnwrapNestedRefs } from 'vue';
import { useRoute } from 'vue-router';
import { getPageState } from '@/utils/page-state-cache';

export const usePageState = <T extends object>(factory: () => T): UnwrapNestedRefs<T> => {
    const route = useRoute();
    const key = String(route.name || route.path);
    return reactive(getPageState(key, factory));
};
