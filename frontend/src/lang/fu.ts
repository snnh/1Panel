type FuLocaleMessage = Record<string, unknown>;

const fuLocales: Record<string, FuLocaleMessage> = {
    en: {
        fu: {
            table: {
                more: 'More',
                custom_table_rows: 'Custom columns',
            },
            steps: {
                cancel: 'Cancel',
                prev: 'Previous',
                next: 'Next',
                finish: 'Finish',
            },
        },
    },
    zh: {
        fu: {
            table: {
                more: '更多',
                custom_table_rows: '自定义列',
            },
            steps: {
                cancel: '取消',
                prev: '上一步',
                next: '下一步',
                finish: '完成',
            },
        },
    },
};

export const getFuLocaleMessage = (locale: string) => {
    return fuLocales[locale] || fuLocales.en;
};
