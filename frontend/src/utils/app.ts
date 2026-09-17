import { jumpToPath } from './router';
import router from '@/routers';

export const jumpToInstall = (type: string) => {
    switch (type) {
        case 'php':
        case 'node':
        case 'java':
        case 'go':
        case 'python':
        case 'dotnet':
            jumpToPath(router, '/websites/runtimes/' + type);
            return true;
    }
    return false;
};
