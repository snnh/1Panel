import { reactive, type UnwrapNestedRefs } from 'vue';
import { useRoute } from 'vue-router';
import { getPageState } from '@/utils/page-state-cache';

/** A new page visit starts from its defaults instead of restoring the previous page number. */
export const usePageState = <T extends object>(factory: () => T): UnwrapNestedRefs<T> => {
    const route = useRoute();
    const key = String(route.name || route.path);
    return reactive(getPageState(key, factory));
};
