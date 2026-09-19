import http from '@/api';
import { ResPage } from '../interface';
import { Log } from '../interface/log';
import { TimeoutEnum } from '@/enums/http-enum';

export const getOperationLogs = (info: Log.SearchOpLog) => {
    return http.post<ResPage<Log.OperationLog>>(`/core/logs/operation`, info);
};

export const getLoginLogs = (info: Log.SearchLgLog) => {
    return http.post<ResPage<Log.LoginLogs>>(`/core/logs/login`, info, undefined, undefined);
};

export const getSystemFiles = () => {
    return http.get<Array<string>>(`/logs/system/files`);
};

export const getSystemLogStatus = () => {
    const query = '';
    return http.get<Log.SystemLogStatus>(`/logs/system/status${query}`);
};

export const readSystemLogs = (params: Log.SystemLogSearch) => {
    const query = '';
    return http.post<Log.SystemLog>(`/logs/system/read${query}`, params);
};

export const listRunningServices = () => {
    const query = '';
    return http.get<string[]>(`/logs/system/services${query}`);
};

export const cleanLogs = (param: Log.CleanLog) => {
    return http.post(`/core/logs/clean`, param);
};

export const searchTasks = (req: Log.SearchTaskReq) => {
    return http.post<ResPage<Log.Task>>(`/logs/tasks/search`, req);
};

export const readTaskLogByLine = (req: Log.TaskLogReadReq) => {
    return http.post<any>(`/logs/tasks/read`, req, TimeoutEnum.T_40S);
};

export const countExecutingTask = () => {
    return http.get<number>(`/logs/tasks/executing/count`);
};
