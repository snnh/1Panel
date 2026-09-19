<template>
    <div class="logo" style="cursor: pointer" @click="goHome">
        <template v-if="isCollapse">
            <MenuLogo />
        </template>
        <template v-else>
            <PrimaryLogo />
        </template>
    </div>
</template>

<script setup lang="ts">
import PrimaryLogo from '@/assets/images/1panel-logo.svg?component';
import MenuLogo from '@/assets/images/1panel-menu-logo.svg?component';
import { routerToNameWithQuery } from '@/utils/router';

defineProps<{ isCollapse: boolean }>();

const goHome = () => {
    routerToNameWithQuery('home', { t: Date.now() });
};
</script>

<style scoped lang="scss">
.logo {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 49px;
    z-index: 1;
    img {
        object-fit: contain;
        width: 95%;
        height: 45px;
    }
}
</style>
