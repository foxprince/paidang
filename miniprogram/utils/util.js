const STATUS_TEXT = {
  pending: '待确认',
  confirmed: '已确认',
  completed: '已完成',
  cancelled: '已取消',
};

const WEEKDAYS = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];

// 分 → 元
function fen2yuan(fen) {
  return ((fen || 0) / 100).toFixed(2);
}

// 今天 YYYY-MM-DD
function todayStr() {
  const d = new Date();
  return fmtDate(d);
}

function fmtDate(d) {
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${d.getFullYear()}-${m}-${day}`;
}

// 未来 n 天（含今天）
function nextDays(n) {
  const out = [];
  const d = new Date();
  for (let i = 0; i < n; i++) {
    const t = new Date(d.getTime() + i * 86400000);
    out.push({
      date: fmtDate(t),
      label: i === 0 ? '今天' : `${t.getMonth() + 1}/${t.getDate()} ${WEEKDAYS[t.getDay()]}`,
    });
  }
  return out;
}

function uuid() {
  return 'xxxxxxxx-xxxx-4xxx'.replace(/x/g, () =>
    Math.floor(Math.random() * 16).toString(16)
  ) + Date.now().toString(16);
}

function toast(msg) {
  wx.showToast({ title: msg || '出错了', icon: 'none' });
}

module.exports = { STATUS_TEXT, WEEKDAYS, fen2yuan, todayStr, fmtDate, nextDays, uuid, toast };
