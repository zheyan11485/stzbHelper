<script setup>
import { ref, onMounted } from "vue";
import { useMessage, NTable, NDatePicker, NSelect, NButton, NSpace } from 'naive-ui';
import { ApiGetWuHistory } from '@/api';

const nmessage = useMessage();
const historydata = ref([]);
const selectedGroup = ref('');
const startDate = ref(null);
const endDate = ref(null);
const allGroups = ref([]);
const loadingHistory = ref(false);

// 获取历史数据
function getHistoryData() {
    loadingHistory.value = true;
    
    const params = {};
    if (selectedGroup.value) params.group = selectedGroup.value;
    if (startDate.value) {
        // 将时间戳转换为日期字符串格式 (YYYY-MM-DD)
        const start = new Date(startDate.value);
        params.start_date = `${start.getFullYear()}-${String(start.getMonth() + 1).padStart(2, '0')}-${String(start.getDate()).padStart(2, '0')}`;
    }
    if (endDate.value) {
        // 将时间戳转换为日期字符串格式 (YYYY-MM-DD)
        const end = new Date(endDate.value);
        params.end_date = `${end.getFullYear()}-${String(end.getMonth() + 1).padStart(2, '0')}-${String(end.getDate()).padStart(2, '0')}`;
    }
    
    ApiGetWuHistory(params).then(v => {
        if (v.status == 200) {
            let resp = v.data;
            if (resp.code == 200) {
                let data = resp.data;
                // 按日期排序，最新的在前
                historydata.value = data.sort((a, b) => new Date(b.record_date) - new Date(a.record_date));
                if (data.length === 0) {
                    nmessage.info('暂无历史数据');
                }
            } else {
                nmessage.error(resp.msg);
            }
        } else {
            nmessage.error("获取历史武勋数据失败");
        }
    }).finally(() => {
        loadingHistory.value = false;
    });
}

// 获取所有分组
function getAllGroups() {
    import('@/api').then(apiModule => {
        apiModule.ApiGetTeamGroup().then(v => {
            if (v.status == 200) {
                let resp = v.data;
                if (resp.code == 200) {
                    allGroups.value = resp.data.map(group => ({ label: group, value: group }));
                }
            }
        });
    });
}

// 计算今天的数据
function getTodayData() {
    // 计算今天的日期
    const today = new Date();
    
    // 设置日期选择器的值为今天
    startDate.value = today.getTime();
    endDate.value = today.getTime();
}

// 计算本周一到周日的数据
function getThisWeekData() {
    // 计算本周周一和周日的日期
    const today = new Date();
    const dayOfWeek = today.getDay(); // 0 (Sunday) to 6 (Saturday)
    const diffToMonday = dayOfWeek === 0 ? -6 : 1 - dayOfWeek; // 计算到周一的差值
    
    const monday = new Date(today);
    monday.setDate(today.getDate() + diffToMonday);
    monday.setHours(0, 0, 0, 0);
    
    const sunday = new Date(monday);
    sunday.setDate(monday.getDate() + 6);
    sunday.setHours(23, 59, 59, 999);
    
    // 设置日期选择器的值
    startDate.value = monday.getTime();
    endDate.value = sunday.getTime();
}

// 计算上周一到周日的数据
function getLastWeekData() {
    // 计算上周周一和周日的日期
    const today = new Date();
    const dayOfWeek = today.getDay(); // 0 (Sunday) to 6 (Saturday)
    const diffToLastMonday = dayOfWeek === 0 ? -6 : 1 - dayOfWeek; // 计算到上周一的差值
    
    const lastMonday = new Date(today);
    lastMonday.setDate(today.getDate() + diffToLastMonday - 7); // 减去7天得到上周一
    lastMonday.setHours(0, 0, 0, 0);
    
    const lastSunday = new Date(lastMonday);
    lastSunday.setDate(lastMonday.getDate() + 6);
    lastSunday.setHours(23, 59, 59, 999);
    
    // 设置日期选择器的值
    startDate.value = lastMonday.getTime();
    endDate.value = lastSunday.getTime();
}

// 计算前一天的总武勋值，用于趋势比较
function getPrevDayWu(index, group) {
    if (index >= historydata.value.length - 1) return null;
    
    // 查找同一分组的前一天数据
    const currentGroup = historydata.value[index].group_name;
    for (let i = index + 1; i < historydata.value.length; i++) {
        if (historydata.value[i].group_name === currentGroup) {
            return historydata.value[i].total_wu;
        }
    }
    return null;
}

onMounted(() => {
    getAllGroups();
});
</script>

<template>
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">历史武勋统计</h2>
                <p>查看武勋历史数据及变化趋势</p>
            </div>
            <div>
                <div style="width: 100%;
                    height: 48px;
                    border-bottom: 1px solid rgba(228, 228, 231, 0.6);
                    display: flex;
                    align-items: center;
                    padding: 0 8px;
                    box-sizing: border-box; flex-wrap: wrap;">
                    <router-link class="button" to="/">返回</router-link>
                    <a class="button" @click="getTodayData">
                        今天数据
                    </a>
                    <a class="button" @click="getThisWeekData">
                        本周数据
                    </a>
                    <a class="button" @click="getLastWeekData">
                        上周数据
                    </a>
                </div>
                
                <div style="padding: 10px; background: #f5f5f5; margin-bottom: 10px; display: flex; flex-wrap: wrap; gap: 10px; align-items: center;">
                    <n-select 
                        v-model:value="selectedGroup" 
                        :options="allGroups" 
                        placeholder="选择分组（可选）"
                        style="min-width: 200px; flex: 1; max-width: 200px;"
                        clearable
                    />
                    <n-date-picker
                        v-model:value="startDate"
                        type="date"
                        placeholder="开始日期"
                        style="min-width: 200px; flex: 1; max-width: 200px;"
                        :allow-input="false"
                        clearable
                    />
                    <n-date-picker
                        v-model:value="endDate"
                        type="date"
                        placeholder="结束日期"
                        style="min-width: 200px; flex: 1; max-width: 200px;"
                        :allow-input="false"
                        clearable
                    />
                    <n-button @click="getHistoryData" type="primary" :loading="loadingHistory">查询</n-button>
                </div>
                
                <div v-if="historydata.length > 0">
                    <n-table>
                        <thead>
                            <tr>
                                <th>分组名称</th>
                                <th>人数</th>
                                <th>总武勋</th>
                                <th>平均武勋</th>
                                <th>0武勋人数</th>
                                <th>记录日期</th>
                                <th>当日变化</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(u, index) in historydata" :key="u.id">
                                <td>{{ u.group_name }}</td>
                                <td>{{ u.member_count }}</td>
                                <td>{{ u.total_wu }}</td>
                                <td>{{ u.average_wu }}</td>
                                <td>{{ u.zero_wu_count }}</td>
                                <td>{{ new Date(u.record_date).toLocaleDateString() }}</td>
                                <td>
                                    <span v-if="getPrevDayWu(index, u.group_name) !== null">
                                        <span v-if="u.total_wu - getPrevDayWu(index, u.group_name) > 0" style="color: green;">+{{ u.total_wu - getPrevDayWu(index, u.group_name) }}</span>
                                        <span v-else-if="u.total_wu - getPrevDayWu(index, u.group_name) < 0" style="color: red;">{{ u.total_wu - getPrevDayWu(index, u.group_name) }}</span>
                                        <span v-else style="color: gray;">0</span>
                                    </span>
                                    <span v-else style="color: gray;">—</span>
                                </td>
                            </tr>
                        </tbody>
                    </n-table>
                    <p style="text-align: center; margin-top: 10px; color: #666;">
                        共 {{ historydata.length }} 条记录
                    </p>
                </div>
                <div v-else style="text-align: center; padding: 20px;">
                    <p>暂无历史数据</p>
                    <p style="color: #999; font-size: 14px; margin-top: 10px;">
                        请选择日期范围并点击"查询"按钮
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
a.button {
    border-bottom: 1px solid rgb(228 228 231 / 60%);
    margin-right: 8px;
    font-size: 14px;
    cursor: pointer;
    padding: 4px 8px;
}
</style>