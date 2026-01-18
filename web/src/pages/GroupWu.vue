<script setup>
import { ref, onMounted, h } from "vue";
import {  useMessage, useDialog, NDataTable } from 'naive-ui'
import { ApiGetGroupWu } from '@/api'
import * as XLSX from 'xlsx';

const nmessage = useMessage()
const groupdata = ref([])

function getData() {
    groupdata.value = [];
    ApiGetGroupWu().then(v => {
        if (v.status == 200) {
            let resp = v.data;
            if (resp.code == 200) {
                let data = resp.data;
                groupdata.value = data;
            } else {
                nmessage.error(resp.msg);
            }
        } else {
            nmessage.error("获取分组武勋数据失败");
        }
    });
}

onMounted(() => {
    getData();
});

// 定义自定义排序图标
function renderSorterIcon({ order }) {
  if (order === false) {
    return h('span', { style: { color: '#666' } }, '↕')
  }
  if (order === 'ascend') {
    return h('span', { style: { color: '#333' } }, '↑')
  }
  if (order === 'descend') {
    return h('span', { style: { color: '#333' } }, '↓')
  }
}
</script>

<template>
    <div class="bikamoeapp">
        <div class="bikamoeapp-content">
            <div class="bikamoeapp-title">
                <h2 style="margin-bottom: 4px;">攻城考勤助手</h2>
                <p>分组本周武勋</p>
                <p>更新武勋数据请同步成员数据</p>
            </div>
            <!-- <div class="bikamoeapp-list"> -->
            <div>
                <div style="width: 100%;
                    height: 48px;
                    border-bottom: 1px solid rgba(228, 228, 231, 0.6);
                    display: flex;
                    align-items: center;
                    padding: 0 8px;
                    box-sizing: border-box;">
                    <router-link class="button" to="/">返回</router-link>
                    <router-link class="button" to="/wuHistory">
                        历史武勋
                    </router-link>
                    <a class="button" @click="getData">
                        刷新
                    </a>
                </div>
                <div>
                  <n-data-table
                    :columns="[
                      { title: '分组名称', key: 'group' },
                      { title: '人数', key: 'member_count' },
                      { 
                        title: '总武勋', 
                        key: 'total_wu', 
                        sorter: (row1, row2) => row1.total_wu - row2.total_wu,
                        defaultSortOrder: false,
                        renderSorterIcon: renderSorterIcon
                      },
                      { 
                        title: '平均武勋', 
                        key: 'average_wu', 
                        sorter: (row1, row2) => row1.average_wu - row2.average_wu,
                        defaultSortOrder: false,
                        renderSorterIcon: renderSorterIcon
                      },
                      { title: '0武勋人数', key: 'zero_wu_count' }
                    ]"
                    :data="groupdata"
                    :pagination="{ pageSize: 100 }"
                    />
                </div>
            </div>

            <!-- </div> -->
        </div>
    </div>
</template>

<style scoped>
a.button {
    border-bottom: 1px solid rgb(228 228 231 / 60%);
    margin-right: 8px;
    font-size: 14px;
}
</style>